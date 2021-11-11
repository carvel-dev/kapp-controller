// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed",
		Aliases: []string{"i"},
		Short:   "PackageInstall",
	}
	return cmd
}
