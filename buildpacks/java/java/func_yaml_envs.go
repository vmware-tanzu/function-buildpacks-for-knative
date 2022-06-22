// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

type FuncYamlEnvs struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger

	Envs map[string]string
}

func NewFuncYamlEnvs(applicationPath string) FuncYamlEnvs {
	envs := getFuncYamlEnvs(applicationPath)
	return FuncYamlEnvs{
		LayerContributor: libpak.NewLayerContributor(
			"func-yaml-envs",
			envs,
			libcnb.LayerTypes{
				Launch: true,
			},
		),
		Envs: envs,
	}
}

func (f FuncYamlEnvs) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.LayerContributor.Logger = f.Logger
	return f.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		for k, v := range f.Envs {
			layer.LaunchEnvironment.Default(k, v)
		}
		return layer, nil
	})
}

func (f FuncYamlEnvs) Name() string {
	return f.LayerContributor.Name
}

func getFuncYamlEnvs(applicationPath string) map[string]string {
	envs := map[string]string{}

	configFile := filepath.Join(applicationPath, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		return envs
	}

	f, err := knfn.NewFunction(applicationPath)
	if err != nil {
		return envs
	}

	return envsToMap(f.Envs)
}

func envsToMap(envs knfn.Envs) map[string]string {
	result := map[string]string{}

	for _, e := range envs {
		key := *e.Name
		val := ""
		if e.Value != nil {
			val = *e.Value
		}
		result[key] = val
	}

	return result
}
