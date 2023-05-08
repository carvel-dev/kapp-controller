// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"testing"

	"github.com/k14s/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	fakeapiserver "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	fakekappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_PackageInstallDeletion(t *testing.T) {
	log := logf.Log.WithName("kc")
	now := metav1.Now()

	t.Run("ensures that deletion related fields on App are kept in sync with PackageInstall "+
		"even after PackageInstall is marked for deletion", func(t *testing.T) {
		existingApp := &v1alpha1.App{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "instl-pkg",
				DeletionTimestamp: &now, // being deleted
			},
			Spec: v1alpha1.AppSpec{
				ServiceAccountName: "",
				Cluster:            nil,
				NoopDelete:         false,
				Paused:             false,
				Canceled:           false,
				Fetch: []v1alpha1.AppFetch{{
					ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{Image: "ver-existing"},
				}},
			},
		}

		pkgInstall := &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "instl-pkg",
				DeletionTimestamp: &now, // being deleted
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				ServiceAccountName: "sa-name",
				Cluster:            &v1alpha1.AppCluster{},
				NoopDelete:         true,
				Paused:             true,
				Canceled:           true,
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "expected-pkg",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "1.0.0",
					},
				},
			},
		}

		pkgClient := fakeapiserver.NewSimpleClientset()
		appClient := fakekappctrl.NewSimpleClientset(pkgInstall, existingApp)
		coreClient := fake.NewSimpleClientset()

		ip := NewPackageInstallCR(pkgInstall, log, appClient, pkgClient, coreClient, FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")})
		_, err := ip.Reconcile()
		assert.Nil(t, err)

		app := getApp(t, appClient, "instl-pkg")
		assert.Equal(t, "sa-name", app.Spec.ServiceAccountName)
		assert.Equal(t, &v1alpha1.AppCluster{}, app.Spec.Cluster)
		assert.Equal(t, true, app.Spec.NoopDelete)
		assert.Equal(t, true, app.Spec.Paused)
		assert.Equal(t, true, app.Spec.Canceled)
	})
}

type FakeComponentInfo struct {
	KCVersion  semver.Version
	K8sVersion semver.Version
}

func (f FakeComponentInfo) KappControllerVersion() (semver.Version, error) {
	return f.KCVersion, nil
}

func (f FakeComponentInfo) KubernetesVersion(_ string, _ *v1alpha1.AppCluster, _ *metav1.ObjectMeta) (semver.Version, error) {
	return f.K8sVersion, nil
}
