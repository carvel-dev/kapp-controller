// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg", "p"},
		Short:   "Package",
		Annotations: map[string]string{
			cmdcore.PackageHelpGroup.Key: cmdcore.PackageHelpGroup.Value,
		},
	}
	return cmd
}
