// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/cppforlife/cobrautil"
)

const (
	cmdGroupKey = "kctrl-group"
)

var (
	PackageHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "package",
		Title: "Package Commands:",
	}
	PackageRepoHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "repository",
		Title: "Package Repository Commands:",
	}
	AppHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "app",
		Title: "App Commands:",
	}
	DevHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "dev",
		Title: "Development Commands:",
	}
	PackageManagementCommandsHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "package-management",
		Title: "Package Management Commands:",
	}
	PackageAuthoringCommandsHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "package-authoring",
		Title: "Package Authoring Commands:",
	}
	AppManagementCommandsHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "app-management",
		Title: "App Management Commands:",
	}
	RestOfCommandsHelpGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "", // default
		Title: "Available/Other Commands:",
	}
)
