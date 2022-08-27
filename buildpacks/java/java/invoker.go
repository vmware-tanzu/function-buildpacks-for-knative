// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

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
}

func NewInvoker(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (Invoker, libcnb.BOMEntry) {
	dependency.CPEs = []string{}
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return Invoker{LayerContributor: contributor}, entry
}

func (i Invoker) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger

	return i.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		i.Logger.Bodyf("Expanding to %s", layer.Path)

		if err := crush.ExtractZip(artifact, layer.Path, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		layer.LaunchEnvironment.Prepend("CLASSPATH", string(os.PathListSeparator), layer.Path)

		return layer, nil
	})
}

func (i Invoker) Name() string {
	return i.LayerContributor.LayerName()
}
