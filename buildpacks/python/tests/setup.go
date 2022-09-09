// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

const (
	HTTPFuncTemplate = `from typing import Any

def {{.Name}}({{range $var, $type := .Arguments}}{{$var}}:{{$type}}{{end}}):
	return "{{.ReturnValue}}"`
)

type HTTPFunction struct {
	Module      string
	Name        string
	Arguments   map[string]string
	ReturnValue string
}

type Function interface {
	Generate(path, funcTemplate string) error
}

func (f HTTPFunction) Generate(path, funcTemplate string) error {
	templ, err := template.New("http-func-template").Parse(funcTemplate)
	if err != nil {
		return err
	}

	modulePath := filepath.Join(path, fmt.Sprintf("%s.py", f.Module))
	file, err := os.Create(modulePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return templ.Execute(file, f)
}

func NewLogger() bard.Logger {
	buf := bytes.NewBuffer(nil)
	return bard.NewLogger(buf)
}

type SetupOpts func(directory string)

func SetupTestDirectory(opts ...SetupOpts) (string, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "python-functions-buildpack-*")
	if err != nil {
		panic(fmt.Sprintf("unable to create test directory: %v", err))
	}

	for _, opt := range opts {
		opt(dir)
	}

	cleanup := func() {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("Failed to delete temp directory %s: %v", dir, err)
		}

	}
	return dir, cleanup
}

func WithFuncYaml() SetupOpts {
	return func(directory string) {
		cfg, err := knfn.NewFunction(directory)
		if err != nil {
			panic(err)
		}

		err = cfg.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncEnvs(envs map[string]string) SetupOpts {
	return func(directory string) {
		cfg, err := knfn.NewFunction(directory)
		if err != nil {
			panic(err)
		}

		for envName, envValue := range envs {
			name := envName
			value := envValue
			cfg.Envs = append(cfg.Envs, knfn.Env{
				Name:  &name,
				Value: &value,
			})
		}

		err = cfg.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncScale(scale knfn.ScaleOptions) SetupOpts {
	return func(directory string) {
		cfg, err := knfn.NewFunction(directory)
		if err != nil {
			panic(err)
		}

		cfg.Options.Scale = &scale

		err = cfg.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncResourceRequests(requests knfn.ResourcesRequestsOptions) SetupOpts {
	return func(directory string) {
		cfg, err := knfn.NewFunction(directory)
		if err != nil {
			panic(err)
		}

		if cfg.Options.Resources == nil {
			cfg.Options.Resources = &knfn.ResourcesOptions{}
		}

		cfg.Options.Resources.Requests = &requests

		err = cfg.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncResourceLimits(limits knfn.ResourcesLimitsOptions) SetupOpts {
	return func(directory string) {
		cfg, err := knfn.NewFunction(directory)
		if err != nil {
			panic(err)
		}

		if cfg.Options.Resources == nil {
			cfg.Options.Resources = &knfn.ResourcesOptions{}
		}

		cfg.Options.Resources.Limits = &limits

		err = cfg.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}

func WithFunctionFile(module, function, funcTemplate string) SetupOpts {
	return func(directory string) {
		f := HTTPFunction{
			Module: module,
			Name:   function,
		}

		err := f.Generate(directory, funcTemplate)
		if err != nil {
			panic(err)
		}
	}
}
