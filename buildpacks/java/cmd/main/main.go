// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"os"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"

	"kn-fn/java-function-buildpack/java"
)

func main() {
	libpak.Main(
		java.Detect{},
		java.Build{Logger: bard.NewLogger(os.Stdout)},
	)
}
