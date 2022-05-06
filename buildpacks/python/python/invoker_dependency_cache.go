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

type InvokerDependencyCache struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger

	CacheDir string
}

func NewInvokerDependencyCache(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (InvokerDependencyCache, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return InvokerDependencyCache{LayerContributor: contributor}, entry
}

func (i InvokerDependencyCache) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger
	i.CacheDir = layer.Path

	return i.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		i.Logger.Bodyf("Installing to %s", artifact.Name())

		if err := crush.ExtractTarGz(artifact, layer.Path, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		return layer, nil
	})
}

func (i InvokerDependencyCache) Name() string {
	return i.LayerContributor.Name()
}
