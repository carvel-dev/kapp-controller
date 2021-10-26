// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageavailable

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "available",
		ValidArgs: []string{"list", "get"},
		Aliases:   []string{"pkga"},
		Short:     "PackageAvailable",
	}
	return cmd
}
