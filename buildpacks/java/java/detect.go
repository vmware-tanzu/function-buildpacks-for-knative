// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"

	"kn-fn/buildpacks/config"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) checkConfigs(cr libpak.ConfigurationResolver) bool {
	if val, defined := cr.Resolve("BP_FUNCTION"); defined && val != "" {
		return true
	}

	return false
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{}

	funcYaml := config.ParseFuncYaml(context.Application.Path, d.Logger)

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
				Metadata: map[string]any{
					"launch":            true,
					"func_yaml_envs":    funcYaml.Envs,
					"func_yaml_options": funcYaml.Options,
				},
			},
			{
				Name: "jre",
				Metadata: map[string]any{
					"launch": true,
				},
			},
			{
				Name: "jvm-application-package",
			},
		},
	})

	result.Pass = funcYaml.Exists || configPass
	return result, nil
}
