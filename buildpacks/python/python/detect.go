// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"

	"kn-fn/buildpacks/config"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{}

	funcYaml := config.ParseFuncYaml(context.Application.Path, d.Logger)

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
				Metadata: map[string]any{
					"launch":            true,
					"func_yaml_envs":    funcYaml.Envs,
					"func_yaml_options": funcYaml.Options,
				},
			},
			{
				Name: "site-packages",
				Metadata: map[string]any{
					"build":  true,
					"launch": true,
				},
			},
			{
				Name: "pip",
				Metadata: map[string]any{
					"build": true,
				},
			},
			{
				Name: "cpython",
				Metadata: map[string]any{
					"build":  true,
					"launch": true,
				},
			},
		},
	})

	result.Pass = funcYaml.Exists || funcFound

	return result, nil
}
