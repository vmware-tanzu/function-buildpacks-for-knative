// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
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

	invokerDepCacheLayer := NewInvokerDependencyCache(invokerDepCache, dependencyCache)
	invokerDepCacheLayer.Logger = b.Logger
	result.Layers = append(result.Layers, &invokerDepCacheLayer)

	invokerLayer := NewInvoker(invoker, dependencyCache)
	invokerLayer.Logger = b.Logger
	result.Layers = append(result.Layers, &invokerLayer)

	functionPlan, ok, err := planResolver.Resolve("python-function")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve python-function plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}

	functionClass, isFuncDefDefault := cr.Resolve("BP_FUNCTION")
	functionLayer := NewFunction(
		WithLogger(b.Logger),
		WithFunctionClass(functionClass, isFuncDefDefault),
		WithFuncYamlEnvs(functionPlan.Metadata["func_yaml_envs"].(map[string]interface{})),
	)
	result.Layers = append(result.Layers, functionLayer)

	validationLayer := NewFunctionValidationLayer(
		context.Application.Path,
		&invokerLayer,
		&invokerDepCacheLayer,
		WithValidationLogger(b.Logger),
		WithValidationFunctionClass(functionClass, isFuncDefDefault),
		WithValidationFunctionEnvs(functionPlan.Metadata["func_yaml_envs"].(map[string]interface{})),
	)
	result.Layers = append(result.Layers, validationLayer)

	for optionName, optionValue := range functionPlan.Metadata["func_yaml_options"].(map[string]interface{}) {
		result.Labels = append(result.Labels, libcnb.Label{
			Key:   optionName,
			Value: optionValue.(string),
		})
	}

	command := "python"
	arguments := []string{"-m", "pyfunc", "start"}
	result.Processes = append(result.Processes,
		libcnb.Process{
			Default:   true,
			Type:      "func",
			Command:   command,
			Arguments: arguments,
		},
	)

	return result, nil
}
