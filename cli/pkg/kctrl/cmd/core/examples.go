// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"
)

type PackageCommandTreeOpts struct {
	BinaryName     string
	PositionalArgs bool

	Color bool
	JSON  bool
}

type Example struct {
	Description string
	Args        []string
}

func (e Example) asString(nameFlag string, opts PackageCommandTreeOpts) string {
	command := opts.BinaryName
	for _, arg := range e.Args {
		if opts.PositionalArgs && arg == nameFlag {
			continue
		}
		command += " " + arg
	}
	return fmt.Sprintf("# %s \n%s", e.Description, command)
}

type Examples []Example

func (es Examples) Description(nameFlag string, opts PackageCommandTreeOpts) string {
	var description string
	for _, example := range es {
		description += example.asString(nameFlag, opts) + "\n\n"
	}
	return strings.TrimSuffix(description, "\n\n")
}
