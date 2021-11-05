// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package install

import "fmt"

const (
	// For use with packages created using kapp-ctrl-cli
	KappPkgAnnotation       = "packaging.carvel.dev/package"
	KappPkgAnnotationPrefix = "package"

	// For use with packages created with Tanzu CLI
	// KappPkgAnnotation       = "tkg.tanzu.vmware.com/tanzu-package"
	// KappPkgAnnotationPrefix = "tanzu-package"

	ClusterRoleBindingName = "%s-%s-cluster-rolebinding"
	ClusterRoleName        = "%s-%s-cluster-role"
	SecretName             = "%s-%s-values"
	ServiceAccountName     = "%s-%s-sa"
)

type CreatedResourceAnnotations struct {
	name      string
	namespace string
}

func NewCreatedResourceAnnotations(name string, namepsace string) *CreatedResourceAnnotations {
	return &CreatedResourceAnnotations{
		name:      name,
		namespace: namepsace,
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
