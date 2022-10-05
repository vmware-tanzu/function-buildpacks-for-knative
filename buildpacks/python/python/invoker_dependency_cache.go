// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"

	"kn-fn/buildpacks/command"
)

type InvokerDependencyCache struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger

	pythonPath    string
	commandRunner command.Runner
}

func NewInvokerDependencyCache(
	dependency libpak.BuildpackDependency,
	cache libpak.DependencyCache,
	commandRunner command.Runner,
) *InvokerDependencyCache {
	contributor := libpak.NewDependencyLayerContributor(dependency, cache, libcnb.LayerTypes{
		Launch: true,
		Cache:  true,
	})

	return &InvokerDependencyCache{
		LayerContributor: contributor,
		commandRunner:    commandRunner,
	}
}

func (i *InvokerDependencyCache) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger
	i.pythonPath = filepath.Join(layer.Path, "install")

	return i.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		i.Logger.Bodyf("Extracting dependency to %s", layer.Path)

		depsDir, err := os.MkdirTemp("", "")
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to create temp directory\n%w", err)
		}
		if err := crush.Extract(artifact, depsDir, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to extract %s\n%w", artifact.Name(), err)
		}

		if err := os.Mkdir(i.PythonPath(), 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to make dependency directory %s\n%w", i.PythonPath(), err)
		}

		args := []string{
			"install",
			"--target", i.PythonPath(),
			"--no-index",
			"--find-links", depsDir,
			"--compile",
			"--disable-pip-version-check",
			"--ignore-installed",
			"--exists-action=w",
		}

		files, err := filepath.Glob(filepath.Join(depsDir, "*"))
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to glob for dependencies: %w", err)
		}

		args = append(args, files...)

		cmd := exec.Command("pip", args...)

		if output, err := i.commandRunner.Run(cmd); err != nil {
			i.Logger.Body("failed to install dependencies: %s", output)
			return layer, err
		}

		layer.LaunchEnvironment.Append("PYTHONPATH", string(os.PathListSeparator), i.PythonPath())

		return layer, nil
	})
}

func (i *InvokerDependencyCache) Name() string {
	return "invoker-deps"
}

func (i *InvokerDependencyCache) PythonPath() string {
	return i.pythonPath
}
