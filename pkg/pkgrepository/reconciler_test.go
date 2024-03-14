// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	fakekappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/pkgrepository"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_PlaceholderSecretCreated_WhenPackageRepositoryHasNoSecret(t *testing.T) {
	pkgr := &v1alpha1.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkgrepo-with-placeholder",
		},
		Spec: v1alpha1.PackageRepositorySpec{
			Fetch: &v1alpha1.PackageRepositoryFetch{
				ImgpkgBundle: &v1alpha12.AppFetchImgpkgBundle{
					Image: "repository/user/image",
				},
			},
		},
		Status: v1alpha1.PackageRepositoryStatus{},
	}

	fakekctrl := fakekappctrl.NewSimpleClientset(pkgr)
	fakek8s := fake.NewSimpleClientset()
	log := logf.Log.WithName("kc")

	pkgri := pkgrepository.NewReconciler(fakekctrl, fakek8s,
		log, pkgrepository.AppFactory{}, nil, nil)

	app, err := pkgrepository.NewPackageRepoApp(pkgr)
	assert.Nil(t, err, "error from creating PackageRepository App: %s", err)

	err = pkgri.ReconcileFetchPlaceholderSecrets(*pkgr, app)
	assert.Nil(t, err)

	placeholderSecretName := pkgr.Name + "-fetch-0"
	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", placeholderSecretName)
	assert.Nil(t, err, "error from checking placeholder secret exists: %s", err)
	require.NotNil(t, obj)
	secret := obj.(*corev1.Secret)
	_, ok := secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)

	require.Equal(t, 1, len(app.Spec.Fetch))
	require.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle)
	assert.Equal(t, placeholderSecretName, app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
}

func Test_PlaceholderSecretNotCreated_WhenPackageRepositoryHasSecret(t *testing.T) {
	pkgr := &v1alpha1.PackageRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkgrepo-with-placeholder",
		},
		Spec: v1alpha1.PackageRepositorySpec{
			Fetch: &v1alpha1.PackageRepositoryFetch{
				ImgpkgBundle: &v1alpha12.AppFetchImgpkgBundle{
					Image:     "repository/user/image",
					SecretRef: &v1alpha12.AppFetchLocalRef{Name: "secret"},
				},
			},
		},
		Status: v1alpha1.PackageRepositoryStatus{},
	}

	fakekctrl := fakekappctrl.NewSimpleClientset(pkgr)
	fakek8s := fake.NewSimpleClientset()
	log := logf.Log.WithName("kc")

	pkgri := pkgrepository.NewReconciler(fakekctrl, fakek8s,
		log, pkgrepository.AppFactory{}, nil, nil)

	app, err := pkgrepository.NewPackageRepoApp(pkgr)
	assert.Nil(t, err, "error from creating PackageRepository App: %s", err)

	err = pkgri.ReconcileFetchPlaceholderSecrets(*pkgr, app)
	assert.Nil(t, err)

	placeholderSecretName := pkgr.Name + "-fetch-0"
	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	_, err = fakek8s.Tracker().Get(gvr, "", placeholderSecretName)
	assert.NotNil(t, err, "error from checking placeholder secret exists: %s", err)
	assert.True(t, errors.IsNotFound(err))

	require.Equal(t, 1, len(app.Spec.Fetch))
	require.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle)
	assert.Equal(t, "secret", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
}
