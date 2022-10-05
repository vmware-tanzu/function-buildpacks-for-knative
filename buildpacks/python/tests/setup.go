// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"kn-fn/buildpacks/tests"
)

const (
	HTTPFuncTemplate = `from typing import Any

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
	Generate(path, funcTemplate string) error
}

func (f HTTPFunction) Generate(path, funcTemplate string) error {
	templ, err := template.New("http-func-template").Parse(funcTemplate)
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

func WithFunctionFile(module, function, funcTemplate string) tests.SetupOpts {
	return func(directory string) {
		f := HTTPFunction{
			Module: module,
			Name:   function,
		}

		err := f.Generate(directory, funcTemplate)
		if err != nil {
			panic(err)
		}
	}
}
