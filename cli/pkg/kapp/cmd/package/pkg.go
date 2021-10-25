// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg"},
		Short:   "Package",
	}
	return cmd
}
