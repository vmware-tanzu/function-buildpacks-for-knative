// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type FunctionValidationLayer struct {
	layerContributor libpak.LayerContributor
	logger           bard.Logger
	commandRunner    CommandRunner

	module          string
	function        string
	applicationPath string

	override bool

	Invoker    Layer
	InvokerDep Layer
}

type FunctionValidationOpts func(*FunctionValidationLayer, map[string]string)

//go:generate mockgen -destination ../mock_python/layer.go . Layer
type Layer interface {
	PythonPath() string
}

func NewFunctionValidationLayer(appPath string, invoker Layer, InvokerDepLayer Layer, commandRunner CommandRunner, opts ...FunctionValidationOpts) *FunctionValidationLayer {
	fvl := &FunctionValidationLayer{
		applicationPath: appPath,
		Invoker:         invoker,
		InvokerDep:      InvokerDepLayer,
		commandRunner:   commandRunner,
	}
	meta := map[string]string{}

	for _, opt := range opts {
		opt(fvl, meta)
	}

	fvl.layerContributor = libpak.NewLayerContributor("validation", meta, libcnb.LayerTypes{
		Build: true,
	})

	return fvl
}

func (f *FunctionValidationLayer) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.layerContributor.Logger = f.logger

	return f.layerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		f.logger.Body("Validating function")

		var pythonPath []string
		if f.Invoker.PythonPath() != "" {
			pythonPath = append(pythonPath, f.Invoker.PythonPath())
		}

		if f.InvokerDep.PythonPath() != "" {
			pythonPath = append(pythonPath, f.InvokerDep.PythonPath())
		}

		if env, found := os.LookupEnv("PYTHONPATH"); found {
			pythonPath = append(pythonPath, env)
		}

		cmd := exec.Command("python", "-m", "pyfunc", "check", "-s", f.applicationPath)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHONPATH=%s", strings.Join(pythonPath, string(os.PathListSeparator))))
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", EnvModuleName, f.module), fmt.Sprintf("%s=%s", EnvFunctionName, f.function))

		if output, err := f.commandRunner.Run(cmd); err != nil {
			return layer, fmt.Errorf("%v: %v", output, err)
		}

		f.logger.Debug("Function was successfully parsed")
		return layer, nil
	})
}

func (f *FunctionValidationLayer) Name() string {
	return f.layerContributor.Name
}

func WithValidationLogger(logger bard.Logger) FunctionValidationOpts {
	return func(vl *FunctionValidationLayer, metadata map[string]string) {
		vl.logger = logger
	}
}

func WithValidationFunctionClass(functionClass string, override bool) FunctionValidationOpts {
	return func(vl *FunctionValidationLayer, metadata map[string]string) {
		vl.override = override

		fSplit := strings.Split(functionClass, ".")
		if override || (vl.module == "" && vl.function == "") {
			vl.module = fSplit[0]
			vl.function = fSplit[1]
			metadata[EnvModuleName] = fSplit[0]
			metadata[EnvFunctionName] = fSplit[1]
		}
	}
}

func WithValidationFunctionEnvs(envs map[string]any) FunctionValidationOpts {
	return func(vl *FunctionValidationLayer, metadata map[string]string) {
		for name, value := range envs {
			if name == EnvModuleName && !vl.override {
				vl.module = value.(string)
			} else if name == EnvFunctionName && !vl.override {
				vl.function = value.(string)
			} else {
				metadata[name] = value.(string)
			}
		}
	}
}
