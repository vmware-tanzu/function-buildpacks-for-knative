package main

import (
	"os"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"

	"kn-fn/python-function-buildpack/python"
)

func main() {
	logger := bard.NewLogger(os.Stdout)
	libpak.Main(
		python.Detect{Logger: logger},
		python.Build{Logger: logger},
	)
}
