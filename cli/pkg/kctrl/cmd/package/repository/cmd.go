// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
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
