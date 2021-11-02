package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	PYTHON_HTTP_TEMPLATE = `
from typing import Any

def handler(req: Any):
	return "Hello world!"`
)

type TestApplication struct {
	ModuleName   string
	FunctionName string

	ApplicationPath  string
	FunctionFilePath string

	cleanup func()
}

func createTestApplication(opts ...func(*TestApplication)) TestApplication {
	appPath, _ := ioutil.TempDir(os.TempDir(), "python-functions-buildpack-test-*")

	app := TestApplication{
		ModuleName:   "undefined",
		FunctionName: "undefined",

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

func withFunction(module string, function string) func(*TestApplication) {
	return func(ta *TestApplication) {
		ta.ModuleName = module
		ta.FunctionName = function

		ta.FunctionFilePath = filepath.Join(ta.ApplicationPath, fmt.Sprintf("%s.py", module))
		file, _ := os.Create(ta.FunctionFilePath)

		file.WriteString(PYTHON_HTTP_TEMPLATE)

		defer file.Close()
	}
}
