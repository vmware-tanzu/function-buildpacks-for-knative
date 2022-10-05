// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"os"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"

	"kn-fn/buildpacks/command"
	"kn-fn/python-function-buildpack/python"
)

func main() {
	logger := bard.NewLogger(os.Stdout)
	libpak.Main(
		python.Detect{Logger: logger},
		python.Build{Logger: logger, CommandRunner: command.NewDefaultRunner()},
	)
}
