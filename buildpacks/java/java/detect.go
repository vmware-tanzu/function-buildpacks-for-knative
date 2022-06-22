// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	knfn "knative.dev/kn-plugin-func"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) checkConfigs(cr libpak.ConfigurationResolver) bool {
	if _, defined := cr.Resolve("BP_FUNCTION"); defined {
		return true
	}

	if _, defined := cr.Resolve("BP_DEFAULT_FUNCTION"); defined {
		return true
	}

	return false
}

func (d Detect) checkFuncYaml(appPath string) bool {
	configFile := filepath.Join(appPath, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		d.Logger.Bodyf("unable to find file '%s'", configFile)
		return false
	}

	return true
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{}

	appPath := context.Application.Path
	funcYamlPass := d.checkFuncYaml(appPath)

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &d.Logger)
	if err != nil {
		return result, fmt.Errorf("unable to create configuration resolver: %v", err)
	}

	configPass := d.checkConfigs(cr)
	if err != nil {
		d.Logger.Bodyf("unable to check buildpack configurations: %v", err)
		return result, nil
	}

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
					"launch":        true,
					"has_func_yaml": funcYamlPass,
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

	result.Pass = funcYamlPass || configPass
	return result, nil
}
