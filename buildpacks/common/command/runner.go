// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package command

import (
	"os/exec"
)

//go:generate mockgen -destination ../mock_command/runner.go . Runner
type Runner interface {
	Run(cmd *exec.Cmd) (output string, err error)
}

type DefaultRunner struct{}

func NewDefaultRunner() *DefaultRunner {
	return &DefaultRunner{}
}

func (dcr *DefaultRunner) Run(cmd *exec.Cmd) (output string, err error) {
	buff, err := cmd.CombinedOutput()
	return string(buff), err
}
