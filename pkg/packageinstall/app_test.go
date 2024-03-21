// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/packageinstall"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// several tests below have no SyncPeriod set so they'll all use the same default.
var defaultSyncPeriod metav1.Duration = metav1.Duration{Duration: 10 * time.Minute}

func TestAppManuallyControlled(t *testing.T) {
	existingApp := &kcv1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/manually-controlled": "",
			},
		},
	}

	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			Values: []pkgingv1alpha1.PackageInstallValues{
				{SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values1"}},
				{SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values2"}},
			},
		},
	}

	pkgVersion := datapkgingv1alpha1.Package{
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expec-pkg",
			Version: "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{
					Template: []kcv1alpha1.AppTemplate{
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
					},
				},
			},
		},
	}

	app, err := packageinstall.NewApp(existingApp, ipkg, pkgVersion, packageinstall.Opts{})
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/manually-controlled": "",
			},
		},
	}

	if !reflect.DeepEqual(expectedApp, app) {
		bs, _ := yaml.Marshal(app)
		t.Fatalf("App does not match expected app: (actual)\n%s", bs)
	}
}

func TestAppCustomFetchSecretNames(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/fetch-0-secret-name": "secret0-name",
				"ext.packaging.carvel.dev/fetch-1-secret-name": "secret1-name",
				// no secret for fetch 2
				"ext.packaging.carvel.dev/fetch-3-secret-name": "secret3-name",
				"ext.packaging.carvel.dev/fetch-4-secret-name": "secret4-name",
				"ext.packaging.carvel.dev/fetch-5-secret-name": "secret5-name",
				"ext.packaging.carvel.dev/fetch-6-secret-name": "secret6-name",
			},
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			SyncPeriod: &metav1.Duration{100 * time.Second},
		},
	}

	pkgVersion := datapkgingv1alpha1.Package{
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expec-pkg",
			Version: "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{
					Fetch: []kcv1alpha1.AppFetch{
						{HelmChart: &kcv1alpha1.AppFetchHelmChart{}},       // 0
						{ImgpkgBundle: &kcv1alpha1.AppFetchImgpkgBundle{}}, // 1
						{Image: &kcv1alpha1.AppFetchImage{}},               // 2
						{Git: &kcv1alpha1.AppFetchGit{}},                   // 3
						{HTTP: &kcv1alpha1.AppFetchHTTP{ // 4
							SecretRef: &kcv1alpha1.AppFetchLocalRef{Name: "existing-secret-name"},
						}},
						{HelmChart: &kcv1alpha1.AppFetchHelmChart{ // 5
							Repository: &kcv1alpha1.AppFetchHelmChartRepo{},
						}},
						{Image: &kcv1alpha1.AppFetchImage{}}, // 6
					},
				},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion, packageinstall.Opts{})
	require.NoError(t, err)

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{100 * time.Second},
			Fetch: []kcv1alpha1.AppFetch{
				{HelmChart: &kcv1alpha1.AppFetchHelmChart{ // 0
					// no repository specified, so no secret set
				}},
				{ImgpkgBundle: &kcv1alpha1.AppFetchImgpkgBundle{ // 1
					SecretRef: &kcv1alpha1.AppFetchLocalRef{Name: "secret1-name"},
				}},
				{Image: &kcv1alpha1.AppFetchImage{ // 2
					// no annotation specified, no secret set
				}},
				{Git: &kcv1alpha1.AppFetchGit{ // 3
					SecretRef: &kcv1alpha1.AppFetchLocalRef{Name: "secret3-name"},
				}},
				{HTTP: &kcv1alpha1.AppFetchHTTP{ // 4
					SecretRef: &kcv1alpha1.AppFetchLocalRef{Name: "secret4-name"},
				}},
				{HelmChart: &kcv1alpha1.AppFetchHelmChart{ // 5
					Repository: &kcv1alpha1.AppFetchHelmChartRepo{
						SecretRef: &kcv1alpha1.AppFetchLocalRef{Name: "secret5-name"},
					},
				}},
				{Image: &kcv1alpha1.AppFetchImage{ // 6
					SecretRef: &kcv1alpha1.AppFetchLocalRef{Name: "secret6-name"},
				}},
			},
		},
	}

	// Not interesting in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	require.Equal(t, expectedApp, app, "App does not match expected app")
}

// TestAppPackageIntallDefaultSyncPeriod tests the creation of an App and expects syncPeriod is set to the default value.
func TestAppPackageIntallDefaultSyncPeriod(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
	}

	pkgVersion := datapkgingv1alpha1.Package{
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expec-pkg",
			Version: "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion, packageinstall.Opts{DefaultSyncPeriod: 10 * time.Minute})
	require.NoError(t, err)

	// Define the expected app object, with the sync period attribute set to the default value
	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &defaultSyncPeriod,
		},
	}

	// Not interested in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	require.Equal(t, expectedApp, app, "App does not match expected app")
}

// TestAppCustomPackageIntallSyncPeriod tests the creation of an App when using PackageInstall with a defined syncPeriod.
func TestAppCustomPackageIntallSyncPeriod(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			SyncPeriod: &metav1.Duration{Duration: 100 * time.Second},
		},
	}

	pkgVersion := datapkgingv1alpha1.Package{
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expec-pkg",
			Version: "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion, packageinstall.Opts{})
	require.NoError(t, err)

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &metav1.Duration{Duration: 100 * time.Second},
		},
	}

	// Not interested in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	require.Equal(t, expectedApp, app, "App does not match expected app")
}

func TestAppPackageDetailsAnnotations(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{},
	}

	pkgVersion := datapkgingv1alpha1.Package{
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expec-pkg",
			Version: "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion, packageinstall.Opts{})
	require.NoError(t, err)

	trueVal := true
	expectedObjectMeta := metav1.ObjectMeta{
		Name:      "app",
		Namespace: "default",
		Annotations: map[string]string{
			"packaging.carvel.dev/package-ref-name": "expec-pkg",
			"packaging.carvel.dev/package-version":  "1.5.0",
		},
		OwnerReferences: []metav1.OwnerReference{{
			APIVersion:         "packaging.carvel.dev/v1alpha1",
			Kind:               "PackageInstall",
			Name:               "app",
			UID:                "",
			Controller:         &trueVal,
			BlockOwnerDeletion: &trueVal,
		}},
	}

	require.Equal(t, expectedObjectMeta, app.ObjectMeta)
}

// TestAppPackageIntallDefaultNamespace tests the creation of an App when using PackageInstall with a defaultNamespace defined.
func TestAppPackageIntallDefaultNamespace(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			DefaultNamespace: "default-namespace",
			SyncPeriod:       &metav1.Duration{100 * time.Second},
		},
	}

	pkgVersion := datapkgingv1alpha1.Package{
		Spec: datapkgingv1alpha1.PackageSpec{
			RefName: "expec-pkg",
			Version: "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion, packageinstall.Opts{})
	require.NoError(t, err)

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			DefaultNamespace: "default-namespace",
			SyncPeriod:       &metav1.Duration{100 * time.Second},
		},
	}

	// Not interested in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	require.Equal(t, expectedApp, app, "App does not match expected app")
}

func TestAppPackageIntallValuesForTemplateSteps(t *testing.T) {
	pkgi := func(values []pkgingv1alpha1.PackageInstallValues) *pkgingv1alpha1.PackageInstall {
		return &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "app",
				Namespace: "default",
			},
			Spec: pkgingv1alpha1.PackageInstallSpec{
				Values: values,
			},
		}
	}

	pkg := func() datapkgingv1alpha1.Package {
		return datapkgingv1alpha1.Package{
			Spec: datapkgingv1alpha1.PackageSpec{
				RefName: "expec-pkg",
				Version: "1.5.0",
				Template: datapkgingv1alpha1.AppTemplateSpec{
					Spec: &kcv1alpha1.AppSpec{
						Template: []kcv1alpha1.AppTemplate{
							{Sops: &kcv1alpha1.AppTemplateSops{}},
							{Ytt: &kcv1alpha1.AppTemplateYtt{}},
							{Ytt: &kcv1alpha1.AppTemplateYtt{}},
							{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
							{Kbld: &kcv1alpha1.AppTemplateKbld{}},
							{Ytt: &kcv1alpha1.AppTemplateYtt{}},
							{Cue: &kcv1alpha1.AppTemplateCue{}},
						},
					},
				},
			},
		}
	}

	testCases := map[string]struct {
		values           []pkgingv1alpha1.PackageInstallValues
		patchPkg         func(*datapkgingv1alpha1.Package)
		patchPkgi        func(*pkgingv1alpha1.PackageInstall)
		expectedTemplate []kcv1alpha1.AppTemplate
		exepectedErrMsg  string
	}{
		"no values": {
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{Sops: &kcv1alpha1.AppTemplateSops{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				{Kbld: &kcv1alpha1.AppTemplateKbld{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{Cue: &kcv1alpha1.AppTemplateCue{}},
			},
		},
		"only add to first step": {
			values: []pkgingv1alpha1.PackageInstallValues{
				{SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-first"}},
			},
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{Sops: &kcv1alpha1.AppTemplateSops{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-first"}},
					},
				}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				{Kbld: &kcv1alpha1.AppTemplateKbld{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{Cue: &kcv1alpha1.AppTemplateCue{}},
			},
		},
		"invalid step index": {
			values: []pkgingv1alpha1.PackageInstallValues{
				{TemplateSteps: []int{100}},
			},
			exepectedErrMsg: "out of range",
		},
		"specific values, but step does not support values": {
			values: []pkgingv1alpha1.PackageInstallValues{
				{TemplateSteps: []int{0}, SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "foo-bar"}},
			},
			exepectedErrMsg: "does not support values",
		},
		"some values, but no steps which takes values": {
			values: []pkgingv1alpha1.PackageInstallValues{
				{SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "foo-bar"}},
			},
			patchPkg: func(pkg *datapkgingv1alpha1.Package) {
				// pkg does define any template steps that support values
				pkg.Spec.Template.Spec.Template = []kcv1alpha1.AppTemplate{
					{Sops: &kcv1alpha1.AppTemplateSops{}},
				}
			},
			exepectedErrMsg: "no template step of class 'takesValues' found",
		},
		"values for specific steps": {
			values: []pkgingv1alpha1.PackageInstallValues{
				{SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-first"}},
				{TemplateSteps: []int{6}, SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-specific-cue"}},
				{TemplateSteps: []int{3, 6}, SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-multiple-steps"}},
				{TemplateSteps: []int{5}, SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-specific-ytt"}},
			},
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{Sops: &kcv1alpha1.AppTemplateSops{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-first"}},
					},
				}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-multiple-steps"}},
					},
				}},
				{Kbld: &kcv1alpha1.AppTemplateKbld{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-specific-ytt"}},
					},
				}},
				{Cue: &kcv1alpha1.AppTemplateCue{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-specific-cue"}},
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-multiple-steps"}},
					},
				}},
			},
		},
		"values for ytt, as inline paths": {
			patchPkgi: func(pkgi *pkgingv1alpha1.PackageInstall) {
				pkgi.ObjectMeta.SetAnnotations(map[string]string{
					"ext.packaging.carvel.dev/ytt-data-values-overlays":   "",
					"ext.packaging.carvel.dev/ytt-5-data-values-overlays": "",
				})
			},
			values: []pkgingv1alpha1.PackageInstallValues{
				{SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-first"}},
				{TemplateSteps: []int{2}, SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-specific-ytt"}},
				{TemplateSteps: []int{5}, SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "values-for-another-ytt"}},
			},
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{Sops: &kcv1alpha1.AppTemplateSops{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							{SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{Name: "values-for-first"}},
						},
					},
				}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "values-for-specific-ytt"}},
					},
				}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				{Kbld: &kcv1alpha1.AppTemplateKbld{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							{SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{Name: "values-for-another-ytt"}},
						},
					},
				}},
				{Cue: &kcv1alpha1.AppTemplateCue{}},
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			pkg := pkg()
			if tc.patchPkg != nil {
				tc.patchPkg(&pkg)
			}

			pkgi := pkgi(tc.values)
			if tc.patchPkgi != nil {
				tc.patchPkgi(pkgi)
			}

			app, err := packageinstall.NewApp(&kcv1alpha1.App{}, pkgi, pkg, packageinstall.Opts{})

			if errMsg := tc.exepectedErrMsg; errMsg != "" {
				require.ErrorContains(t, err, errMsg)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedTemplate, app.Spec.Template, "App template does not match expected template")
		})
	}
}

func TestAppPackageIntallAnnotationsForTemplateSteps(t *testing.T) {
	pkgi := func(annotations map[string]string) *pkgingv1alpha1.PackageInstall {
		pkgi := &pkgingv1alpha1.PackageInstall{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "app",
				Namespace: "default",
			},
		}

		pkgi.ObjectMeta.SetAnnotations(annotations)

		return pkgi
	}

	pkg := func(templateSteps []kcv1alpha1.AppTemplate) datapkgingv1alpha1.Package {
		return datapkgingv1alpha1.Package{
			Spec: datapkgingv1alpha1.PackageSpec{
				RefName: "expec-pkg",
				Version: "1.5.0",
				Template: datapkgingv1alpha1.AppTemplateSpec{
					Spec: &kcv1alpha1.AppSpec{
						Template: templateSteps,
					},
				},
			},
		}
	}

	someHelmTemplateSteps := func() []kcv1alpha1.AppTemplate {
		return []kcv1alpha1.AppTemplate{
			{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
			{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
			{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
		}
	}
	someYttTemplateSteps := func() []kcv1alpha1.AppTemplate {
		return []kcv1alpha1.AppTemplate{
			{Ytt: &kcv1alpha1.AppTemplateYtt{}},
			{Ytt: &kcv1alpha1.AppTemplateYtt{}},
			{Ytt: &kcv1alpha1.AppTemplateYtt{}},
		}
	}

	testCases := map[string]struct {
		pkgiAnnotations  map[string]string
		pkgTemplateSteps []kcv1alpha1.AppTemplate
		expectedTemplate []kcv1alpha1.AppTemplate
	}{
		"no annotations": {
			pkgTemplateSteps: someHelmTemplateSteps(),
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
			},
		},
		"helm annotations": {
			pkgTemplateSteps: someHelmTemplateSteps(),
			pkgiAnnotations: map[string]string{
				"ext.packaging.carvel.dev/helm-template-name":          "some-default-helm-name",
				"ext.packaging.carvel.dev/helm-1-template-name":        "some-specific-helm-name",
				"ext.packaging.carvel.dev/helm-2-template-namespace":   "some-specific-helm-namespace",
				"ext.packaging.carvel.dev/helm-100-template-namespace": "no-such-step",

				"ext.packaging.carvel.dev/helm-template-values-from-secret-name":       "some-helm-secret",
				"ext.packaging.carvel.dev/helm-template-values-from-secret-name.blipp": "some-helm-secret.blipp",
				"ext.packaging.carvel.dev/helm-1-template-values-from-secret-name.foo": "some-helm-secret.foo",
				"ext.packaging.carvel.dev/helm-1-template-values-from-secret-name.bar": "some-helm-secret.bar",
			},
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					Name: "some-default-helm-name",
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "some-helm-secret"}},
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "some-helm-secret.blipp"}},
					},
				}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					Name: "some-specific-helm-name",
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "some-helm-secret.bar"}},
						{SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "some-helm-secret.foo"}},
					},
				}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					Namespace: "some-specific-helm-namespace",
				}},
			},
		},
		"ytt annotations": {
			pkgTemplateSteps: someYttTemplateSteps(),
			pkgiAnnotations: map[string]string{
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name":          "some-ytt-secret",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.1":        "some-other-ytt-secret",
				"ext.packaging.carvel.dev/ytt-0-paths-from-secret-name.foobar": "some-third-ytt-secret",
				"ext.packaging.carvel.dev/ytt-2-paths-from-secret-name":        "some-specific-ytt-secret",
				"ext.packaging.carvel.dev/ytt-100-paths-from-secret-name":      "no such step",
			},
			expectedTemplate: []kcv1alpha1.AppTemplate{
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							{SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{Name: "some-ytt-secret"}},
							{SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{Name: "some-other-ytt-secret"}},
							{SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{Name: "some-third-ytt-secret"}},
						},
					},
				}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							{SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{Name: "some-specific-ytt-secret"}},
						},
					},
				}},
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			app, err := packageinstall.NewApp(
				&kcv1alpha1.App{},
				pkgi(tc.pkgiAnnotations),
				pkg(tc.pkgTemplateSteps),
				packageinstall.Opts{},
			)
			require.NoError(t, err)
			require.Equal(t, tc.expectedTemplate, app.Spec.Template, "App template does not match expected template")
		})
	}
}
