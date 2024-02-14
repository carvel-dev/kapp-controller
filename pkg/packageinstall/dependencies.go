// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"

	packagingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	versionsv1alpha1 "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (pi *PackageInstallCR) reconcileDependencies(pkg v1alpha1.Package) error {
	var dependencyList []v1alpha1.Package
	// Check if packages exist in the cluster
	for _, dep := range pkg.Spec.Dependencies {

		if dep.Package != nil {
			pkg, err := pi.getPkgVersion(dep.Package.RefName, &versionsv1alpha1.VersionSelectionSemver{
				Constraints: dep.Package.Version,
			})
			if err != nil {
				if errors.IsNotFound(err) {
					continue
				}
				return err
			}
			dependencyList = append(dependencyList, pkg)
		}

	}

	// Check if the packages are installed
	for _, dep := range dependencyList {
		pkgiName := dep.Spec.RefName + "." + dep.Spec.Version
		pkgi, err := pi.kcclient.PackagingV1alpha1().PackageInstalls(pi.model.Namespace).Get(context.TODO(), pkgiName, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				// do we need to ignore this and continue?
				return err
			}

			pkgi = &packagingv1alpha1.PackageInstall{
				TypeMeta: pi.model.TypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      pkgiName,
					Namespace: pi.model.Namespace,
					Annotations: map[string]string{
						"kapp-controller.carvel.dev/owner": "PackageInstall/" + pi.model.Name,
					},
				},
				Spec: packagingv1alpha1.PackageInstallSpec{
					ServiceAccountName: pi.model.Spec.ServiceAccountName,
					PackageRef: &packagingv1alpha1.PackageRef{
						RefName: dep.Spec.RefName,
						VersionSelection: &versionsv1alpha1.VersionSelectionSemver{
							Constraints: dep.Spec.Version,
						},
					},
					DefaultNamespace: pi.model.Spec.DefaultNamespace,
					Dependencies:     packagingv1alpha1.Dependencies{},
				},
			}

			_, err = pi.kcclient.PackagingV1alpha1().PackageInstalls(pi.model.Namespace).Create(context.TODO(), pkgi, metav1.CreateOptions{})
			if err != nil {
				// do we need to ignore this and continue?
				return fmt.Errorf("Unable to create the packageinstall for the package %s: %s", dep.Name, err)
			}
		}

	}
	return nil
}
