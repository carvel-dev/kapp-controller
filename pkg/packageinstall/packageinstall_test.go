// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"reflect"
	"testing"

	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This test was developed for issue:
// https://github.com/vmware-tanzu/carvel-kapp-controller/issues/116
func Test_PackageRefWithPrerelease_IsFound(t *testing.T) {
	// PackageMetadata with prerelease version
	expectedPackageVersion := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkg.test.carvel.dev",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "pkg.test.carvel.dev",
			Version: "3.0.0-rc.1",
		},
	}

	// Load package into fake client
	fakePkgClient := fake.NewSimpleClientset(&expectedPackageVersion)

	// PackageInstall that has PackageRef with prerelease
	ip := PackageInstallCR{
		model: &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-prerelease",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "pkg.test.carvel.dev",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "3.0.0-rc.1",
						Prereleases: &versions.VersionSelectionSemverPrereleases{
							Identifiers: []string{"rc"},
						},
					},
				},
			},
		},
		pkgclient: fakePkgClient,
	}

	out, err := ip.referencedPkgVersion()
	if err != nil {
		t.Fatalf("\nExpected no error from getting PackageRef with prerelease\nBut got:\n%v\n", err)
	}

	if !reflect.DeepEqual(out, expectedPackageVersion) {
		t.Fatalf("\nPackageVersion is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedPackageVersion, out)
	}
}

func Test_PackageRefUsesName(t *testing.T) {
	// PackageMetadata with prerelease version
	expectedPackageVersion := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
		},
	}

	alternatePackageVersion := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "alternate-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "alternate-pkg",
			Version: "1.0.0",
		},
	}

	// Load package into fake client
	fakePkgClient := fake.NewSimpleClientset(&expectedPackageVersion, &alternatePackageVersion)

	// PackageInstall that has PackageRef with prerelease
	ip := PackageInstallCR{
		model: &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "expected-pkg",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "1.0.0",
					},
				},
			},
		},
		pkgclient: fakePkgClient,
	}

	out, err := ip.referencedPkgVersion()
	if err != nil {
		t.Fatalf("\nExpected no error from resolving referenced package\nBut got:\n%v\n", err)
	}

	if !reflect.DeepEqual(out, expectedPackageVersion) {
		t.Fatalf("\nPackageVersion is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedPackageVersion, out)
	}
}
