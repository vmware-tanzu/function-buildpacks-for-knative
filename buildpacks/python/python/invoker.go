package python

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Invoker struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger

	modulePath string

	handler  string
	module   string
	function string
}

func NewInvokerFromPlan(plan libcnb.BuildpackPlanEntry, buildpackPath string) Invoker {
	contributor := libpak.NewLayerContributor(plan.Name, map[string]interface{}{}, libcnb.LayerTypes{
		Launch: true,
	})

	// Assumption is that build always comes after a successful detection which will add the handler key
	handler := plan.Metadata["handler"].(map[string]interface{})

	return Invoker{
		LayerContributor: contributor,
		modulePath:       filepath.Join(buildpackPath, "invoker"),

		module:   handler["module"].(string),
		function: handler["function"].(string),
		handler:  handler["raw"].(string),
	}
}

func (i Invoker) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger

	return i.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		i.Logger.Bodyf("Creating layer %s with path %s", layer.Name, layer.Path)

		// TODO: Maybe do the copying in GO instead of shelling out to cp.
		pyfuncDir := filepath.Join(layer.Path, "pyfunc")
		cp := exec.Command("cp", "-a", i.modulePath, pyfuncDir)
		if err := cp.Run(); err != nil {
			return layer, err
		}

		var stderr bytes.Buffer
		cmd := exec.Command("python", "setup.py", "install")
		cmd.Dir = pyfuncDir
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			i.Logger.Body("Installation of pyfunc-invoker failed:\n%s", stderr.String())
			return layer, err
		}

		if i.handler != "" {
			layer.LaunchEnvironment.Override("PYTHON_HANDLER", i.handler)
		}
		return layer, nil
	})
}

func (i Invoker) Name() string {
	return i.LayerContributor.Name
}
