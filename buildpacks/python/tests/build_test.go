// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"

	"kn-fn/python-function-buildpack/python"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", testBuild, spec.Report(report.Terminal{}))
}

func testBuild(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		build         python.Build
		cleanupAppDir func()
		context       libcnb.BuildContext
	)

	it.Before(func() {
		build = python.Build{
			Logger: NewLogger(),
		}
	})

	it.After(func() {
		cleanupAppDir()
	})

	when("#Build", func() {
		var result libcnb.BuildResult

		it.Before(func() {
			var (
				appDir string
				err    error
			)

			Expect(os.Setenv("BP_FUNCTION", "some_module.some_function"))
			appDir, cleanupAppDir = SetupTestDirectory(
				WithFuncYaml(),
				WithFunctionFile("some_module", "some_function", HTTPFuncTemplate),
			)

			context = makeBuildContext(
				withBuildApplicationPath(appDir),
				withDependencies([]map[string]any{
					{"id": "invoker-deps", "version": "1.2.3"},
					{"id": "invoker", "version": "2.3.4"},
				}),
				withOptions(map[string]any{
					"some-option":       "some-value",
					"some-other-option": "some-other-value",
				}),
			)

			result, err = build.Build(context)
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_FUNCTION"))
		})

		it("adds expected layers", func() {
			var layers []string
			for _, l := range result.Layers {
				layers = append(layers, l.Name())
			}

			Expect(layers).To(Equal([]string{
				"invoker-deps",
				"invoker",
				"python-function",
				"validation",
			}))
		})

		it("adds expected labels", func() {
			Expect(result.Labels).To(Equal([]libcnb.Label{
				{Key: "some-option", Value: "some-value"},
				{Key: "some-other-option", Value: "some-other-value"},
			}))
		})

		it("adds launch command", func() {
			Expect(result.Processes).To(Equal([]libcnb.Process{
				{
					Type:             "func",
					Command:          "python",
					Arguments:        []string{"-m", "pyfunc", "start"},
					Direct:           false,
					WorkingDirectory: "",
					Default:          true,
				},
			}))
		})
	})
}

func makeBuildContext(opts ...func(*libcnb.BuildContext)) libcnb.BuildContext {
	ctx := libcnb.BuildContext{
		Application: libcnb.Application{},
		Buildpack: libcnb.Buildpack{
			Metadata: map[string]any{},
		},
		Platform: libcnb.Platform{
			Environment: map[string]string{},
		},
		Plan: libcnb.BuildpackPlan{
			Entries: []libcnb.BuildpackPlanEntry{
				{
					Name: "python-function",
					Metadata: map[string]any{
						"func_yaml_envs": map[string]any{},
					},
				},
			},
		},
	}

	for _, opt := range opts {
		opt(&ctx)
	}

	return ctx
}

func withBuildApplicationPath(path string) func(*libcnb.BuildContext) {
	return func(bc *libcnb.BuildContext) {
		bc.Application.Path = path
	}
}

func withDependencies(deps []map[string]any) func(ctx *libcnb.BuildContext) {
	return func(bc *libcnb.BuildContext) {
		bc.Buildpack.Metadata["dependencies"] = deps
	}
}

func withOptions(options map[string]any) func(ctx *libcnb.BuildContext) {
	return func(bc *libcnb.BuildContext) {
		bc.Plan.Entries[0].Metadata["func_yaml_options"] = options
	}
}
