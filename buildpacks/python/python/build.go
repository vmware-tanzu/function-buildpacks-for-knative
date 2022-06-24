// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	_, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)

	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	planResolver := libpak.PlanEntryResolver{Plan: context.Plan}

	dependencyCache, err := libpak.NewDependencyCache(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}
	dependencyCache.Logger = b.Logger

	dependencyResolver, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	invokerDepCache, err := dependencyResolver.Resolve("invoker-deps", "")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	invoker, err := dependencyResolver.Resolve("invoker", "")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	invokerDepCacheLayer, invokerDepCacheBom := NewInvokerDependencyCache(invokerDepCache, dependencyCache)
	invokerDepCacheLayer.Logger = b.Logger
	result.Layers = append(result.Layers, invokerDepCacheLayer)
	result.BOM.Entries = append(result.BOM.Entries, invokerDepCacheBom)

	invokerLayer, invokerBOM := NewInvoker(invoker, dependencyCache)
	invokerLayer.Logger = b.Logger
	invokerLayer.DependencyCacheLayer = &invokerDepCacheLayer
	result.Layers = append(result.Layers, invokerLayer)
	result.BOM.Entries = append(result.BOM.Entries, invokerBOM)

	functionPlan, ok, err := planResolver.Resolve("python-function")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve python-function plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}
	functionLayer := NewFunction(functionPlan)
	result.Layers = append(result.Layers, functionLayer)

	validationLayer := NewFunctionValidationLayer(functionPlan, context.Application.Path)
	result.Layers = append(result.Layers, validationLayer)
	result.Labels = b.getFuncYamlOptions(context.Application.Path)

	command := "python"
	arguments := []string{"-m", "pyfunc", "start"}
	result.Processes = append(result.Processes,
		libcnb.Process{
			Default:   true,
			Type:      "func",
			Command:   command,
			Arguments: arguments,
		},
		libcnb.Process{
			Type:    "shell",
			Command: "bash",
		},
	)

	return result, nil
}

func (b Build) getFuncYamlOptions(appPath string) []libcnb.Label {
	configFile := filepath.Join(appPath, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		b.Logger.Bodyf("'%s' not detected", knfn.ConfigFile)
		return []libcnb.Label{}
	}

	f, err := knfn.NewFunction(appPath)
	if err != nil {
		b.Logger.Bodyf("unable to parse '%s': %v", knfn.ConfigFile, err)
		return []libcnb.Label{}
	}
	return b.optionsToLabels(f.Options)
}

func (b Build) optionsToLabels(options knfn.Options) []libcnb.Label {
	labels := []libcnb.Label{}

	scaleJson, err := json.Marshal(options.Scale)
	if err != nil {
		b.Logger.Bodyf("unable to marshal func.yaml options.Scale")
	}
	requestsJson, err := json.Marshal(options.Resources.Requests)
	if err != nil {
		b.Logger.Bodyf("unable to marshal func.yaml options.Resources.Requests")

	}
	limitsJson, err := json.Marshal(options.Resources.Limits)
	if err != nil {
		b.Logger.Bodyf("unable to marshal func.yaml options.Resources.Limits")
	}
	labels = append(labels,
		libcnb.Label{
			Key:   "options-scale",
			Value: string(scaleJson),
		},
		libcnb.Label{
			Key:   "options-resources-requests",
			Value: string(requestsJson),
		},
		libcnb.Label{
			Key:   "options-resources-limits",
			Value: string(limitsJson),
		})

	return labels
}
