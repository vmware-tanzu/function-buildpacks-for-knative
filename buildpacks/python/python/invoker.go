package python

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Invoker struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewInvoker(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (Invoker, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return Invoker{LayerContributor: contributor}, entry
}

func (i Invoker) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger
	return i.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		i.Logger.Bodyf("Installing to %s", artifact.Name())

		var stderr bytes.Buffer
		cmd := exec.Command("pip", "install", artifact.Name())
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			i.Logger.Body("failed to install invoker: %s", stderr.String())
			return layer, err
		}

		return layer, nil
	})
}

func (i Invoker) Name() string {
	return i.LayerContributor.Name()
}
