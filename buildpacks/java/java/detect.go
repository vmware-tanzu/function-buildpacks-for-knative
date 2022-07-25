// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
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

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{}

	funcYaml := ParseFuncYaml(context.Application.Path, d.Logger)

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
					"launch":                  true,
					"func_yaml_envs":          funcYaml.Envs,
					"func_yaml_function_name": funcYaml.Name,
					"func_yaml_options":       funcYaml.Options,
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

	result.Pass = funcYaml.Exists || configPass
	return result, nil
}
