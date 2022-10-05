// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"os"
	"path/filepath"

	"kn-fn/buildpacks/tests"
)

func WithTomcatJar() tests.SetupOpts {
	return func(directory string) {
		path := filepath.Join(directory, "tomcat-embed-core-10.1.0.jar")
		_, err := os.Create(path)
		if err != nil {
			panic(err)
		}
	}
}
