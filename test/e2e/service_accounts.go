// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
)

type ServiceAccounts struct {
	Namespace string
}

func (sa ServiceAccounts) ForNamespaceYAML() string {
	return fmt.Sprintf(`
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kappctrl-e2e-ns-sa
  annotations:
    kapp.k14s.io/change-rule.apps: "delete after deleting kappctrl-e2e.k14s.io/apps"
    kapp.k14s.io/change-rule.instpkgs: "delete after deleting kappctrl-e2e.k14s.io/packageinstalls"
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role
  annotations:
    kapp.k14s.io/change-rule.apps: "delete after deleting kappctrl-e2e.k14s.io/apps"
    kapp.k14s.io/change-rule.instpkgs: "delete after deleting kappctrl-e2e.k14s.io/packageinstalls"
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role-binding
  annotations:
    kapp.k14s.io/change-rule.apps: "delete after deleting kappctrl-e2e.k14s.io/apps"
    kapp.k14s.io/change-rule.instpkgs: "delete after deleting kappctrl-e2e.k14s.io/packageinstalls"
subjects:
- kind: ServiceAccount
  name: kappctrl-e2e-ns-sa
  namespace: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kappctrl-e2e-ns-role
`, sa.Namespace)
}
