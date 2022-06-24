// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc, err := libpak.NewDependencyCache(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}
	dc.Logger = b.Logger

	e, ok, err := pr.Resolve("java-function")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve java-function plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}
	if e.Metadata["has_func_yaml"] == true {
		envs := NewFuncYamlEnvs(context.Application.Path)
		envs.Logger = b.Logger
		result.Layers = append(result.Layers, envs)
		result.Labels = b.getFuncYamlOptions(context.Application.Path)
	}

	dep, err := dr.Resolve("invoker", "")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	r, err := regexp.Compile(`tomcat-embed-core-[\d\.]+\.jar`)
	if err != nil {
		return libcnb.BuildResult{}, err
	}
	dependencyPath, err := findPath(context.Application.Path, r)
	if err != nil {
		return libcnb.BuildResult{}, err
	}

	if dependencyPath != "" {
		b.Logger.Info("embedded tomcat dependency found at: ", dependencyPath, " (skipping invoker layer)")
	} else {
		i, be := NewInvoker(dep, dc)
		i.Logger = b.Logger
		result.Layers = append(result.Layers, i)
		result.BOM.Entries = append(result.BOM.Entries, be)
	}

	funcDef, _ := cr.Resolve("BP_FUNCTION")
	defaultDef, _ := cr.Resolve("BP_DEFAULT_FUNCTION")

	f := NewFunction(funcDef, defaultDef, context.Application.Path)
	f.Logger = b.Logger
	result.Layers = append(result.Layers, f)

	command := "java"
	arguments := []string{"org.springframework.boot.loader.JarLauncher"}
	result.Processes = append(result.Processes,
		libcnb.Process{Type: "func", Command: command, Arguments: arguments, Default: true},
	)

	return result, nil
}

func findPath(path string, r *regexp.Regexp) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			if found, err := findPath(strings.Join([]string{path, file.Name()}, "/"), r); found != "" {
				return found, err
			}
		} else {
			if r.MatchString(file.Name()) {
				return strings.Join([]string{path, file.Name()}, "/"), nil
			}
		}
	}
	return "", nil
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
