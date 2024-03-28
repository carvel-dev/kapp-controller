// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package installed

import "fmt"

const (
	// For use with packages created using kapp-ctrl-cli
	KctrlPkgAnnotation       = "packaging.carvel.dev/package"
	KctrlPkgAnnotationPrefix = "package-"

	// For use with packages created with Tanzu CLI. To be deprecated
	TanzuPkgAnnotation       = "tkg.tanzu.vmware.com/tanzu-package"
	TanzuPkgAnnotationPrefix = "tanzu-package-"

	ClusterRoleBindingName = "%s-%s-cluster-rolebinding"
	ClusterRoleName        = "%s-%s-cluster-role"
	SecretName             = "%s-%s-values"
	ServiceAccountName     = "%s-%s-sa"
)

type CreatedResourceAnnotations struct {
	name      string
	namespace string
}

func NewCreatedResourceAnnotations(name string, namespace string) *CreatedResourceAnnotations {
	return &CreatedResourceAnnotations{
		name:      name,
		namespace: namespace,
	}
}

func (a *CreatedResourceAnnotations) SecretAnnValue() string {
	return fmt.Sprintf(SecretName, a.name, a.namespace)
}

func (a *CreatedResourceAnnotations) ServiceAccountAnnValue() string {
	return fmt.Sprintf(ServiceAccountName, a.name, a.namespace)
}

func (a *CreatedResourceAnnotations) ClusterRoleAnnValue() string {
	return fmt.Sprintf(ClusterRoleName, a.name, a.namespace)
}

func (a *CreatedResourceAnnotations) ClusterRoleBindingAnnValue() string {
	return fmt.Sprintf(ClusterRoleBindingName, a.name, a.namespace)
}

func (a *CreatedResourceAnnotations) PackageAnnValue() string {
	return fmt.Sprintf("%s-%s", a.name, a.namespace)
}
