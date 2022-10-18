// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

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

func createYaml(dir string) knfn.Function {
	cfg, err := knfn.NewFunction(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	} else if err != nil && errors.Is(err, os.ErrNotExist) {
		cfg = knfn.Function{
			Name:    "test",
			Runtime: "unknown",
			Root:    dir,
		}
	}

	err = cfg.Write()
	if err != nil {
		panic(err)
	}
	return cfg
}

func WithFuncYaml() SetupOpts {
	return func(directory string) {
		createYaml(directory)
	}
}

func WithFuncEnvs(envs map[string]string) SetupOpts {
	return func(directory string) {
		cfg := createYaml(directory)

		for envName, envValue := range envs {
			name := envName
			value := envValue
			cfg.Run.Envs = append(cfg.Run.Envs, knfn.Env{
				Name:  &name,
				Value: &value,
			})
		}

		err := cfg.Write()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncScale(scale knfn.ScaleOptions) SetupOpts {
	return func(directory string) {
		cfg := createYaml(directory)

		cfg.Deploy.Options.Scale = &scale

		err := cfg.Write()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncResourceRequests(requests knfn.ResourcesRequestsOptions) SetupOpts {
	return func(directory string) {
		cfg := createYaml(directory)

		if cfg.Deploy.Options.Resources == nil {
			cfg.Deploy.Options.Resources = &knfn.ResourcesOptions{}
		}

		cfg.Deploy.Options.Resources.Requests = &requests

		err := cfg.Write()
		if err != nil {
			panic(err)
		}
	}
}

func WithFuncResourceLimits(limits knfn.ResourcesLimitsOptions) SetupOpts {
	return func(directory string) {
		cfg := createYaml(directory)

		if cfg.Deploy.Options.Resources == nil {
			cfg.Deploy.Options.Resources = &knfn.ResourcesOptions{}
		}

		cfg.Deploy.Options.Resources.Limits = &limits

		err := cfg.Write()
		if err != nil {
			panic(err)
		}
	}
}
