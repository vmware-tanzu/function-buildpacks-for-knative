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
}

func NewFunction(plan libcnb.BuildpackPlanEntry) Function {
	contributor := libpak.NewLayerContributor(plan.Name, map[string]interface{}{}, libcnb.LayerTypes{
		Launch: true,
	})

	// Assumption is that build always comes after a successful detection which will add the appropriate envs
	handler := plan.Metadata["envs"].(map[string]interface{})

	return Function{
		LayerContributor: contributor,
		module:           handler[EnvModuleName].(string),
		function:         handler[EnvFunctionName].(string),
	}
}

func (i Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger

	return i.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		layer.LaunchEnvironment.Override(EnvModuleName, i.module)
		layer.LaunchEnvironment.Override(EnvFunctionName, i.function)
		return layer, nil
	})
}

func (i Function) Name() string {
	return i.LayerContributor.Name
}
