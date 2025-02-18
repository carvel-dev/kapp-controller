// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	cmdcore "carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed",
		Aliases: []string{"i"},
		Short:   "Manage installed packages",
		Annotations: map[string]string{
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value,
		},
	}
	return cmd
}
