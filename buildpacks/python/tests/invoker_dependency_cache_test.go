// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpacks/libcnb"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"kn-fn/python-function-buildpack/mock_python"
	"kn-fn/python-function-buildpack/python"
)

func TestInvokerDependencyCache(t *testing.T) {
	spec.Run(t, "InvokerDependencyCache", testInvokerDependencyCache, spec.Report(report.Terminal{}), spec.Parallel())
}

func testInvokerDependencyCache(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		ctrl              *gomock.Controller
		mockCommandRunner *mock_python.MockCommandRunner
		invokerDepCache   *python.InvokerDependencyCache
		layer             libcnb.Layer
		layersDir         string
	)

	it.Before(func() {
		ctrl = gomock.NewController(t)
		mockCommandRunner = mock_python.NewMockCommandRunner(ctrl)

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
		it("returns invoker dep cache layer name", func() {
			invokerDepCache = python.NewInvokerDependencyCache(
				libpak.BuildpackDependency{},
				libpak.DependencyCache{},
				python.NewDefaultCommandRunner(),
			)

			Expect(invokerDepCache.Name()).To(Equal("invoker-deps"))
		})
	})

	when("#Contribute", func() {
		setupMockPipInstall := func(targetDir string) {
			mockCommandRunner.EXPECT().Run(gomock.Any()).DoAndReturn(func(cmd *exec.Cmd) (string, error) {
				depsDir := cmd.Args[6]
				invokerDepsDir := filepath.Join(depsDir, "sample-invoker-deps")

				// Validate command called correctly
				cmdString := strings.Join(cmd.Args, " ")
				Expect(cmdString).To(MatchRegexp(
					"^pip install --target %s --no-index --find-links %s --compile --disable-pip-version-check --ignore-installed --exists-action=w %s",
					targetDir,
					depsDir,
					invokerDepsDir,
				))

				// Pretend to install deps to target dir
				entries, err := os.ReadDir(invokerDepsDir)
				Expect(err).NotTo(HaveOccurred())
				for _, entry := range entries {
					src, err := os.ReadFile(filepath.Join(invokerDepsDir, entry.Name()))
					Expect(err).NotTo(HaveOccurred())
					Expect(os.WriteFile(filepath.Join(targetDir, entry.Name()), src, os.ModePerm)).To(Succeed())
				}

				return "", nil
			})
		}

		it.Before(func() {
			invokerDepCache = python.NewInvokerDependencyCache(
				libpak.BuildpackDependency{
					ID:     "invoker-deps",
					URI:    "https://www.example.com/sample-invoker-deps.zip",
					SHA256: "some-sha256",
				},
				libpak.DependencyCache{
					CachePath: filepath.Join("testdata", "invoker-dep-cache", "cache"),
				},
				mockCommandRunner,
			)
		})

		it("installs invoker dependencies to layer path", func() {
			targetDir := filepath.Join(layer.Path, "install")
			setupMockPipInstall(targetDir)

			var err error
			layer, err = invokerDepCache.Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			contents, err := os.ReadFile(filepath.Join(targetDir, "sample-dep.whl"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(Equal("sample dependency"))
		})

		it("sets layer as cached launch layer", func() {
			mockCommandRunner.EXPECT().Run(gomock.Any()).Return("", nil)

			var err error
			layer, err = invokerDepCache.Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			Expect(layer.LayerTypes).To(Equal(libcnb.LayerTypes{
				Build:  false,
				Cache:  true,
				Launch: true,
			}))
		})

		it("sets python path to install directory", func() {
			mockCommandRunner.EXPECT().Run(gomock.Any()).Return("", nil)

			var err error
			layer, err = invokerDepCache.Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			expectedPath := filepath.Join(layer.Path, "install")

			Expect(invokerDepCache.PythonPath()).To(Equal(expectedPath))

			Expect(layer.LaunchEnvironment).To(Equal(libcnb.Environment{
				"PYTHONPATH.append": expectedPath,
				"PYTHONPATH.delim":  string(os.PathListSeparator),
			}))
		})
	})
}
