// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"k8s.io/utils/pointer"
	function "knative.dev/kn-plugin-func"

	"kn-fn/buildpacks/tests"
	"kn-fn/java-function-buildpack/java"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		detect        java.Detect
		cleanupAppDir func()
		context       libcnb.DetectContext
	)

	it.Before(func() {
		detect = java.Detect{
			Logger: tests.NewLogger(),
		}
	})

	it.After(func() {
		cleanupAppDir()
	})

	when("func.yaml exists", func() {
		it.Before(func() {
			var appDir string
			appDir, cleanupAppDir = tests.SetupTestDirectory(
				tests.WithFuncYaml(),
			)
			context = makeDetectContext(
				withApplicationPath(appDir),
			)
		})

		it("passes detection", func() {
			result, err := detect.Detect(context)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Pass).To(BeTrue())
		})
	})

	when("BP_FUNCTION is configured without func.yaml", func() {
		it.Before(func() {
			t.Setenv("BP_FUNCTION", "function.Handler")

			var appDir string
			appDir, cleanupAppDir = tests.SetupTestDirectory()
			context = makeDetectContext(
				withApplicationPath(appDir),
			)
		})

		it("passes detection", func() {
			result, err := detect.Detect(context)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Pass).To(BeTrue())
		})
	})

	when("func.yaml does not exist", func() {
		it.Before(func() {
			var appDir string
			appDir, cleanupAppDir = tests.SetupTestDirectory()
			context = makeDetectContext(
				withApplicationPath(appDir),
			)
		})

		when("BP_FUNCTION is not configured", func() {
			it.Before(func() {
				Expect(os.Unsetenv("BP_FUNCTION")).To(Succeed())
			})

			it("fails detection", func() {
				result, err := detect.Detect(context)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Pass).To(BeFalse())
			})
		})
	})

	when("func.yaml has configuration for envs or options", func() {
		it.Before(func() {
			var appDir string
			appDir, cleanupAppDir = tests.SetupTestDirectory(
				tests.WithFuncEnvs(map[string]string{
					"SOME_VAR": "SOME_VALUE",
				}),
				tests.WithFuncScale(function.ScaleOptions{
					Min: pointer.Int64(1),
					Max: pointer.Int64(42),
				}),
			)
			context = makeDetectContext(
				withApplicationPath(appDir),
			)
		})

		it("includes configuration in build plan", func() {
			result, err := detect.Detect(context)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Plans).To(Equal([]libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{
							Name: "java-function",
						},
					},
					Requires: []libcnb.BuildPlanRequire{
						{
							Name: "java-function",
							Metadata: map[string]any{
								"launch": true,
								"func_yaml_envs": map[string]string{
									"SOME_VAR": "SOME_VALUE",
								},
								"func_yaml_options": map[string]string{
									"options-scale": "min: 1\nmax: 42\n",
								},
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
				}}))
		})
	})
}

func makeDetectContext(opts ...func(*libcnb.DetectContext)) libcnb.DetectContext {
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
