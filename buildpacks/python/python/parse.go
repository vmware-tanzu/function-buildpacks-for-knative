// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package python

import (
	"fmt"
	"strings"
)

func parseFunctionClass(value string) (moduleName string, functionName string, err error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid function class '%s', expected format: '<module name>.<function name>'", value)
	}
	return parts[0], parts[1], nil
}
