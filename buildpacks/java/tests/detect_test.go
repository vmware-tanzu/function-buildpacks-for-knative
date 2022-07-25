// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"bytes"
	"kn-fn/java-function-buildpack/java"
	"reflect"
	"testing"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
)

type result struct {
	plan libcnb.DetectResult
	err  error
}

func TestDetectNoEnvironmentWithValidFile(t *testing.T) {
	d := getDetector()
	appDir, cleanup := SetupTestDirectory(
		WithFuncName("my-function"),
	)
	defer cleanup()

	plan, err := d.Detect(getContext(
		withApplicationPath(appDir),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"func_yaml_function_name": "my-function",
			"func_yaml_envs":          map[string]string{},
			"func_yaml_options":       map[string]string{},
		},
	}
	expectations.Check(t)
}

func TestDetectEnvironmentWithValidFile(t *testing.T) {
	d := getDetector()
	appDir, cleanup := SetupTestDirectory(
		WithFuncEnvs(map[string]string{
			"MODULE_NAME":   "other",
			"FUNCTION_NAME": "handler2",
		}),
	)
	defer cleanup()

	plan, err := d.Detect(getContext(
		withApplicationPath(appDir),
		withModuleName("other"),
		withFunctionName("handler2"),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"func_yaml_function_name": "handler2",
			"func_yaml_envs": map[string]string{
				"MODULE_NAME":   "other",
				"FUNCTION_NAME": "handler2",
			},
			"func_yaml_options": map[string]string{},
		},
	}
	expectations.Check(t)
}

func TestDetectNameFromEnvsBeforeFuncYamlFile(t *testing.T) {
	d := getDetector()
	appDir, cleanup := SetupTestDirectory(
		WithFuncName("not-the-best-name"),
		WithFuncEnvs(map[string]string{
			"MODULE_NAME":   "other",
			"FUNCTION_NAME": "goodName",
		}),
	)
	defer cleanup()

	plan, err := d.Detect(getContext(
		withApplicationPath(appDir),
		withModuleName("other"),
		withFunctionName("goodName"),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"func_yaml_function_name": "goodName",
			"func_yaml_envs": map[string]string{
				"MODULE_NAME":   "other",
				"FUNCTION_NAME": "goodName",
			},
			"func_yaml_options": map[string]string{},
		},
	}
	expectations.Check(t)
}

func getDetector() java.Detect {
	buf := bytes.NewBuffer(nil)
	logger := bard.NewLogger(buf)
	return java.Detect{
		Logger: logger,
	}
}

func getContext(opts ...func(*libcnb.DetectContext)) libcnb.DetectContext {
	ctx := libcnb.DetectContext{
		Application: libcnb.Application{},
		Buildpack:   libcnb.Buildpack{},
		Platform: libcnb.Platform{
			Environment: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(&ctx)
	}

	return ctx
}

func withApplicationPath(path string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Application.Path = path
	}
}

func withModuleName(value string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Platform.Environment["MODULE_NAME"] = value
	}
}

func withFunctionName(value string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Platform.Environment["FUNCTION_NAME"] = value
	}
}

func withEnvironmentVariable(key string, value string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Platform.Environment[key] = value
	}
}

type DetectExpectations struct {
	Result result

	ExpectError  bool
	Err          error
	Pass         bool
	Metadata     map[string]interface{}
	SkipProvides bool
	SkipRequires bool
}

func (e DetectExpectations) Check(t *testing.T) {
	if !e.ExpectError && e.Result.err != nil {
		if e.Err != nil && e.Result.err != e.Err {
			t.Errorf("expected error %v, but received error %v", e.Err, e.Result.err)
		} else {
			t.Errorf("unexpected error received: %v", e.Result.err)
		}
		return
	}

	if e.Pass != e.Result.plan.Pass {
		t.Errorf("expected detection to pass but it failed")
		return
	}

	// Find the invoker requires
	planName := "java-function"
	for _, plan := range e.Result.plan.Plans {
		for _, require := range plan.Requires {
			if require.Name == planName {
				for k, v := range e.Metadata {
					val, ok := require.Metadata[k]
					if !ok {
						t.Errorf("plan requirement '%s' did not find metadata with key %s", planName, k)
						return
					}

					if !reflect.DeepEqual(v, val) {
						t.Errorf("unexpected value in '%s' requires metadata for key %s. Expected %s but got %s", planName, k, v, val)
						return
					}
				}
			}
		}
	}
}
