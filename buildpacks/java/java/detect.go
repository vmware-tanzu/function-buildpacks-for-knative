package java

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{}

	configFile := filepath.Join(context.Application.Path, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		d.logf(fmt.Sprintf("unable to find file '%s'", configFile))
		return result, nil
	}

	f, err := knfn.NewFunction(context.Application.Path)
	if err != nil {
		return result, fmt.Errorf("parsing function config: %v", err)
	}

	envs := envsToMap(f.Envs)

	result.Plans = append(result.Plans, libcnb.BuildPlan{
		Provides: []libcnb.BuildPlanProvide{
			{
				Name: "java-function",
			},
		},
		Requires: []libcnb.BuildPlanRequire{
			{
				Name: "java-function",
				Metadata: map[string]interface{}{
					"launch": true,
					"envs":   envs,
				},
			},
			{
				Name: "jre",
				Metadata: map[string]interface{}{
					"launch": true,
				},
			},
			{
				Name: "jvm-application",
			},
		},
	})

	result.Pass = true
	return result, nil
}

func (d Detect) logf(format string, args ...interface{}) {
	d.Logger.Infof(format, args...)
}

func envsToMap(envs knfn.Envs) map[string]string {
	result := map[string]string{}

	for _, e := range envs {
		key := *e.Name
		val := ""
		if e.Value != nil {
			val = *e.Value
		}
		result[key] = val
	}

	return result
}
