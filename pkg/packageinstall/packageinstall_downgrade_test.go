// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	fakeapiserver "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	fakekappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_PackageInstallVersionDowngrades(t *testing.T) {
	log := logf.Log.WithName("kc")

	pkg1 := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg.1.0.0",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{{
						ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{Image: "ver-1.0.0"},
					}},
				},
			},
		},
	}
	pkg2 := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg.2.0.0",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "2.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{{
						ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{Image: "ver-2.0.0"},
					}},
				},
			},
		},
	}
	existingApp := &v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "instl-pkg",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{{
				ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{Image: "ver-existing"},
			}},
		},
	}

	t.Run("succeeds with same version", func(t *testing.T) {
		pkgInstall := &pkgingv1alpha1.PackageInstall{
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
			Status: pkgingv1alpha1.PackageInstallStatus{
				LastAttemptedVersion: "1.0.0",
			},
		}

		pkgClient := fakeapiserver.NewSimpleClientset(&pkg1, &pkg2)
		appClient := fakekappctrl.NewSimpleClientset(pkgInstall, existingApp)
		coreClient := fake.NewSimpleClientset()

		ip := NewPackageInstallCR(pkgInstall, log, appClient, pkgClient, coreClient)
		_, err := ip.Reconcile()
		assert.Nil(t, err)

		assert.Equal(t, pkgingv1alpha1.PackageInstallStatus{
			Version:              "1.0.0",
			LastAttemptedVersion: "1.0.0",
			GenericStatus: v1alpha1.GenericStatus{
				ObservedGeneration: 0,
				Conditions: []v1alpha1.Condition{{
					Type:   "Reconciling",
					Status: "True",
				}},
				FriendlyDescription: "Reconciling",
			},
		}, getPackageInstall(t, appClient, "instl-pkg").Status)

		assert.Equal(t, "ver-1.0.0", getApp(t, appClient, "instl-pkg").Spec.Fetch[0].ImgpkgBundle.Image)
	})

	t.Run("succeeds with higher version", func(t *testing.T) {
		pkgInstall := &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "expected-pkg",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "2.0.0",
					},
				},
			},
			Status: pkgingv1alpha1.PackageInstallStatus{
				LastAttemptedVersion: "1.0.0",
			},
		}

		pkgClient := fakeapiserver.NewSimpleClientset(&pkg1, &pkg2)
		appClient := fakekappctrl.NewSimpleClientset(pkgInstall, existingApp)
		coreClient := fake.NewSimpleClientset()

		ip := NewPackageInstallCR(pkgInstall, log, appClient, pkgClient, coreClient)
		_, err := ip.Reconcile()
		assert.Nil(t, err)

		assert.Equal(t, pkgingv1alpha1.PackageInstallStatus{
			Version:              "2.0.0",
			LastAttemptedVersion: "2.0.0",
			GenericStatus: v1alpha1.GenericStatus{
				ObservedGeneration: 0,
				Conditions: []v1alpha1.Condition{{
					Type:   "Reconciling",
					Status: "True",
				}},
				FriendlyDescription: "Reconciling",
			},
		}, getPackageInstall(t, appClient, "instl-pkg").Status)

		assert.Equal(t, "ver-2.0.0", getApp(t, appClient, "instl-pkg").Spec.Fetch[0].ImgpkgBundle.Image)
	})

	t.Run("errors when trying to install lower version", func(t *testing.T) {
		pkgInstall := &pkgingv1alpha1.PackageInstall{
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
			Status: pkgingv1alpha1.PackageInstallStatus{
				LastAttemptedVersion: "2.0.0", // higher than 1.0.0
			},
		}

		pkgClient := fakeapiserver.NewSimpleClientset(&pkg1, &pkg2)
		appClient := fakekappctrl.NewSimpleClientset(pkgInstall, existingApp)
		coreClient := fake.NewSimpleClientset()

		ip := NewPackageInstallCR(pkgInstall, log, appClient, pkgClient, coreClient)
		_, err := ip.Reconcile()
		assert.Nil(t, err)

		assert.Equal(t, pkgingv1alpha1.PackageInstallStatus{
			Version:              "1.0.0",
			LastAttemptedVersion: "2.0.0",
			GenericStatus: v1alpha1.GenericStatus{
				ObservedGeneration: 0,
				Conditions: []v1alpha1.Condition{{
					Type:    "ReconcileFailed",
					Status:  "True",
					Reason:  "",
					Message: "Error (see .status.usefulErrorMessage for details)",
				}},
				FriendlyDescription: "Reconcile failed: Error (see .status.usefulErrorMessage for details)",
				UsefulErrorMessage:  "Stopped installing matched version '1.0.0' since last attempted version '2.0.0' is higher.\nhint: Add annotation packaging.carvel.dev/downgradable: \"\" to PackageInstall to proceed with downgrade",
			},
		}, getPackageInstall(t, appClient, "instl-pkg").Status)

		assert.Equal(t, "ver-existing", getApp(t, appClient, "instl-pkg").Spec.Fetch[0].ImgpkgBundle.Image)
	})

	t.Run("succeeds when trying to install lower version but is allowed to downgrade", func(t *testing.T) {
		pkgInstall := &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg",
				Annotations: map[string]string{
					"packaging.carvel.dev/downgradable": "",
				},
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "expected-pkg",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "1.0.0",
					},
				},
			},
			Status: pkgingv1alpha1.PackageInstallStatus{
				LastAttemptedVersion: "2.0.0", // higher than 1.0.0
			},
		}

		pkgClient := fakeapiserver.NewSimpleClientset(&pkg1, &pkg2)
		appClient := fakekappctrl.NewSimpleClientset(pkgInstall, existingApp)
		coreClient := fake.NewSimpleClientset()

		ip := NewPackageInstallCR(pkgInstall, log, appClient, pkgClient, coreClient)
		_, err := ip.Reconcile()
		assert.Nil(t, err)

		assert.Equal(t, pkgingv1alpha1.PackageInstallStatus{
			Version:              "1.0.0",
			LastAttemptedVersion: "1.0.0",
			GenericStatus: v1alpha1.GenericStatus{
				ObservedGeneration: 0,
				Conditions: []v1alpha1.Condition{{
					Type:   "Reconciling",
					Status: "True",
				}},
				FriendlyDescription: "Reconciling",
			},
		}, getPackageInstall(t, appClient, "instl-pkg").Status)

		assert.Equal(t, "ver-1.0.0", getApp(t, appClient, "instl-pkg").Spec.Fetch[0].ImgpkgBundle.Image)
	})
}

func getPackageInstall(t *testing.T, clientset *fakekappctrl.Clientset, name string) *pkgingv1alpha1.PackageInstall {
	gvr := schema.GroupVersionResource{"packaging.carvel.dev", "v1alpha1", "packageinstalls"}
	obj, err := clientset.Tracker().Get(gvr, "", name)
	assert.Nil(t, err)
	require.NotNil(t, obj)
	return obj.(*pkgingv1alpha1.PackageInstall)
}

func getApp(t *testing.T, clientset *fakekappctrl.Clientset, name string) *v1alpha1.App {
	gvr := schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err := clientset.Tracker().Get(gvr, "", name)
	assert.Nil(t, err)
	require.NotNil(t, obj)
	return obj.(*v1alpha1.App)
}
