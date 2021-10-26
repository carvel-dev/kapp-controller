// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed",
		Aliases: []string{"pkgi"},
		Short:   "PackageInstall",
	}
	return cmd
}
