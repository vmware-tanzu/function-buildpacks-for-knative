package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	knfunc "knative.dev/kn-plugin-func"
)

const (
	EnvModuleName   = "FUNCTION_MODULE"
	EnvFunctionName = "FUNCTION_NAME"

	HTTPTemplate = `from typing import Any

def {{.Name}}({{range $var, $type := .Arguments}}{{$var}}:{{$type}}{{end}}):
	return "{{.ReturnValue}}"`
)

type HTTPFunction struct {
	Module      string
	Name        string
	Arguments   map[string]string
	ReturnValue string
}

type Function interface {
	Generate(path string) error
}

func (f HTTPFunction) Generate(path string) error {
	templ, err := template.New("http-func-template").Parse(HTTPTemplate)
	if err != nil {
		return err
	}

	modulePath := filepath.Join(path, fmt.Sprintf("%s.py", f.Module))
	file, err := os.Create(modulePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return templ.Execute(file, f)
}

type TestApplication struct {
	Function Function

	ApplicationPath  string
	FunctionFilePath string

	cleanup func()
}

func createTestApplication(opts ...func(*TestApplication)) TestApplication {
	appPath, _ := ioutil.TempDir(os.TempDir(), "python-functions-buildpack-test-*")

	app := TestApplication{
		ApplicationPath: appPath,
	}

	for _, opt := range opts {
		opt(&app)
	}

	app.cleanup = func() {
		os.RemoveAll(appPath)
	}

	return app
}

func (a TestApplication) Finish() {
	if a.cleanup != nil {
		a.cleanup()
	}
}

func withDefaultHTTPFunction() func(*TestApplication) {
	return withHTTPFunction("func", "main")
}

func withHTTPFunction(module string, function string) func(*TestApplication) {
	return func(ta *TestApplication) {
		if ta.Function != nil {
			panic(fmt.Errorf("test application already defined function"))
		}

		ta.Function = HTTPFunction{
			Module: module,
			Name:   function,
		}

		cfg, err := knfunc.NewFunction(ta.ApplicationPath)
		if err != nil {
			panic(err)
		}

		cfg.Envs = append(cfg.Envs,
			knfunc.Env{Name: strptr(EnvModuleName), Value: &module},
			knfunc.Env{Name: strptr(EnvFunctionName), Value: &function},
		)

		err = cfg.WriteConfig()
		if err != nil {
			panic(err)
		}

		err = ta.Function.Generate(ta.ApplicationPath)
		if err != nil {
			panic(err)
		}
	}
}

func strptr(str string) *string {
	return &str
}
