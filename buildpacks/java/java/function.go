package java

import (
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Function struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger

	ApplicationPath string
	Handler         string
	Envs            map[string]interface{}
}

func NewFunction(plan libcnb.BuildpackPlanEntry, applicationPath string) Function {
	envs := plan.Metadata["envs"].(map[string]interface{})

	return Function{
		ApplicationPath: applicationPath,
		LayerContributor: libpak.NewLayerContributor(
			plan.Name,
			envs,
			libcnb.LayerTypes{
				Launch: true,
			},
		),
		Envs: envs,
	}
}

func (f Function) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	f.LayerContributor.Logger = f.Logger

	return f.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		if len(f.Handler) > 0 {
			if strings.ContainsAny(f.Handler, ".") {
				layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_FUNCTION_CLASS", f.Handler)
			} else {
				layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_DEFINITION", f.Handler)
			}
		}

		layer.LaunchEnvironment.Default("SPRING_CLOUD_FUNCTION_LOCATION", f.ApplicationPath)

		for k, v := range f.Envs {
			layer.LaunchEnvironment.Default(k, v)
		}

		return layer, nil
	})
}

func (f Function) Name() string {
	return f.LayerContributor.Name
}
