// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
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

	functionClass, functionClassSet := cr.Resolve("BP_FUNCTION")
	if _, _, err := parseFunctionClass(functionClass); err != nil {
		d.Logger.Body(err.Error())
		return libcnb.DetectResult{}, nil
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

	result.Pass = funcYaml.Exists || functionClassSet

	return result, nil
}
