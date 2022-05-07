// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "Create packages",
	}
	return cmd
}
