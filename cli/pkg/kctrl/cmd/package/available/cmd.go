// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package available

import (
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "available",
		Aliases: []string{"a"},
		Short:   "Manage available packages",
		Annotations: map[string]string{
			cmdcore.PackageManagementCommandsHelpGroup.Key: cmdcore.PackageManagementCommandsHelpGroup.Value,
		},
	}
	return cmd
}
