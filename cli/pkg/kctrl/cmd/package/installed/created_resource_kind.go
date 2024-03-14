// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package installed

import "fmt"

type CreatedResourceKind string

const (
	KindClusterRole        CreatedResourceKind = "ClusterRole"
	KindClusterRoleBinding CreatedResourceKind = "ClusterRoleBinding"
	KindSecret             CreatedResourceKind = "Secret"
	KindServiceAccount     CreatedResourceKind = "ServiceAccount"
	KindNamespace          CreatedResourceKind = "Namespace"
)

func (k CreatedResourceKind) AsString() string {
	return string(k)
}

func (k CreatedResourceKind) Resource() string {
	switch k {
	case KindClusterRole:
		return "clusterroles"
	case KindClusterRoleBinding:
		return "clusterrolebindings"
	case KindSecret:
		return "secrets"
	case KindServiceAccount:
		return "serviceaccounts"
	}
	return ""
}

func (k CreatedResourceKind) Name(pkgiName string, pkgiNamespace string) string {
	switch k {
	case KindClusterRole:
		return fmt.Sprintf(ClusterRoleName, pkgiName, pkgiNamespace)
	case KindClusterRoleBinding:
		return fmt.Sprintf(ClusterRoleBindingName, pkgiName, pkgiNamespace)
	case KindSecret:
		return fmt.Sprintf(SecretName, pkgiName, pkgiNamespace)
	case KindServiceAccount:
		return fmt.Sprintf(ServiceAccountName, pkgiName, pkgiNamespace)
	}
	return ""
}
