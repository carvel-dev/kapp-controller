// Copyright 2024 The Carvel Authors.
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

// ForClusterYAML can be used to get service account with cluster wide permissions
func (sa ServiceAccounts) ForClusterYAML() string {
	return fmt.Sprintf(`
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kappctrl-e2e-ns-sa
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role-binding
subjects:
- kind: ServiceAccount
  name: kappctrl-e2e-ns-sa
  namespace: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kappctrl-e2e-ns-role
`, sa.Namespace)
}

// ForDefaultNamespaceYAML can be used to get service account with permissions in defaultNamespace
func (sa ServiceAccounts) ForDefaultNamespaceYAML(defaultNamespace, pkgiNamespace string) string {
	return fmt.Sprintf(`
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kappctrl-e2e-ns-sa
  namespace: %[2]s
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role
  namespace: %[1]s
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role-binding
  namespace: %[1]s
subjects:
- kind: ServiceAccount
  name: kappctrl-e2e-ns-sa
  namespace: %[2]s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kappctrl-e2e-ns-role
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role
  namespace: %[2]s
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: ["get", "list", "create", "delete", "update"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role-binding
  namespace: %[2]s
subjects:
- kind: ServiceAccount
  name: kappctrl-e2e-ns-sa
  namespace: %[2]s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kappctrl-e2e-ns-role
`, defaultNamespace, pkgiNamespace)
}
