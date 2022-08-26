// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
)

type Invoker struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger

	pythonPath string
}

func NewInvoker(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) Invoker {
	contributor := libpak.NewDependencyLayerContributor(dependency, cache, libcnb.LayerTypes{
		Launch: true,
		Cache:  true,
	})
	return Invoker{LayerContributor: contributor}
}

func (i *Invoker) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger
	i.pythonPath = layer.Path

	return i.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		i.Logger.Bodyf("Extracting invoker to %s", layer.Path)

		if err := crush.Extract(artifact, layer.Path, 1); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		layer.LaunchEnvironment.Append("PYTHONPATH", string(os.PathListSeparator), i.PythonPath())

		return layer, nil
	})
}

func (i *Invoker) Name() string {
	return "invoker"
}

func (i *Invoker) PythonPath() string {
	return i.pythonPath
}
