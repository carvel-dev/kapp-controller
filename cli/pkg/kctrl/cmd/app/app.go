// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"a"},
		Short:   "App",
	}
	return cmd
}
