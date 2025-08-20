// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"fmt"
	"reflect"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	fakeapiserver "carvel.dev/kapp-controller/pkg/apiserver/client/clientset/versioned/fake"
	fakekappctrl "carvel.dev/kapp-controller/pkg/client/clientset/versioned/fake"
	"carvel.dev/kapp-controller/pkg/metrics"
	versions "carvel.dev/vendir/pkg/vendir/versions/v1alpha1"
	"github.com/k14s/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// This test was developed for issue:
// https://carvel.dev/kapp-controller/issues/116
func Test_PackageRefWithPrerelease_IsFound(t *testing.T) {
	log := logf.Log.WithName("kc")

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
				ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
			},
		},
		pkgclient: fakePkgClient,
		log:       log,
		compInfo:  FakeComponentInfo{K8sVersion: semver.MustParse("0.20.0")},
	}

	out, err := ip.referencedPkgVersion()
	if err != nil {
		t.Fatalf("\nExpected no error from getting PackageRef with prerelease\nBut got:\n%v\n", err)
	}

	if !reflect.DeepEqual(out, expectedPackageVersion) {
		t.Fatalf("\nPackageVersion is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedPackageVersion, out)
	}
}

func Test_PackageWithConstraints(t *testing.T) {
	const (
		kubernetesVersionOverrideAnnotation     = "packaging.carvel.dev/ignore-kubernetes-version-selection"
		kappControllerVersionOverrideAnnotation = "packaging.carvel.dev/ignore-kapp-controller-version-selection"
	)

	log := logf.Log.WithName("kc")
	fakek8s := fake.NewSimpleClientset()
	pkg := generatePackageWithConstraints("pkg.test.carvel.dev", "0.0.0", ">1.0.0 <2.0.0", ">0.15.0")
	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	ip := PackageInstallCR{
		model: &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-ignore-kc-constraint",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "pkg.test.carvel.dev",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "0.0.0",
					},
				},
				ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
			},
		},
		pkgclient:  fakePkgClient,
		log:        log,
		coreClient: fakek8s,
		compInfo:   FakeComponentInfo{KCVersion: semver.MustParse("1.5.0"), K8sVersion: semver.MustParse("0.20.0")},
	}

	// all constraints met
	_, err := ip.referencedPkgVersion()
	require.NoError(t, err)

	// kapp-controller version constraint fail
	ip.compInfo = FakeComponentInfo{KCVersion: semver.MustParse("3.0.0"), K8sVersion: semver.MustParse("0.20.0")}
	_, err = ip.referencedPkgVersion()
	require.Error(t, err)
	assert.ErrorContains(t, err, "after-kubernetes-version-check=1")
	assert.ErrorContains(t, err, "after-kapp-controller-version-check=0")

	// kapp-controller version override annotation
	ip.model.ObjectMeta.Annotations = map[string]string{
		kappControllerVersionOverrideAnnotation: "",
	}
	_, err = ip.referencedPkgVersion()
	require.NoError(t, err)

	// kubernetes version constraint fail
	ip.compInfo = FakeComponentInfo{KCVersion: semver.MustParse("1.5.0"), K8sVersion: semver.MustParse("0.0.0")}
	_, err = ip.referencedPkgVersion()
	require.Error(t, err)
	assert.ErrorContains(t, err, "after-kubernetes-version-check=0")

	// kubernetes version override annotation
	ip.model.ObjectMeta.Annotations = map[string]string{
		kappControllerVersionOverrideAnnotation: "",
		kubernetesVersionOverrideAnnotation:     "",
	}
	_, err = ip.referencedPkgVersion()
	require.NoError(t, err)
}

func Test_Package_NotFound(t *testing.T) {
	log := logf.Log.WithName("kc")
	fakek8s := fake.NewSimpleClientset()
	fakePkgClient := fakeapiserver.NewSimpleClientset()
	pkgName := "pkg.test.carvel.dev"

	ip := PackageInstallCR{
		model: &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-no-pkg-found",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: pkgName,
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "0.0.0",
					},
				},
			},
		},
		pkgclient:  fakePkgClient,
		compInfo:   FakeComponentInfo{KCVersion: semver.MustParse("0.42.0")},
		log:        log,
		coreClient: fakek8s,
	}

	_, err := ip.referencedPkgVersion()
	require.Error(t, err)
	assert.ErrorContains(t, err, fmt.Sprintf("Package %s not found", pkgName))
}

func Test_Package_ConstraintNotGiven_ErrorDoesNotContainMessage(t *testing.T) {
	log := logf.Log.WithName("kc")
	fakek8s := fake.NewSimpleClientset()
	pkg := generatePackageWithConstraints("pkg.test.carvel.dev", "0.0.0", "1.0.0", "")
	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	ip := PackageInstallCR{
		model: &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-ignore-kc-constraint",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: "pkg.test.carvel.dev",
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: "0.0.0",
					},
				},
				ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
			},
		},
		pkgclient:  fakePkgClient,
		compInfo:   FakeComponentInfo{KCVersion: semver.MustParse("1.5.0")},
		log:        log,
		coreClient: fakek8s,
	}

	_, err := ip.referencedPkgVersion()
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "after-kubernetes-version-check=1")
	assert.ErrorContains(t, err, "after-kapp-controller-version-check=0")
}

func Test_PackageWithConstraints_HighestMatch(t *testing.T) {
	log := logf.Log.WithName("kc")
	fakek8s := fake.NewSimpleClientset()
	pkgName := "pkg.test.carvel.dev"
	pkg1 := generatePackageWithConstraints(pkgName, "0.4.0", ">0.1.0", ">0.1.0") // this one is the lowest version but installable
	pkg2 := generatePackageWithConstraints(pkgName, "0.5.0", ">0.1.0", ">0.1.0") // this one is the highest installable version
	pkg3 := generatePackageWithConstraints(pkgName, "1.4.1", ">2.0.0", "")       // higher version uninstallable
	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg1, &pkg2, &pkg3)

	ip := PackageInstallCR{
		model: &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name: "instl-pkg-multi-version-constraints",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				PackageRef: &pkgingv1alpha1.PackageRef{
					RefName: pkgName,
					VersionSelection: &versions.VersionSelectionSemver{
						Constraints: ">0.0.0",
					},
				},
				ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
			},
		},
		pkgclient:  fakePkgClient,
		compInfo:   FakeComponentInfo{KCVersion: semver.MustParse("1.5.0"), K8sVersion: semver.MustParse("0.20.0")},
		log:        log,
		coreClient: fakek8s,
	}

	out, err := ip.referencedPkgVersion()
	assert.Equal(t, out, pkg2, "Highest version of Package meeting constraints not chosen: \nExpected:\n%#v\nGot:\n%#v\n", pkg2, out)
	require.NoError(t, err)
}

func Test_PackageRefWithPrerelease_DoesNotRequirePrereleaseMarker(t *testing.T) {
	expectedPackageVersion := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pkg.test.carvel.dev",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "pkg.test.carvel.dev",
			Version: "3.0.0-rc.1",
		},
	}

	log := logf.Log.WithName("kc")
	fakek8s := fake.NewSimpleClientset()
	fakePkgClient := fakeapiserver.NewSimpleClientset(&expectedPackageVersion)

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

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
					},
				},
				ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
			},
		},
		pkgclient: fakePkgClient,
		compInfo:  FakeComponentInfo{KCVersion: semver.MustParse("1.5.0")},
		log:       log,
	}

	out, err := ip.referencedPkgVersion()
	require.NoError(t, err)
	require.Equal(t, out, expectedPackageVersion)
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
	log := logf.Log.WithName("kc")
	fakek8s := fake.NewSimpleClientset()

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

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
				ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
			},
		},
		pkgclient: fakePkgClient,
		log:       log,
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
			ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
		},
	}
	log := logf.Log.WithName("kc")
	fakekctrl := fakekappctrl.NewSimpleClientset(model)
	fakek8s := fake.NewSimpleClientset()

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s,
		FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")}, Opts{},
		metrics.NewMetrics())

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch-0")
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

	assert.Equal(t, "instl-pkg-fetch-0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
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
			ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
		},
	}
	log := logf.Log.WithName("kc")
	fakekctrl := fakekappctrl.NewSimpleClientset(model)
	fakek8s := fake.NewSimpleClientset()
	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s,
		FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")}, Opts{},
		metrics.NewMetrics())

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch-0")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	secret := obj.(*corev1.Secret)
	_, ok := secret.Annotations["secretgen.carvel.dev/image-pull-secret"]
	assert.True(t, ok)

	gvr = schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err = fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch-1")
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

	assert.Equal(t, "instl-pkg-fetch-0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
	assert.Equal(t, "instl-pkg-fetch-1", app.Spec.Fetch[1].Image.SecretRef.Name)
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
			ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
		},
	}
	log := logf.Log.WithName("kc")
	fakekctrl := fakekappctrl.NewSimpleClientset(model)
	fakek8s := fake.NewSimpleClientset()

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s,
		FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")}, Opts{},
		metrics.NewMetrics())

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
			ServiceAccountName: "use-local-cluster-sa", // saname being present indicates use local cluster version
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
	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s,
		FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")}, Opts{},
		metrics.NewMetrics())

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

	_, err := ip.Reconcile()
	assert.Nil(t, err)

	gvr := schema.GroupVersionResource{"", "v1", "secrets"}
	obj, err := fakek8s.Tracker().Get(gvr, "", "instl-pkg-fetch-0")
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

	assert.Equal(t, "instl-pkg-fetch-0", app.Spec.Fetch[0].ImgpkgBundle.SecretRef.Name)
}

// Test_StatusUpdaterClosureWithAppUpdateFromReconcileFailed tests this exact scenario:
// 1. App exists in ReconcileFailed state (generation == observedGeneration)
// 2. PackageInstall gets updated to reference a new package version
// 3. This updates the App spec (and in real K8s would increment generation)
// 4. PackageInstall status must now reflect the updated App state
func Test_StatusUpdaterClosureWithAppUpdateFromReconcileFailed(t *testing.T) {
	log := logf.Log.WithName("kc")

	// Create a package
	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pkg.2.0.0",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "test-pkg",
			Version: "2.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{
						{
							ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
								Image: "test-pkg:2.0.0",
							},
						},
					},
					Template: []v1alpha1.AppTemplate{
						{
							Ytt: &v1alpha1.AppTemplateYtt{},
						},
					},
					Deploy: []v1alpha1.AppDeploy{
						{
							Kapp: &v1alpha1.AppDeployKapp{},
						},
					},
				},
			},
		},
	}

	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	model := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pkg-install",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			PackageRef: &pkgingv1alpha1.PackageRef{
				RefName: "test-pkg",
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: "2.0.0",
				},
			},
			ServiceAccountName: "use-local-cluster-sa",
		},
	}

	// Create an existing App that will be updated
	existingApp := &v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-pkg-install",
			Generation: 3,
			Annotations: map[string]string{
				"packaging.carvel.dev/package-ref-name": "test-pkg",
				"packaging.carvel.dev/package-version":  "1.0.0", // Old version
			},
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				{
					ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
						Image: "test-pkg:1.0.0", // Old image
					},
				},
			},
			Template: []v1alpha1.AppTemplate{
				{
					Ytt: &v1alpha1.AppTemplateYtt{},
				},
			},
			Deploy: []v1alpha1.AppDeploy{
				{
					Kapp: &v1alpha1.AppDeployKapp{},
				},
			},
			ServiceAccountName: "use-local-cluster-sa",
		},
		Status: v1alpha1.AppStatus{
			GenericStatus: v1alpha1.GenericStatus{
				ObservedGeneration: 3, // Generation == ObservedGeneration
				Conditions: []v1alpha1.Condition{
					{
						Type:   v1alpha1.ReconcileFailed,
						Status: corev1.ConditionTrue,
					},
				},
				FriendlyDescription: "Reconcile failed",
				UsefulErrorMessage:  "Original error from v1.0.0",
			},
		},
	}

	// Create a custom fake client that increments generation on update (like real K8s)
	fakekctrl := fakekappctrl.NewSimpleClientset(model, existingApp)

	// Add a reactor to simulate generation increment on App updates
	fakekctrl.PrependReactor("update", "apps", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		updateAction := action.(k8stesting.UpdateAction)
		if app, ok := updateAction.GetObject().(*v1alpha1.App); ok {
			// Simulate generation increment when App is updated (like real K8s)
			app.Generation = app.Generation + 1
		}
		return false, nil, nil // Let the default handler process the update
	})

	fakek8s := fake.NewSimpleClientset()
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s,
		FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")}, Opts{},
		metrics.NewMetrics())

	// Reconcile should update the App and set PackageInstall status
	_, err := ip.Reconcile()
	assert.Nil(t, err)

	// Verify the App was updated
	gvr := schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err := fakekctrl.Tracker().Get(gvr, "", "test-pkg-install")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	updatedApp := obj.(*v1alpha1.App)

	// Verify App was updated with new package content
	assert.Equal(t, "test-pkg:2.0.0", updatedApp.Spec.Fetch[0].ImgpkgBundle.Image)
	assert.Equal(t, "2.0.0", updatedApp.Annotations["packaging.carvel.dev/package-version"])
	assert.Equal(t, int64(4), updatedApp.Generation, "Generation should have incremented from 3 to 4")

	// Verify PackageInstall status reflects the updated App state
	assert.Len(t, ip.model.Status.Conditions, 1)
	assert.Equal(t, v1alpha1.Reconciling, ip.model.Status.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, ip.model.Status.Conditions[0].Status)
	assert.Equal(t, "Reconciling", ip.model.Status.FriendlyDescription)
}

// Test_StatusUpdaterClosureWithNoAppUpdate verifies the closure works when App doesn't need updating:
// 1. App exists in ReconcileFailed state (generation == observedGeneration)
// 2. PackageInstall points to same package version as App
// 3. No App update needed, so PackageInstall status should reflect existing App state
func Test_StatusUpdaterClosureWithNoAppUpdate(t *testing.T) {
	log := logf.Log.WithName("kc")

	// Create a package
	pkg := datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pkg.1.0.0",
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "test-pkg",
			Version: "1.0.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &v1alpha1.AppSpec{
					Fetch: []v1alpha1.AppFetch{
						{
							ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
								Image: "test-pkg:1.0.0",
							},
						},
					},
					Template: []v1alpha1.AppTemplate{
						{
							Ytt: &v1alpha1.AppTemplateYtt{},
						},
					},
					Deploy: []v1alpha1.AppDeploy{
						{
							Kapp: &v1alpha1.AppDeployKapp{},
						},
					},
				},
			},
		},
	}

	fakePkgClient := fakeapiserver.NewSimpleClientset(&pkg)

	// Create a PackageInstall pointing to v1.0.0
	model := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pkg-install",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			PackageRef: &pkgingv1alpha1.PackageRef{
				RefName: "test-pkg",
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: "1.0.0",
				},
			},
			ServiceAccountName: "use-local-cluster-sa",
		},
	}

	// Create an existing App in ReconcileFailed state with matching spec (no update needed)
	existingApp := &v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-pkg-install",
			Generation: 3,
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				{
					ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{
						Image: "test-pkg:1.0.0", // Same as package spec
					},
				},
			},
			Template: []v1alpha1.AppTemplate{
				{
					Ytt: &v1alpha1.AppTemplateYtt{},
				},
			},
			Deploy: []v1alpha1.AppDeploy{
				{
					Kapp: &v1alpha1.AppDeployKapp{},
				},
			},
			ServiceAccountName: "use-local-cluster-sa",
		},
		Status: v1alpha1.AppStatus{
			GenericStatus: v1alpha1.GenericStatus{
				ObservedGeneration: 3, // Generation == ObservedGeneration (stable failed state)
				Conditions: []v1alpha1.Condition{
					{
						Type:   v1alpha1.ReconcileFailed,
						Status: corev1.ConditionTrue,
					},
				},
				FriendlyDescription: "Reconcile failed",
				UsefulErrorMessage:  "Deploy failed for v1.0.0",
			},
		},
	}

	fakekctrl := fakekappctrl.NewSimpleClientset(model, existingApp)
	fakek8s := fake.NewSimpleClientset()

	// mock the kubernetes server version
	fakeDiscovery, _ := fakek8s.Discovery().(*fakediscovery.FakeDiscovery)
	fakeDiscovery.FakedServerVersion = &version.Info{
		GitVersion: "v0.20.0",
	}

	ip := NewPackageInstallCR(model, log, fakekctrl, fakePkgClient, fakek8s,
		FakeComponentInfo{KCVersion: semver.MustParse("0.42.31337")}, Opts{},
		metrics.NewMetrics())

	// Reconcile should NOT update the App since specs match
	_, err := ip.Reconcile()
	assert.Nil(t, err)

	// Verify the App was NOT updated (generation should remain the same)
	gvr := schema.GroupVersionResource{"kappctrl.k14s.io", "v1alpha1", "apps"}
	obj, err := fakekctrl.Tracker().Get(gvr, "", "test-pkg-install")
	assert.Nil(t, err)
	require.NotNil(t, obj)
	appAfterReconcile := obj.(*v1alpha1.App)

	// App should be unchanged
	assert.Equal(t, int64(3), appAfterReconcile.Generation, "App generation should not have changed")
	assert.Equal(t, "test-pkg:1.0.0", appAfterReconcile.Spec.Fetch[0].ImgpkgBundle.Image)

	// PackageInstall status should reflect the existing App state (ReconcileFailed)
	assert.Len(t, ip.model.Status.Conditions, 1)
	assert.Equal(t, v1alpha1.ReconcileFailed, ip.model.Status.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, ip.model.Status.Conditions[0].Status)
	assert.Equal(t, "Deploy failed for v1.0.0", ip.model.Status.UsefulErrorMessage)
}

func generatePackageWithConstraints(name, version, kcConstraint string, k8sConstraint string) datapkgingv1alpha1.Package {
	return datapkgingv1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: name + "." + version,
		},
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: name,
			Version: version,
			KappControllerVersionSelection: &datapkgingv1alpha1.VersionSelection{
				Constraints: kcConstraint,
			},
			KubernetesVersionSelection: &datapkgingv1alpha1.VersionSelection{
				Constraints: k8sConstraint,
			},
		},
	}
}
