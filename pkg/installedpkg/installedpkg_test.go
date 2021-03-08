// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg

import (
	"reflect"
	"testing"

	ipv1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	packagev1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/package/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This test was developed for issue:
// https://github.com/vmware-tanzu/carvel-kapp-controller/issues/116
func Test_PackageRefWithPrerelease_IsFound(t *testing.T) {
	// Package with prerelease version
	expectedPkg := packagev1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkg.test.carvel.dev",
		},
		Spec: packagev1.PackageSpec{
			PublicName: "pkg.test.carvel.dev",
			Version:    "3.0.0-rc.1",
		},
	}

	// Load package into fake client
	kappcs := fake.NewSimpleClientset(&expectedPkg)

	// InstalledPackage that has PackageRef with prerelease
	ip := InstalledPackageCR{
		model: &ipv1.InstalledPackage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-prerelease",
			},
			Spec: ipv1.InstalledPackageSpec{
				PkgRef: &ipv1.PackageRef{
					PublicName: "pkg.test.carvel.dev",
					Version:    "3.0.0-rc.1",
				},
			},
		},
		client: kappcs,
	}

	out, err := ip.referencedPkg()
	if err != nil {
		t.Fatalf("\nExpected no error from getting PackageRef with prerelease\nBut got:\n%v\n", err)
	}

	if !reflect.DeepEqual(out, expectedPkg) {
		t.Fatalf("\nPackage is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedPkg, out)
	}
}

func Test_PackageRefUsesName(t *testing.T) {
	// Package with prerelease version
	expectedPkg := packagev1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: packagev1.PackageSpec{
			PublicName: "expected-pkg",
			Version:    "1.0.0",
		},
	}

	alternatePkg := packagev1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "alternate-pkg",
		},
		Spec: packagev1.PackageSpec{
			PublicName: "alternate-pkg",
			Version:    "1.0.0",
		},
	}

	// Load package into fake client
	kappcs := fake.NewSimpleClientset(&expectedPkg, &alternatePkg)

	// InstalledPackage that has PackageRef with prerelease
	ip := InstalledPackageCR{
		model: &ipv1.InstalledPackage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg",
			},
			Spec: ipv1.InstalledPackageSpec{
				PkgRef: &ipv1.PackageRef{
					PublicName: "expected-pkg",
					Version:    "1.0.0",
				},
			},
		},
		client: kappcs,
	}

	out, err := ip.referencedPkg()
	if err != nil {
		t.Fatalf("\nExpected no error from resolving referenced package\nBut got:\n%v\n", err)
	}

	if !reflect.DeepEqual(out, expectedPkg) {
		t.Fatalf("\nPackage is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedPkg, out)
	}
}
