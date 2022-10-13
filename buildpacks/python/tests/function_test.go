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
				function := python.NewFunction()

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

		when("func.yaml envs are configured", func() {
			it.Before(func() {
				function := python.NewFunction(
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
					"SOME_VAR.default": "SOME_VALUE",
				}))
			})

			it("adds env vars to layer metadata", func() {
				Expect(layer.Metadata).To(Equal(map[string]any{
					"SOME_VAR": "SOME_VALUE",
				}))
			})
		})
	})
}
