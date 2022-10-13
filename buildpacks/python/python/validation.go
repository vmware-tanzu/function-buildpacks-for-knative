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

	"kn-fn/buildpacks/command"
)

type FunctionValidationLayer struct {
	layerContributor libpak.LayerContributor
	logger           bard.Logger
	commandRunner    command.Runner

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

func NewFunctionValidationLayer(appPath string, invoker Layer, InvokerDepLayer Layer, commandRunner command.Runner, opts ...FunctionValidationOpts) *FunctionValidationLayer {
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

		cmd := exec.Command(
			"python", "-m", "pyfunc",
			"check",
			"-s", f.applicationPath,
			"-m", f.module,
			"-f", f.function,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHONPATH=%s", strings.Join(pythonPath, string(os.PathListSeparator))))

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

func WithValidationFunctionClass(moduleName string, functionName string) FunctionValidationOpts {
	return func(vl *FunctionValidationLayer, metadata map[string]string) {
		vl.module = moduleName
		vl.function = functionName
	}
}

func WithValidationFunctionEnvs(envs map[string]any) FunctionValidationOpts {
	return func(vl *FunctionValidationLayer, metadata map[string]string) {
		for name, value := range envs {
			metadata[name] = value.(string)
		}
	}
}
