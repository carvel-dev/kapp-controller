// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CreatedResourceKind string

const (
	// For use with packages created using kapp-ctrl-cli
	// KappPkgAnnotation       = "packaging.carvel.dev/package"
	// KappPkgAnnotationPrefix = "package"

	// For use with packages created with Tanzu CLI
	KappPkgAnnotation       = "tkg.tanzu.vmware.com/tanzu-package"
	KappPkgAnnotationPrefix = "tanzu-package"

	KindClusterRole        CreatedResourceKind = "ClusterRole"
	KindClusterRoleBinding CreatedResourceKind = "ClusterRoleBinding"
	KindSecret             CreatedResourceKind = "Secret"
	KindServiceAccount     CreatedResourceKind = "ServiceAccount"

	ClusterRoleBindingName = "%s-%s-cluster-rolebinding"
	ClusterRoleName        = "%s-%s-cluster-role"
	SecretName             = "%s-%s-values"
	ServiceAccountName     = "%s-%s-sa"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "installed",
		Aliases: []string{"pkgi"},
		Short:   "PackageInstall",
	}
	return cmd
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
