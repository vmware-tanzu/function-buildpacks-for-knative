// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

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

	scaleMap := scaleToMap(*options.Scale)
	requestsMap := requestsToMap(options.Resources.Requests)
	limitsMap := limitsToMap(options.Resources.Limits)

	for k, v := range scaleMap {
		labels = append(labels, libcnb.Label{
			Key:   "scale-" + k,
			Value: v,
		})
	}

	for k, v := range requestsMap {
		labels = append(labels, libcnb.Label{
			Key:   "resource-requests-" + k,
			Value: v,
		})
	}

	for k, v := range limitsMap {
		labels = append(labels, libcnb.Label{
			Key:   "resource-limits-" + k,
			Value: v,
		})
	}

	return labels
}

func scaleToMap(input knfn.ScaleOptions) map[string]string {
	result := map[string]string{}
	fields := reflect.TypeOf(input)
	values := reflect.ValueOf(input)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if value.IsValid() {
			result[field.Name] = value.String()

			if value.Type() == reflect.TypeOf(100) {
				result[field.Name] = strconv.FormatInt(value.Int(), 10)
			}
			if value.Type() == reflect.TypeOf(100.1) {
				result[field.Name] = fmt.Sprintf("%f", value.Float())
			}
		}
	}
	return result
}

func requestsToMap(input *knfn.ResourcesRequestsOptions) map[string]string {
	result := map[string]string{}
	fields := reflect.TypeOf(input)
	values := reflect.ValueOf(input)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if value.IsValid() {
			result[field.Name] = value.String()

			if value.Type() == reflect.TypeOf(100) {
				result[field.Name] = strconv.FormatInt(value.Int(), 10)
			}
			if value.Type() == reflect.TypeOf(100.1) {
				result[field.Name] = fmt.Sprintf("%f", value.Float())
			}
		}
	}
	return result
}

func limitsToMap(input *knfn.ResourcesLimitsOptions) map[string]string {
	result := map[string]string{}
	fields := reflect.TypeOf(input)
	values := reflect.ValueOf(input)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if value.IsValid() {
			result[field.Name] = value.String()

			if value.Type() == reflect.TypeOf(100) {
				result[field.Name] = strconv.FormatInt(value.Int(), 10)
			}
			if value.Type() == reflect.TypeOf(100.1) {
				result[field.Name] = fmt.Sprintf("%f", value.Float())
			}
		}
	}
	return result
}
