// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"bytes"
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

	module          string
	function        string
	applicationPath string

	override bool

	Invoker    Layer
	InvokerDep Layer
}

type FunctionValidationOpts func(*FunctionValidationLayer, map[string]string)

type Layer interface {
	PythonPath() string
}

func NewFunctionValidationLayer(appPath string, invoker Layer, InvokerDepLayer Layer, opts ...FunctionValidationOpts) FunctionValidationLayer {
	fvl := FunctionValidationLayer{
		applicationPath: appPath,
		Invoker:         invoker,
		InvokerDep:      InvokerDepLayer,
	}
	meta := map[string]string{}

	for _, opt := range opts {
		opt(&fvl, meta)
	}

	fvl.layerContributor = libpak.NewLayerContributor("validation", meta, libcnb.LayerTypes{
		Build: true,
	})

	return fvl
}

func (i FunctionValidationLayer) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.layerContributor.Logger = i.logger

	return i.layerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		i.logger.Body("Validating function")

		var pythonPath []string
		if i.Invoker.PythonPath() != "" {
			pythonPath = append(pythonPath, i.Invoker.PythonPath())
		}

		if i.InvokerDep.PythonPath() != "" {
			pythonPath = append(pythonPath, i.InvokerDep.PythonPath())
		}

		if env, found := os.LookupEnv("PYTHONPATH"); found {
			pythonPath = append(pythonPath, env)
		}

		buffer := bytes.NewBuffer(nil)
		cmd := exec.Command("python", "-m", "pyfunc", "check", "-s", i.applicationPath)
		// cmd := exec.Command("env")
		cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHONPATH=%s", strings.Join(pythonPath, string(os.PathListSeparator))))
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", EnvModuleName, i.module), fmt.Sprintf("%s=%s", EnvFunctionName, i.function))
		cmd.Stderr = buffer
		cmd.Stdout = buffer

		if err := cmd.Run(); err != nil {
			return layer, fmt.Errorf("%v: %v", buffer.String(), err)
		}

		i.logger.Debug("Function was successfully parsed")
		return layer, nil
	})
}

func (i FunctionValidationLayer) Name() string {
	return i.layerContributor.Name
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
