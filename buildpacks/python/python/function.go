// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Function struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger

	module   string
	function string
	envs     map[string]interface{}
}

func NewFunction(plan libcnb.BuildpackPlanEntry) Function {
	envs := plan.Metadata["envs"].(map[string]interface{})
	contributor := libpak.NewLayerContributor(plan.Name, envs, libcnb.LayerTypes{
		Launch: true,
	})
	return Function{
		LayerContributor: contributor,
		envs:             envs,
	}
}

func (f Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.LayerContributor.Logger = f.Logger

	return f.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		for k, v := range f.envs {
			layer.LaunchEnvironment.Default(k, v)
		}
		return layer, nil
	})
}

func (i Function) Name() string {
	return i.LayerContributor.Name
}
