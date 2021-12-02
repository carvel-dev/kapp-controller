// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/cppforlife/cobrautil"
)

const (
	cmdGroupKey = "kapp-group"
)

var (
	PackageHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "package",
		Title: "Package Commands:",
	}
	RestOfCommandsHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "", // default
		Title: "Available/Other Commands:",
	}
)
