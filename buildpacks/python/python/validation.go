package python

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type FunctionValidationLayer struct {
	LayerContributor libpak.LayerContributor
	Logger           bard.Logger

	module          string
	function        string
	applicationPath string
}

func NewFunctionValidationLayer(plan libcnb.BuildpackPlanEntry, appPath string) FunctionValidationLayer {
	contributor := libpak.NewLayerContributor("validation", map[string]interface{}{}, libcnb.LayerTypes{})
	
	handler := plan.Metadata["envs"].(map[string]interface{})

	return FunctionValidationLayer{
		LayerContributor: contributor,
		module:           handler[EnvModuleName].(string),
		function:         handler[EnvFunctionName].(string),
		applicationPath:  appPath,
	}
}

func (i FunctionValidationLayer) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	i.LayerContributor.Logger = i.Logger

	return i.LayerContributor.Contribute(layer, func() (libcnb.Layer, error) {
		i.Logger.Body("Validating function")

		var stderr bytes.Buffer
		cmd := exec.Command("python", "-m", "pyfunc", "check", "-s", i.applicationPath)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", EnvModuleName, i.module), fmt.Sprintf("%s=%s", EnvFunctionName, i.function))
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return layer, fmt.Errorf("%v: %v", stderr.String(), err)
		}

		i.Logger.Body("Function was successfully parsed")
		return layer, nil
	})
}

func (i FunctionValidationLayer) Name() string {
	return i.LayerContributor.Name
}
