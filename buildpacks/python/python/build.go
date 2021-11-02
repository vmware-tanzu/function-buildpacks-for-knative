package python

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	_, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)

	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	invokerPlan, ok, err := pr.Resolve("kn-fn-python-invoker")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve kn-fn-python-invoker plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}

	i := NewInvokerFromPlan(invokerPlan, context.Buildpack.Path)
	i.Logger = b.Logger
	result.Layers = append(result.Layers, i)

	command := "python"
	arguments := []string{"-m", "pyfunc"}
	result.Processes = append(result.Processes,
		libcnb.Process{
			Default:   true,
			Type:      "func",
			Command:   command,
			Arguments: arguments,
		},
		libcnb.Process{
			Type:    "shell",
			Command: "bash",
		},
	)

	return result, nil
}
