// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"kn-fn/buildpacks/tests"
	"kn-fn/java-function-buildpack/java"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", testBuild, spec.Report(report.Terminal{}))
}

func testBuild(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		build         java.Build
		cleanupAppDir func()
		context       libcnb.BuildContext
	)

	it.Before(func() {
		build = java.Build{
			Logger: tests.NewLogger(),
		}
	})

	it.After(func() {
		cleanupAppDir()
	})

	when("#Build", func() {
		var result libcnb.BuildResult

		when("always", func() {
			it.Before(func() {
				var (
					appDir string
					err    error
				)

				appDir, cleanupAppDir = tests.SetupTestDirectory(
					tests.WithFuncYaml(),
				)

				context = makeBuildContext(
					withBuildApplicationPath(appDir),
					withDependencies([]map[string]any{
						{"id": "invoker", "version": "2.3.4"},
					}),
					withOptions(map[string]any{
						"some-other-option": "some-other-value",
						"some-option":       "some-value",
					}),
				)

				result, err = build.Build(context)
				Expect(err).NotTo(HaveOccurred())
			})

			it("adds expected labels sorted", func() {
				Expect(result.Labels).To(Equal([]libcnb.Label{
					{Key: "some-option", Value: "some-value"},
					{Key: "some-other-option", Value: "some-other-value"},
				}))
			})

			it("adds launch command", func() {
				Expect(result.Processes).To(Equal([]libcnb.Process{
					{
						Type:             "func",
						Command:          "java",
						Arguments:        []string{"org.springframework.boot.loader.JarLauncher"},
						Direct:           false,
						WorkingDirectory: "",
						Default:          true,
					},
				}))
			})

		})

		when("without tomcat", func() {

			it.Before(func() {
				var (
					appDir string
					err    error
				)

				appDir, cleanupAppDir = tests.SetupTestDirectory(
					tests.WithFuncYaml(),
				)

				context = makeBuildContext(
					withBuildApplicationPath(appDir),
					withDependencies([]map[string]any{
						{"id": "invoker", "version": "2.3.4"},
					}),
					withOptions(map[string]any{
						"some-other-option": "some-other-value",
						"some-option":       "some-value",
					}),
				)

				result, err = build.Build(context)
				Expect(err).NotTo(HaveOccurred())
			})

			it("adds invoker layer", func() {
				var layers []string
				for _, l := range result.Layers {
					layers = append(layers, l.Name())
				}

				Expect(layers).To(Equal([]string{
					"java-function",
					"invoker",
				}))
			})

		})

		when("with tomcat", func() {

			it.Before(func() {
				var (
					appDir string
					err    error
				)

				appDir, cleanupAppDir = tests.SetupTestDirectory(
					tests.WithFuncYaml(),
					WithTomcatJar(),
				)

				context = makeBuildContext(
					withBuildApplicationPath(appDir),
					withDependencies([]map[string]any{
						{"id": "invoker", "version": "2.3.4"},
					}),
					withOptions(map[string]any{
						"some-other-option": "some-other-value",
						"some-option":       "some-value",
					}),
				)

				result, err = build.Build(context)
				Expect(err).NotTo(HaveOccurred())
			})

			it("does not add the invoker layer", func() {
				var layers []string
				for _, l := range result.Layers {
					layers = append(layers, l.Name())
				}

				Expect(layers).To(Equal([]string{
					"java-function",
				}))
			})
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
					Name: "java-function",
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
