// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"kn-fn/java-function-buildpack/java"
	"testing"

	"gopkg.in/yaml.v3"
	knfn "knative.dev/kn-plugin-func"
	"knative.dev/pkg/ptr"
)

func TestParseFuncYaml_FileDoesNotExist(t *testing.T) {
	appDir, cleanup := SetupTestDirectory()
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())
	if result.Exists {
		t.Logf("File should not exists but was detected")
		t.Fail()
	}
}

func TestParseFuncYaml_FileExistsButEmpty(t *testing.T) {
	appDir, cleanup := SetupTestDirectory(WithFuncYaml())
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())
	if !result.Exists {
		t.Logf("File should exists but was not detected")
		t.Fail()
	}
}

func TestParseFuncYaml_HasFuncClass(t *testing.T) {
	expectedFunctionName := "functionName"
	appDir, cleanup := SetupTestDirectory(WithFuncName(expectedFunctionName))
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())

	if result.Name != "functionName" {
		t.Logf("Expected function name to be %s but received %s", expectedFunctionName, result.Name)
		t.Fail()
	}
}

func TestParseFuncYaml_HasEnvs(t *testing.T) {
	envs := map[string]string{
		"my-env":   "my-env-value",
		"your-env": "your-env-value",
	}

	appDir, cleanup := SetupTestDirectory(WithFuncEnvs(envs))
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())

	for k, v := range result.Envs {
		expected, found := envs[k]
		if !found {
			t.Logf("Key %s was not in expected results", k)
			t.Fail()
		}

		if v != expected {
			t.Logf("Expected value %s but received %s", expected, v)
			t.Fail()
		}
	}
}

func TestParseFuncYaml_HasScale(t *testing.T) {
	scaleOption := knfn.ScaleOptions{
		Min:         ptr.Int64(4),
		Max:         ptr.Int64(9),
		Metric:      ptr.String("rps"),
		Target:      ptr.Float64(0.5),
		Utilization: ptr.Float64(50),
	}

	appDir, cleanup := SetupTestDirectory(WithFuncScale(scaleOption))
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())
	resultScaleOptions := &knfn.ScaleOptions{}
	yaml.Unmarshal([]byte(result.Options["options-scale"]), resultScaleOptions)

	if *resultScaleOptions.Min != *scaleOption.Min {
		t.Logf("Expected Min %d but received %d", *scaleOption.Min, *resultScaleOptions.Min)
		t.Fail()
	}
	if *resultScaleOptions.Max != *scaleOption.Max {
		t.Logf("Expected Max %d but received %d", *scaleOption.Max, *resultScaleOptions.Max)
		t.Fail()
	}
	if *resultScaleOptions.Metric != *scaleOption.Metric {
		t.Logf("Expected Metric %s but received %s", *scaleOption.Metric, *resultScaleOptions.Metric)
		t.Fail()
	}
	if *resultScaleOptions.Target != *scaleOption.Target {
		t.Logf("Expected Target %d but received %d", scaleOption.Target, resultScaleOptions.Target)
		t.Fail()
	}
	if *resultScaleOptions.Utilization != *scaleOption.Utilization {
		t.Logf("Expected Utilization %d but received %d", scaleOption.Utilization, resultScaleOptions.Utilization)
		t.Fail()
	}
}

func TestParseFuncYaml_HasRequests(t *testing.T) {
	requestOptions := knfn.ResourcesRequestsOptions{
		CPU:    ptr.String("1"),
		Memory: ptr.String("40m"),
	}

	appDir, cleanup := SetupTestDirectory(WithFuncResourceRequests(requestOptions))
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())
	resultRequestOptions := &knfn.ResourcesRequestsOptions{}
	yaml.Unmarshal([]byte(result.Options["options-resources-requests"]), resultRequestOptions)

	if *resultRequestOptions.CPU != *requestOptions.CPU {
		t.Logf("Expected CPU %s but received %s", *requestOptions.CPU, *resultRequestOptions.CPU)
		t.Fail()
	}
	if *resultRequestOptions.Memory != *requestOptions.Memory {
		t.Logf("Expected Memory %s but received %s", *requestOptions.Memory, *resultRequestOptions.Memory)
		t.Fail()
	}
}

func TestParseFuncYaml_HasLimits(t *testing.T) {
	limitOptions := knfn.ResourcesLimitsOptions{
		CPU:         ptr.String("1"),
		Memory:      ptr.String("40m"),
		Concurrency: ptr.Int64(5),
	}

	appDir, cleanup := SetupTestDirectory(WithFuncResourceLimits(limitOptions))
	defer cleanup()
	result := java.ParseFuncYaml(appDir, NewLogger())
	resultLimitOptions := &knfn.ResourcesLimitsOptions{}
	yaml.Unmarshal([]byte(result.Options["options-resources-limits"]), resultLimitOptions)

	if *resultLimitOptions.CPU != *limitOptions.CPU {
		t.Logf("Expected CPU %s but received %s", *limitOptions.CPU, *resultLimitOptions.CPU)
		t.Fail()
	}
	if *resultLimitOptions.Memory != *limitOptions.Memory {
		t.Logf("Expected Memory %s but received %s", *limitOptions.Memory, *resultLimitOptions.Memory)
		t.Fail()
	}
	if *resultLimitOptions.Concurrency != *limitOptions.Concurrency {
		t.Logf("Expected Memory %s but received %s", *limitOptions.Memory, *resultLimitOptions.Memory)
		t.Fail()
	}
}
