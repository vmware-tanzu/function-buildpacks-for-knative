// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpacks/libcnb"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"kn-fn/python-function-buildpack/mock_python"
	"kn-fn/python-function-buildpack/python"
)

func TestFunctionValidator(t *testing.T) {
	spec.Run(t, "FunctionValidator", testFunctionValidator, spec.Report(report.Terminal{}), spec.Sequential())
}

func testFunctionValidator(t *testing.T, when spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	var (
		layer                libcnb.Layer
		layersDir            string
		ctrl                 *gomock.Controller
		mockCommandRunner    *mock_python.MockCommandRunner
		mockInvokerLayer     *mock_python.MockLayer
		mockInvokerDepsLayer *mock_python.MockLayer
	)

	it.Before(func() {
		ctrl = gomock.NewController(t)
		mockCommandRunner = mock_python.NewMockCommandRunner(ctrl)
		mockInvokerLayer = mock_python.NewMockLayer(ctrl)
		mockInvokerDepsLayer = mock_python.NewMockLayer(ctrl)

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
		it("returns function validator layer name", func() {
			function := python.NewFunctionValidationLayer(
				"",
				mockInvokerLayer,
				mockInvokerDepsLayer,
				mockCommandRunner,
			)

			Expect(function.Name()).To(Equal("validation"))
		})
	})

	when("#Contribute", func() {
		setupMockPythonCheck := func(mockOutput string, mockError error) {
			mockCommandRunner.EXPECT().Run(gomock.Any()).DoAndReturn(func(cmd *exec.Cmd) (string, error) {
				// Validate command called correctly
				cmdString := strings.Join(cmd.Args, " ")
				Expect(cmdString).To(Equal("python -m pyfunc check -s some/app/path -m some_module -f some_function"))
				Expect(cmd.Env).To(ContainElements(
					"PYTHONPATH=invoker/python/path:invoker/deps/python/path:env/python/path",
				))

				// Pretend to run check
				return mockOutput, mockError
			})
		}

		it.Before(func() {
			mockInvokerDepsLayer.EXPECT().PythonPath().Return("invoker/deps/python/path").AnyTimes()
			mockInvokerLayer.EXPECT().PythonPath().Return("invoker/python/path").AnyTimes()
			t.Setenv("PYTHONPATH", "env/python/path")
		})

		when("function is valid", func() {
			it.Before(func() {
				function := python.NewFunctionValidationLayer(
					"some/app/path",
					mockInvokerLayer,
					mockInvokerDepsLayer,
					mockCommandRunner,
					python.WithValidationFunctionClass("some_module", "some_function"),
				)

				setupMockPythonCheck("", nil)

				var err error
				layer, err = function.Contribute(layer)
				Expect(err).NotTo(HaveOccurred())
			})

			it("sets layer as uncached build layer", func() {
				Expect(layer.LayerTypes).To(Equal(libcnb.LayerTypes{
					Build:  true,
					Cache:  false,
					Launch: false,
				}))
			})
		})

		when("function is invalid", func() {
			var err error

			it.Before(func() {
				setupMockPythonCheck("check failed", errors.New("invalid function"))

				function := python.NewFunctionValidationLayer(
					"some/app/path",
					mockInvokerLayer,
					mockInvokerDepsLayer,
					mockCommandRunner,
					python.WithValidationFunctionClass("some_module", "some_function"),
				)

				_, err = function.Contribute(layer)
			})

			it("returns check error", func() {
				Expect(err).To(MatchError("check failed: invalid function"))
			})
		})
	})
}
