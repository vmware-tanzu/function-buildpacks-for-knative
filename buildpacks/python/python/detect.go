// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

const (
	EnvModuleName   = "MODULE_NAME"
	EnvFunctionName = "FUNCTION_NAME"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) getFuncYamlEnvs(appPath string) (map[string]string, bool) {
	configFile := filepath.Join(appPath, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		d.Logger.Bodyf("'%s' not detected", knfn.ConfigFile)
		return make(map[string]string), false
	}

	f, err := knfn.NewFunction(appPath)
	if err != nil {
		d.Logger.Bodyf("unable to parse '%s': %v", knfn.ConfigFile, err)
		return make(map[string]string), false
	}
	return envsToMap(f.Envs), true
}

func (d Detect) getFuncYamlOptions(appPath string) (map[string]string, bool) {
	configFile := filepath.Join(appPath, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		d.Logger.Bodyf("'%s' not detected", knfn.ConfigFile)
		return make(map[string]string), false
	}

	f, err := knfn.NewFunction(appPath)
	if err != nil {
		d.Logger.Bodyf("unable to parse '%s': %v", knfn.ConfigFile, err)
		return make(map[string]string), false
	}
	return optionsToMap(f.Options), true
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{}

	envs, hasValidFuncYaml := d.getFuncYamlEnvs(context.Application.Path)

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &d.Logger)
	if err != nil {
		return libcnb.DetectResult{}, err
	}

	functionHandler, funcFound := cr.Resolve("BP_FUNCTION")
	funcParts := strings.Split(functionHandler, ".")
	if funcFound {
		if len(funcParts) != 2 || len(funcParts[0]) == 0 || len(funcParts[1]) == 0 { // We're expecting the format of module.func_name
			d.Logger.Bodyf("BP_FUNCTION detected but is invalid, it should be in the form of `module.function_name`")
			return libcnb.DetectResult{}, nil
		}

		envs[EnvModuleName] = funcParts[0]
		envs[EnvFunctionName] = funcParts[1]
	} else {
		if _, found := envs[EnvModuleName]; !found {
			envs[EnvModuleName] = funcParts[0]
		}

		if _, found := envs[EnvFunctionName]; !found {
			envs[EnvFunctionName] = funcParts[1]
		}
	}

	result.Plans = append(result.Plans, libcnb.BuildPlan{
		Provides: []libcnb.BuildPlanProvide{
			{
				Name: "python-function",
			},
		},
		Requires: []libcnb.BuildPlanRequire{
			{
				Name: "python-function",
				Metadata: map[string]interface{}{
					"launch": true,
					"envs":   envs,
				},
			},
			{
				Name: "site-packages",
				Metadata: map[string]interface{}{
					"build":  true,
					"launch": true,
				},
			},
			{
				Name: "pip",
				Metadata: map[string]interface{}{
					"build": true,
				},
			},
		},
	})

	result.Pass = hasValidFuncYaml || funcFound

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
