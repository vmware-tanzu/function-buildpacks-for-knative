package java

import (
	"fmt"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Function struct {
	ApplicationPath  string
	Handler          string
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger
}

func NewFunction(applicationPath string, handler string) (Function, error) {
	return Function{
		ApplicationPath: applicationPath,
		Handler:         handler,
		LayerContributor: libpak.NewLayerContributor(
			fmt.Sprintf("%s %s", "Java", handler),
			map[string]interface{}{"handler": handler},
			libcnb.LayerTypes{
				Launch: true,
			},
		),
	}, nil
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

		return layer, nil
	})
}

func (Function) Name() string {
	return "function"
}