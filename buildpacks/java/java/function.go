// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"strconv"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Function struct {
	layerContributor libpak.LayerContributor
	logger           bard.Logger

	functionClass         string
	overrideFunctionClass bool

	defaultFunctionName         string
	overrideDefaultFunctionName bool

	funcYamlEnvs map[string]string

	applicationPath string
}

func NewFunction(applicationPath string, opts ...FunctionOpt) *Function {
	f := &Function{
		applicationPath: applicationPath,
	}
	meta := map[string]string{}

	for _, opt := range opts {
		opt(f, meta)
	}

	f.layerContributor = libpak.NewLayerContributor(
		"java-function",
		meta,
		libcnb.LayerTypes{
			Launch: true,
		},
	)

	return f
}

func (f *Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.layerContributor.Logger = f.logger

	return f.layerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_LOCATION", f.applicationPath)

		if f.overrideFunctionClass {
			layer.LaunchEnvironment.Override("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", f.functionClass)
		} else {
			layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", f.functionClass)
		}

		for envName, envValue := range f.funcYamlEnvs {
			layer.LaunchEnvironment.Default(envName, envValue)
		}

		return layer, nil
	})
}

func (f *Function) Name() string {
	return f.layerContributor.Name
}

type FunctionContribOpt func(layer *libcnb.Layer)
type FunctionOpt func(fun *Function, metadata map[string]string)

func WithLogger(logger bard.Logger) FunctionOpt {
	return func(fun *Function, metadata map[string]string) {
		fun.logger = logger
	}
}

func WithFunctionClass(functionClass string, override bool) FunctionOpt {
	return func(fun *Function, metadata map[string]string) {
		fun.functionClass = functionClass
		fun.overrideFunctionClass = override

		metadata["bp-function-class"] = functionClass
		metadata["bp-function-class-override"] = strconv.FormatBool(override)
	}
}

func WithFuncYamlEnvs(funcYamlEnvs map[string]any) FunctionOpt {
	return func(fun *Function, metadata map[string]string) {
		fun.funcYamlEnvs = map[string]string{}

		for envName, envValue := range funcYamlEnvs {
			value := envValue.(string)
			metadata[envName] = value
			fun.funcYamlEnvs[envName] = value
		}
	}
}
