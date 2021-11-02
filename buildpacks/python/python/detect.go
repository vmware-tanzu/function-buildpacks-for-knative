package python

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Detect struct {
	Logger bard.Logger

	shouldError bool
}

const (
	DEFAULT_HANDLER = "handler.handler"
)

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	handler, ok := context.Platform.Environment["PYTHON_HANDLER"]
	if !ok {
		handler = DEFAULT_HANDLER
		d.Logger.Info("Using default handler 'handler.handler'")
	}
	d.shouldError = ok

	// Validate that the handler is defined correctly
	sp := strings.Split(handler, ".")
	if len(sp) != 2 {
		return libcnb.DetectResult{}, d.logOrError("expected PYTHON_HANDLER environment variable to be in the form of 'module_name.function_name'")
	}

	module := strings.TrimSpace(sp[0])
	function := strings.TrimSpace(sp[1])

	if len(module) == 0 {
		return libcnb.DetectResult{}, d.logOrError("expected PYTHON_HANDLER environment variable to be in the form of 'module_name.function_name', but module name was ''")
	}

	if len(function) == 0 {
		return libcnb.DetectResult{}, d.logOrError("expected PYTHON_HANDLER environment variable to be in the form of 'module_name.function_name', but function name was ''")
	}

	file := fmt.Sprintf("%s.py", module)
	_, err := os.Stat(filepath.Join(context.Application.Path, file))
	if err != nil {
		return libcnb.DetectResult{}, d.logOrError(fmt.Sprintf("unable to find file '%s'", file))
	}

	result := libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{
						Name: "kn-fn-python-invoker",
					},
				},
				Requires: []libcnb.BuildPlanRequire{
					{
						Name: "kn-fn-python-invoker",
						Metadata: map[string]interface{}{
							"launch": true,
							"build":  true,
							"handler": map[string]string{
								"module":   module,
								"function": function,
								"raw":      fmt.Sprintf("%s.%s", module, function),
							},
						},
					},
					{
						Name: "site-packages",
						Metadata: map[string]interface{}{
							"build":  true,
							"launch": true,
						},
					},
				},
			},
		},
	}

	return result, nil
}

func (d Detect) logOrError(message string) error {
	if d.shouldError {
		return fmt.Errorf(message)
	} else {
		d.Logger.Info(message)
	}
	return nil
}
