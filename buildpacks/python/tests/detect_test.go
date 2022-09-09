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
	function "knative.dev/kn-plugin-func"
	"knative.dev/pkg/ptr"

	"kn-fn/python-function-buildpack/python"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		detect        python.Detect
		cleanupAppDir func()
		context       libcnb.DetectContext
	)

	it.Before(func() {
		detect = python.Detect{
			Logger: NewLogger(),
		}
	})

	it.After(func() {
		cleanupAppDir()
	})

	when("#Detect", func() {
		when("func.yaml exists", func() {
			it.Before(func() {
				var appDir string
				appDir, cleanupAppDir = SetupTestDirectory(
					WithFuncYaml(),
				)
				context = makeDetectContext(
					withDetectApplicationPath(appDir),
				)
			})

			it("passes detection", func() {
				result, err := detect.Detect(context)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Pass).To(BeTrue())
			})
		})

		when("BP_FUNCTION is configured", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_FUNCTION", "some.function")).To(Succeed())

				var appDir string
				appDir, cleanupAppDir = SetupTestDirectory()
				context = makeDetectContext(
					withDetectApplicationPath(appDir),
				)
			})

			it.After(func() {
				Expect(os.Unsetenv("BP_FUNCTION")).To(Succeed())
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
				appDir, cleanupAppDir = SetupTestDirectory()
				context = makeDetectContext(
					withDetectApplicationPath(appDir),
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

			when("BP_FUNCTION is configured incorrectly", func() {
				it.Before(func() {
					Expect(os.Setenv("BP_FUNCTION", "invalid function")).To(Succeed())
				})

				it.After(func() {
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
				appDir, cleanupAppDir = SetupTestDirectory(
					WithFuncEnvs(map[string]string{
						"SOME_VAR": "SOME_VALUE",
					}),
					WithFuncScale(function.ScaleOptions{
						Min: ptr.Int64(1),
						Max: ptr.Int64(42),
					}),
				)
				context = makeDetectContext(
					withDetectApplicationPath(appDir),
				)
			})

			it("includes configuration in build plan", func() {
				result, err := detect.Detect(context)
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Plans).To(Equal([]libcnb.BuildPlan{{
					Provides: []libcnb.BuildPlanProvide{
						{
							Name: "python-function",
						},
					},
					Requires: []libcnb.BuildPlanRequire{
						{
							Name: "python-function",
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
				}}))
			})
		})
	})
}

func makeDetectContext(opts ...func(*libcnb.DetectContext)) libcnb.DetectContext {
	ctx := libcnb.DetectContext{
		Application: libcnb.Application{},
		Buildpack:   libcnb.Buildpack{},
		Platform: libcnb.Platform{
			Environment: map[string]string{},
		},
	}

	for _, opt := range opts {
		opt(&ctx)
	}

	return ctx
}

func withDetectApplicationPath(path string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Application.Path = path
	}
}
