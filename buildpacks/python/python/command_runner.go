// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"os/exec"
)

//go:generate mockgen -destination ../mock_python/command_runner.go . CommandRunner
type CommandRunner interface {
	Run(cmd *exec.Cmd) (output string, err error)
}

type DefaultCommandRunner struct{}

func NewDefaultCommandRunner() *DefaultCommandRunner {
	return &DefaultCommandRunner{}
}

func (dcr *DefaultCommandRunner) Run(cmd *exec.Cmd) (output string, err error) {
	buff, err := cmd.CombinedOutput()
	return string(buff), err
}
