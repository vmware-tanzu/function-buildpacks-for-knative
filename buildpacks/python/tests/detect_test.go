// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"bytes"
	"kn-fn/python-function-buildpack/python"
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
	app := createTestApplication(withDefaultHTTPFunction())
	defer app.Finish()

	plan, err := d.Detect(getContext(
		withApplicationPath(app.ApplicationPath),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"envs": map[string]string{
				EnvModuleName:   "func",
				EnvFunctionName: "main",
			},
		},
	}
	expectations.Check(t)
}

func TestDetectEnvironmentWithValidFile(t *testing.T) {
	d := getDetector()
	app := createTestApplication(withHTTPFunction("other", "handler2"))
	defer app.Finish()

	plan, err := d.Detect(getContext(
		withApplicationPath(app.ApplicationPath),
		withModuleName("other"),
		withFunctionName("handler2"),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"envs": map[string]string{
				EnvModuleName:   "other",
				EnvFunctionName: "handler2",
			},
		},
	}
	expectations.Check(t)
}

func TestDetectNoEnvironmentWithResources(t *testing.T) {
	d := getDetector()
	app := createTestApplication(withDefaultHTTPFunction())
	defer app.Finish()

	plan, err := d.Detect(getContext(
		withApplicationPath(app.ApplicationPath),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"envs": map[string]string{
				EnvModuleName:   "func",
				EnvFunctionName: "main",
			},
			"options": map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]string{
						"cpu":    "100m",
						"memory": "128Mi",
					},
					"limits": map[string]string{
						"cpu":         "1000m",
						"memory":      "256Mi",
						"concurrency": "100",
					},
				},
			},
		},
	}
	expectations.Check(t)
}

func getDetector() python.Detect {
	buf := bytes.NewBuffer(nil)
	logger := bard.NewLogger(buf)
	return python.Detect{
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
		dc.Platform.Environment[EnvModuleName] = value
	}
}

func withFunctionName(value string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Platform.Environment[EnvFunctionName] = value
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
	planName := "python-function"
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
						t.Errorf("unexpected value in '%s' requires metadata for key %s", planName, k)
						return
					}
				}
			}
		}
	}
}
