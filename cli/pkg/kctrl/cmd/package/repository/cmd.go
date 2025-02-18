// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	cmdcore "carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "repository",
		Aliases: []string{"repo", "r"},
		Short:   "Manage package repositories",
		Annotations: map[string]string{
			cmdcore.PackageRepoHelpGroup.Key: cmdcore.PackageRepoHelpGroup.Value,
		},
	}
	return cmd
}
