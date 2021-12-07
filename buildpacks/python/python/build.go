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

	planResolver := libpak.PlanEntryResolver{Plan: context.Plan}

	dependencyCache, err := libpak.NewDependencyCache(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}
	dependencyCache.Logger = b.Logger

	dependencyResolver, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	invokerDep, err := dependencyResolver.Resolve("invoker", "")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}

	invokerLayer, invokerBOM := NewInvoker(invokerDep, dependencyCache)
	invokerLayer.Logger = b.Logger
	result.Layers = append(result.Layers, invokerLayer)
	result.BOM.Entries = append(result.BOM.Entries, invokerBOM)

	functionPlan, ok, err := planResolver.Resolve("python-function")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve python-function plan entry\n%w", err)
	}
	if !ok {
		return result, nil
	}
	functionLayer := NewFunction(functionPlan)
	result.Layers = append(result.Layers, functionLayer)

	validationLayer := NewFunctionValidationLayer(functionPlan, context.Application.Path)
	result.Layers = append(result.Layers, validationLayer)

	command := "python"
	arguments := []string{"-m", "pyfunc", "start"}
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
