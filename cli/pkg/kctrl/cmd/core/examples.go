// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"
)

type Example struct {
	Description string
	Args        []string
}

func (e Example) asString(binaryName, nameFlag string, positionalNameArg bool) string {
	command := binaryName
	for _, arg := range e.Args {
		if positionalNameArg && arg == nameFlag {
			continue
		}
		command += " " + arg
	}
	return fmt.Sprintf("# %s \n%s", e.Description, command)
}

type Examples []Example

func (es Examples) Description(binaryName, nameFlag string, positionalNameArg bool) string {
	var description string
	for _, example := range es {
		description += example.asString(binaryName, nameFlag, positionalNameArg) + "\n\n"
	}
	return strings.TrimSuffix(description, "\n\n")
}
