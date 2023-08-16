// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"sort"

	"kn-fn/buildpacks/command"
)

type Build struct {
	Logger        bard.Logger
	CommandRunner command.Runner
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	functionClass, _ := cr.Resolve("BP_FUNCTION")
	moduleName, functionName, err := parseFunctionClass(functionClass)
	if err != nil {
		return libcnb.BuildResult{}, err
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

	invokerDeps, err := dependencyResolver.Resolve("invoker-deps", "")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	invoker, err := dependencyResolver.Resolve("invoker", "")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	invokerDepCacheLayer := NewInvokerDependencyCache(invokerDeps, dependencyCache, b.CommandRunner)
	invokerDepCacheLayer.Logger = b.Logger
	result.Layers = append(result.Layers, invokerDepCacheLayer)

	invokerLayer := NewInvoker(invoker, dependencyCache)
	result.Layers = append(result.Layers, invokerLayer)

	functionPlan, ok, err := planResolver.Resolve("python-function")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve python-function plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}

	if functionPlan.Metadata["func_yaml_envs"] == nil {
		functionPlan.Metadata["func_yaml_envs"] = map[string]any{}
	}

	functionLayer := NewFunction(
		WithLogger(b.Logger),
		WithFuncYamlEnvs(functionPlan.Metadata["func_yaml_envs"].(map[string]any)),
	)
	result.Layers = append(result.Layers, functionLayer)

	validationLayer := NewFunctionValidationLayer(
		context.Application.Path,
		invokerLayer,
		invokerDepCacheLayer,
		b.CommandRunner,
		WithValidationLogger(b.Logger),
		WithValidationFunctionClass(moduleName, functionName),
		WithValidationFunctionEnvs(functionPlan.Metadata["func_yaml_envs"].(map[string]any)),
	)
	result.Layers = append(result.Layers, validationLayer)

	if functionPlan.Metadata["func_yaml_options"] != nil {
		for optionName, optionValue := range functionPlan.Metadata["func_yaml_options"].(map[string]any) {
			result.Labels = append(result.Labels, libcnb.Label{
				Key:   optionName,
				Value: optionValue.(string),
			})
		}
		sort.Slice(result.Labels, func(i, j int) bool { return result.Labels[i].Key < result.Labels[j].Key })
	}

	cmd := "python"
	arguments := []string{"-m", "pyfunc", "start", "-m", moduleName, "-f", functionName}
	result.Processes = append(result.Processes,
		libcnb.Process{
			Default:   true,
			Type:      "func",
			Command:   cmd,
			Arguments: arguments,
		},
	)

	return result, nil
}
