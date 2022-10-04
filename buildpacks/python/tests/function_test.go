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

	"kn-fn/python-function-buildpack/python"
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
			function := python.NewFunction()

			Expect(function.Name()).To(Equal("python-function"))
		})
	})

	when("#Contribute", func() {
		when("always", func() {
			it.Before(func() {
				function := python.NewFunction(
					python.WithFunctionClass("some_module.some_function", true),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
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
				function := python.NewFunction(
					python.WithFunctionClass("some_module.some_function", true),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("sets function launch env vars as overrides", func() {
				Expect(layer.LaunchEnvironment).To(Equal(libcnb.Environment{
					"MODULE_NAME.override":   "some_module",
					"FUNCTION_NAME.override": "some_function",
				}))
			})

			it("sets layer metadata accordingly", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"bp-function-class":          "some_module.some_function",
					"bp-function-class-override": "true",
				}))
			})
		})

		when("not overriding function class", func() {
			it.Before(func() {
				function := python.NewFunction(
					python.WithFunctionClass("some_module.some_function", false),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("sets function launch env vars as defaults", func() {
				Expect(layer.LaunchEnvironment).To(Equal(libcnb.Environment{
					"MODULE_NAME.default":   "some_module",
					"FUNCTION_NAME.default": "some_function",
				}))
			})

			it("sets layer metadata accordingly", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"bp-function-class":          "some_module.some_function",
					"bp-function-class-override": "false",
				}))
			})
		})

		when("func.yaml envs are configured", func() {
			it.Before(func() {
				function := python.NewFunction(
					python.WithFunctionClass("some_module.some_function", true),
					python.WithFuncYamlEnvs(map[string]any{
						"SOME_VAR": "SOME_VALUE",
					}),
				)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("adds env vars to launch environment as defaults", func() {
				Expect(layer.LaunchEnvironment).To(Equal(libcnb.Environment{
					"MODULE_NAME.override":   "some_module",
					"FUNCTION_NAME.override": "some_function",
					"SOME_VAR.default":       "SOME_VALUE",
				}))
			})

			it("sets layer metadata accordingly", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"bp-function-class":          "some_module.some_function",
					"bp-function-class-override": "true",
					"SOME_VAR":                   "SOME_VALUE",
				}))
			})
		})
	})
}
