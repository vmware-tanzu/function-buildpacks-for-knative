// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Function struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger

	ApplicationPath string
	Function        string
	DefaultFunction string
}

func NewFunction(function string, defaultFunc string, applicationPath string) Function {
	return Function{
		ApplicationPath: applicationPath,
		LayerContributor: libpak.NewLayerContributor(
			"java-function",
			map[string]string{
				"function": function,
			},
			libcnb.LayerTypes{
				Launch: true,
			},
		),
		Function:        function,
		DefaultFunction: defaultFunc,
	}
}

func (f Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.LayerContributor.Logger = f.Logger

	return f.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_LOCATION", f.ApplicationPath)
		layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", f.Function) // Function lives here

		if len(f.DefaultFunction) > 0 {
			layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_DEFINITION", f.DefaultFunction)
		}

		return layer, nil
	})
}

func (f Function) Name() string {
	return f.LayerContributor.Name
}
