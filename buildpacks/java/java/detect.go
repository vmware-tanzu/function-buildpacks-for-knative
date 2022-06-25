// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	function "knative.dev/kn-plugin-func"
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
	var labels []libcnb.Label
	result := libcnb.DetectResult{}

	appPath := context.Application.Path
	funcYamlPass := d.checkFuncYaml(appPath)

	if funcYamlPass {
		labels = d.getFuncYamlOptions(appPath)
	}

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
					"labels":        labels,
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

func (d Detect) getFuncYamlOptions(appPath string) []libcnb.Label {
	configFile := filepath.Join(appPath, knfn.ConfigFile)
	_, err := os.Stat(configFile)
	if err != nil {
		d.Logger.Bodyf("'%s' not detected", knfn.ConfigFile)
		return []libcnb.Label{}
	}

	f, err := knfn.NewFunction(appPath)
	if err != nil {
		d.Logger.Bodyf("unable to parse '%s': %v", knfn.ConfigFile, err)
		return []libcnb.Label{}
	}
	labels := d.optionsToLabels(f.Options)
	for _, l := range labels {
		f.Labels = append(f.Labels, function.Label{
			Key:   &l.Key,
			Value: &l.Value,
		})
	}
	return labels
}

func (d Detect) optionsToLabels(options knfn.Options) []libcnb.Label {
	labels := []libcnb.Label{}

	scaleJson, err := json.Marshal(options.Scale)
	if err != nil {
		d.Logger.Bodyf("unable to marshal func.yaml options.Scale")
	}
	requestsJson, err := json.Marshal(options.Resources.Requests)
	if err != nil {
		d.Logger.Bodyf("unable to marshal func.yaml options.Resources.Requests")

	}
	limitsJson, err := json.Marshal(options.Resources.Limits)
	if err != nil {
		d.Logger.Bodyf("unable to marshal func.yaml options.Resources.Limits")
	}
	labels = append(labels,
		libcnb.Label{
			Key:   "options-scale",
			Value: string(scaleJson),
		},
		libcnb.Label{
			Key:   "options-resources-requests",
			Value: string(requestsJson),
		},
		libcnb.Label{
			Key:   "options-resources-limits",
			Value: string(limitsJson),
		})

	return labels
}
