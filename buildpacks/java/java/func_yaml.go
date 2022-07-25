// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package java

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/libpak/bard"
	"gopkg.in/yaml.v3"
	knfn "knative.dev/kn-plugin-func"
)

type FuncYaml struct {
	Name    string
	Options map[string]string
	Envs    map[string]string
	Exists  bool
}

func ParseFuncYaml(filedir string, logger bard.Logger) FuncYaml {
	file := filepath.Join(filedir, knfn.ConfigFile)
	_, err := os.Stat(file)
	if err != nil {
		logger.Bodyf("unable to find file '%s'", file)
		return FuncYaml{}
	}

	cfg, err := knfn.NewFunction(filedir)
	if err != nil {
		logger.Bodyf("unable to parse '%s': %v", knfn.ConfigFile, err)
		return FuncYaml{}
	}

	options := optionsToMap(cfg.Options, logger)
	envs := envsToMap(cfg.Envs, logger)
	name := getName(cfg.Envs, cfg.Name)
	return FuncYaml{
		Name:    name,
		Options: options,
		Envs:    envs,
		Exists:  true,
	}
}

func getName(envs knfn.Envs, nameFromYaml string) string {
	result := nameFromYaml

	if envs == nil {
		return result
	}

	for _, e := range envs {
		key := *e.Name
		if key == "FUNCTION_NAME" {
			result = *e.Value
			break
		}
	}

	return result
}

func envsToMap(envs knfn.Envs, logger bard.Logger) map[string]string {
	result := map[string]string{}
	if envs == nil {
		return result
	}

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

func optionsToMap(options knfn.Options, logger bard.Logger) map[string]string {
	result := map[string]string{}

	if options.Scale != nil {
		scaleJson, err := yaml.Marshal(options.Scale)
		if err != nil {
			logger.Bodyf("unable to marshal func.yaml options.Scale")
		} else {
			result["options-scale"] = string(scaleJson)
		}
	}

	if options.Resources != nil {
		if options.Resources.Requests != nil {
			requestsJson, err := yaml.Marshal(options.Resources.Requests)
			if err != nil {
				logger.Bodyf("unable to marshal func.yaml options.Resources.Requests")
			} else {
				result["options-resources-requests"] = string(requestsJson)
			}
		}

		if options.Resources.Limits != nil {
			limitsJson, err := yaml.Marshal(options.Resources.Limits)
			if err != nil {
				logger.Bodyf("unable to marshal func.yaml options.Resources.Limits")
			} else {
				result["options-resources-limits"] = string(limitsJson)
			}
		}
	}

	return result
}
