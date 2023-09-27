// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
)

// Secrets represents a secret with target cluster kubeconfig
type Secrets struct {
	Name       string
	Namespace  string
	Kubeconfig string
}

// ForTargetCluster can be used to get secret with target cluster kubeconfig
func (s Secrets) ForTargetCluster() string {
	indentedKubeconfig := ""
	for _, line := range strings.Split(s.Kubeconfig, "\n") {
		if line != "" {
			indentedKubeconfig += "    " + line + "\n"
		}
	}
	return fmt.Sprintf(`
---
apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-rule.apps: "delete after deleting kappctrl-e2e.k14s.io/apps"
    kapp.k14s.io/change-rule.instpkgs: "delete after deleting kappctrl-e2e.k14s.io/packageinstalls"
type: Opaque
stringData:
  value: |
%s
`, s.Name, s.Namespace, indentedKubeconfig)
}
