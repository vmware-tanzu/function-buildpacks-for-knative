// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"kn-fn/java-function-buildpack/java"
)

func TestInvoker(t *testing.T) {
	spec.Run(t, "Invoker", testInvoker, spec.Report(report.Terminal{}), spec.Parallel())
}

func testInvoker(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		invoker   *java.Invoker
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
		it("returns invoker layer name", func() {
			invoker = java.NewInvoker(
				libpak.BuildpackDependency{},
				libpak.DependencyCache{},
			)

			Expect(invoker.Name()).To(Equal("invoker"))
		})
	})

	when("#Contribute", func() {
		it.Before(func() {
			invoker = java.NewInvoker(
				libpak.BuildpackDependency{
					ID:     "invoker",
					URI:    "https://www.example.com/sample-invoker.zip",
					SHA256: "some-sha256",
				},
				libpak.DependencyCache{
					CachePath: filepath.Join("testdata", "invoker", "cache"),
				},
			)

			var err error
			layer, err = invoker.Contribute(layer)
			Expect(err).NotTo(HaveOccurred())
		})

		it("extracts invoker to layer path", func() {
			contents, err := os.ReadFile(filepath.Join(layer.Path, "sample-invoker-exec"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal("sample invoker"))
		})

		it("sets layer as uncached launch layer", func() {
			Expect(layer.LayerTypes).To(Equal(libcnb.LayerTypes{
				Build:  false,
				Cache:  false,
				Launch: true,
			}))
		})

		it("prepends layer path to class path", func() {
			Expect(layer.LaunchEnvironment).To(Equal(libcnb.Environment{
				"CLASSPATH.prepend": layer.Path,
				"CLASSPATH.delim":   string(os.PathListSeparator),
			}))
		})
	})
}
