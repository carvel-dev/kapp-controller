// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sharedNamespaces = []string{
	"default",
	"kube-public",
}

type SecureNamespaceFlags struct {
	AllowedSharedNamespaces bool
}

func (s *SecureNamespaceFlags) Set(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&s.AllowedSharedNamespaces, "dangerous-allow-use-of-shared-namespace", false, "Allow use of shared namespaces")
}

func (s *SecureNamespaceFlags) CheckForDisallowedSharedNamespaces(namespace string) error {
	if s.AllowedSharedNamespaces {
		return nil
	}
	for _, ns := range sharedNamespaces {
		if namespace == ns {
			return fmt.Errorf("Creating sensitive resources in a shared namespace (%s)"+
				"(hint: Specify a namespace using the '-n' flag or use kubeconfig to change default namespace 'kubectl config set-context --current --namespace=private-namespace'."+
				"Or use '--dangerous-allow-allow-use-of-shared-namespace=%s' to allow use of shared namespace)", namespace, namespace)
		}
	}
	return nil
}
