// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	fakeapiserver "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	fakekappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
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
	fakePkgClient := fakeapiserver.NewSimpleClientset(&expectedPackageVersion)

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
	fakePkgClient := fakeapiserver.NewSimpleClientset(&expectedPackageVersion, &alternatePackageVersion)

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

func Test_PlaceHolderSecretCreated_WhenPackageHasNoSecretRef(t *testing.T) {
	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{
						{
							ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
								// Since no secretRef, we expect a placeholder secret
								// to be created by kapp-controller.
								Image: "foo/bar",
							},
						},
					},
				},
			},
		},
	}

	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	model := &pkgingv1alpha1.PackageInstall{
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
	}
	log := logf.Log.WithName("kc")
	fakekctrl := fakekappctrl.NewSimpleClientset(model)
	fakek8s := fake.NewSimpleClientset()
	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s)

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch0")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	secret := obj.(*corev1.Secret)
	_, ok := secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)

	gvr = schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err = fakekctrl.Tracker().Get(gvr, "", "instl-pkg")
	require.NotNil(t, obj)
	assert.Nil(t, err)
	app := obj.(*v1alpha1.App)

	assert.Equal(t, 1, len(app.Spec.Fetch))
	assert.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle.SecretRef, "expected ImgpkgBundle secretRef to be non nil but was nil")

	assert.Equal(t, "instl-pkg-fetch0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
}

func Test_PlaceHolderSecretsCreated_WhenPackageHasMultipleFetchStages(t *testing.T) {
	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{
						{
							ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
								// Since no secretRef, we expect a placeholder secret
								// to be created by kapp-controller.
								Image: "foo/bar",
							},
						},
						{
							Image: &v1alpha1.AppFetchImage{
								URL: "foo/bar",
							},
						},
					},
				},
			},
		},
	}

	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	model := &pkgingv1alpha1.PackageInstall{
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
	}
	log := logf.Log.WithName("kc")
	fakekctrl := fakekappctrl.NewSimpleClientset(model)
	fakek8s := fake.NewSimpleClientset()
	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s)

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch0")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	secret := obj.(*corev1.Secret)
	_, ok := secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)

	gvr = schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err = fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch1")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	secret = obj.(*corev1.Secret)
	_, ok = secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)

	gvr = schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err = fakekctrl.Tracker().Get(gvr, "", "instl-pkg")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	app := obj.(*v1alpha1.App)

	assert.Equal(t, 2, len(app.Spec.Fetch))
	assert.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle.SecretRef, "expected ImgpkgBundle secretRef to be non nil but was nil")
	assert.NotNil(t, app.Spec.Fetch[1].Image.SecretRef, "expected Image secretRef to be non nil but was nil")

	assert.Equal(t, "instl-pkg-fetch0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
	assert.Equal(t, "instl-pkg-fetch1", app.Spec.Fetch[1].Image.SecretRef.Name)
}

func Test_PlaceHolderSecretsNotCreated_WhenFetchStagesHaveSecrets(t *testing.T) {
	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{
						{
							ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
								Image:     "foo/bar",
								SecretRef: &v1alpha1.AppFetchLocalRef{"foo"},
							},
						},
						{
							Image: &v1alpha1.AppFetchImage{
								URL:       "foo/bar",
								SecretRef: &v1alpha1.AppFetchLocalRef{"foo1"},
							},
						},
					},
				},
			},
		},
	}

	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	model := &pkgingv1alpha1.PackageInstall{
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
	}
	log := logf.Log.WithName("kc")
	fakekctrl := fakekappctrl.NewSimpleClientset(model)
	fakek8s := fake.NewSimpleClientset()
	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s)

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err := fakekctrl.Tracker().Get(gvr, "", "instl-pkg")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	app := obj.(*v1alpha1.App)

	assert.Equal(t, 2, len(app.Spec.Fetch))
	assert.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle.SecretRef, "expected ImgpkgBundle secretRef to be non nil but was nil")
	assert.NotNil(t, app.Spec.Fetch[1].Image.SecretRef, "expected Image secretRef to be non nil but was nil")

	assert.Equal(t, "foo", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
	assert.Equal(t, "foo1", app.Spec.Fetch[1].Image.SecretRef.Name)
}

func Test_PlaceHolderSecretCreated_WhenPackageInstallUpdated(t *testing.T) {
	appSpec := v1alpha1.AppSpec{
		Fetch: []v1alpha1.AppFetch{
			{
				ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
					// Since no secretRef, we expect a placeholder secret
					// to be created by kapp-controller.
					Image: "foo/bar",
				},
			},
		},
	}

	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &appSpec,
			},
		},
	}

	// Load package into fake client
	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	model := &pkgingv1alpha1.PackageInstall{
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
	}
	log := logf.Log.WithName("kc")

	// The existing App in this test should have a secret
	// so when the App is updated based on Package definition we should
	// see the placeholder secret used by the App instead of older
	// secret.
	appSpec.Fetch[0].ImgpkgBundle.SecretRef = &v1alpha1.AppFetchLocalRef{"secret-update"}
	existingApp := &v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "instl-pkg",
		},
		Spec: appSpec,
	}

	fakekctrl := fakekappctrl.NewSimpleClientset(model, existingApp)
	fakek8s := fake.NewSimpleClientset()
	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s)

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch0")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	secret := obj.(*corev1.Secret)
	_, ok := secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)

	gvr = schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err = fakekctrl.Tracker().Get(gvr, "", "instl-pkg")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	app := obj.(*v1alpha1.App)

	assert.Equal(t, 1, len(app.Spec.Fetch))
	assert.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle.SecretRef, "expected ImgpkgBundle secretRef to be non nil but was nil")

	assert.Equal(t, "instl-pkg-fetch0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
}

func Test_PlaceHolderSecretUpdated_WhenPlaceHolderSecretAlreadyExists(t *testing.T) {
	appSpec := v1alpha1.AppSpec{
		Fetch: []v1alpha1.AppFetch{
			{
				ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
					// Since no secretRef, we expect a placeholder secret
					// to be created by kapp-controller.
					Image: "foo/bar",
				},
			},
		},
	}

	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "expected-pkg",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expected-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &appSpec,
			},
		},
	}

	// Load package into fake client
	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	model := &pkgingv1alpha1.PackageInstall{
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
	}
	log := logf.Log.WithName("kc")

	existingApp := &v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "instl-pkg",
		},
		Spec: appSpec,
	}

	fakekctrl := fakekappctrl.NewSimpleClientset(model, existingApp)

	placeHolderSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "instl-pkg-fetch0",
			Annotations: map[string]string{
				"secretgen.carvel.dev/image-pull-secret": "",
			},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: []byte("dXNlcm5hbWU6cGFzc3dvcmQ="),
		},
	}
	fakek8s := fake.NewSimpleClientset(&placeHolderSecret)

	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s)

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch0")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	secret := obj.(*corev1.Secret)
	_, ok := secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)
	// Assert that key is updated to empty auth
	auths, ok := secret.Data[corev1.DockerConfigJsonKey]
	assert.Equal(t, `{"auths":{}}`, string(auths))

	gvr = schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err = fakekctrl.Tracker().Get(gvr, "", "instl-pkg")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	app := obj.(*v1alpha1.App)

	assert.Equal(t, 1, len(app.Spec.Fetch))
	assert.NotNil(t, app.Spec.Fetch[0].ImgpkgBundle.SecretRef, "expected ImgpkgBundle secretRef to be non nil but was nil")

	assert.Equal(t, "instl-pkg-fetch0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
}
