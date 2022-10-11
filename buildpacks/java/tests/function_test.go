// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"kn-fn/java-function-buildpack/java"
)

func TestFunction(t *testing.T) {
	spec.Run(t, "Function", testFunction, spec.Report(report.Terminal{}), spec.Sequential())
}

func testFunction(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		layer     libcnb.Layer
		layersDir string
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		layer.Path = filepath.Join(layersDir, "layer")
		layer.LaunchEnvironment = libcnb.Environment{}
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
	})

	when("#Name", func() {
		it("returns function layer name", func() {
			function := java.NewFunction("")

			Expect(function.Name()).To(Equal("java-function"))
		})
	})

	when("#Contribute", func() {
		when("always", func() {
			it.Before(func() {
				function := java.NewFunction(
					"some/app/path",
					java.WithFunctionClass("", true),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("sets function location to app path", func() {
				Expect(layer.LaunchEnvironment).To(HaveKeyWithValue(
					"SPRING_CLOUD_FUNCTION_LOCATION.default", "some/app/path",
				))
			})

			it("sets layer as uncached launch layer", func() {
				Expect(layer.LayerTypes).To(Equal(libcnb.LayerTypes{
					Build:  false,
					Cache:  false,
					Launch: true,
				}))
			})
		})

		when("overriding function class", func() {
			it.Before(func() {
				function := java.NewFunction(
					"",
					java.WithFunctionClass("somePackage.someFunction", true),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("sets function class env var as override", func() {
				Expect(layer.LaunchEnvironment).To(HaveKeyWithValue(
					"SPRING_CLOUD_FUNCTION_FUNCTION_CLASS.override", "somePackage.someFunction",
				))
			})

			it("sets layer metadata accordingly", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"bp-function-class":          "somePackage.someFunction",
					"bp-function-class-override": "true",
				}))
			})
		})

		when("not overriding function class", func() {
			it.Before(func() {
				function := java.NewFunction(
					"",
					java.WithFunctionClass("somePackage.someFunction", false),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("sets function class env var as default", func() {
				Expect(layer.LaunchEnvironment).To(HaveKeyWithValue(
					"SPRING_CLOUD_FUNCTION_FUNCTION_CLASS.default", "somePackage.someFunction",
				))
			})

			it("sets layer metadata accordingly", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"bp-function-class":          "somePackage.someFunction",
					"bp-function-class-override": "false",
				}))
			})
		})

		when("func.yaml envs are configured", func() {
			it.Before(func() {
				function := java.NewFunction(
					"",
					java.WithFunctionClass("somePackage.someFunction", true),
					java.WithFuncYamlEnvs(map[string]any{
						"SOME_VAR": "SOME_VALUE",
					}),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("adds env vars to launch environment as defaults", func() {
				Expect(layer.LaunchEnvironment).To(HaveKeyWithValue(
					"SOME_VAR.default", "SOME_VALUE",
				))
			})

			it("sets layer metadata accordingly", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"bp-function-class":          "somePackage.someFunction",
					"bp-function-class-override": "true",
					"SOME_VAR":                   "SOME_VALUE",
				}))
			})
		})
	})
}
