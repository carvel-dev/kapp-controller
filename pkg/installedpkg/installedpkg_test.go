// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg

import (
	"reflect"
	"testing"

	ipv1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	packagev1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/packages/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This test was developed for issue:
// https://github.com/vmware-tanzu/carvel-kapp-controller/issues/116
func Test_PackageRefWithPrerelease_IsFound(t *testing.T) {
	// Package with prerelease version
	expectedPackageVersion := packagev1.PackageVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkg.test.carvel.dev",
		},
		Spec: packagev1.PackageVersionSpec{
			PackageName: "pkg.test.carvel.dev",
			Version:     "3.0.0-rc.1",
		},
	}

	// Load package into fake client
	fakePkgClient := fake.NewSimpleClientset(&expectedPackageVersion)

	// InstalledPackage that has PackageRef with prerelease
	ip := InstalledPackageCR{
		model: &ipv1.InstalledPackage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-prerelease",
			},
			Spec: ipv1.InstalledPackageSpec{
				PackageVersionRef: &ipv1.PackageVersionRef{
					PackageName: "pkg.test.carvel.dev",
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
	// Package with prerelease version
	expectedPackageVersion := packagev1.PackageVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: packagev1.PackageVersionSpec{
			PackageName: "expected-pkg",
			Version:     "1.0.0",
		},
	}

	alternatePackageVersion := packagev1.PackageVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "alternate-pkg",
		},
		Spec: packagev1.PackageVersionSpec{
			PackageName: "alternate-pkg",
			Version:     "1.0.0",
		},
	}

	// Load package into fake client
	fakePkgClient := fake.NewSimpleClientset(&expectedPackageVersion, &alternatePackageVersion)

	// InstalledPackage that has PackageRef with prerelease
	ip := InstalledPackageCR{
		model: &ipv1.InstalledPackage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg",
			},
			Spec: ipv1.InstalledPackageSpec{
				PackageVersionRef: &ipv1.PackageVersionRef{
					PackageName: "expected-pkg",
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
