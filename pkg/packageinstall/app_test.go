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
var defaultSyncPeriod metav1.Duration = metav1.Duration{10 * time.Minute}

func TestAppExtPathsFromSecretNameAnn(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name":                 "ytt-no-suffix",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.4":               "ytt-suffix-4",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.2":               "ytt-suffix-2",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.text":            "ytt-suffix-text",
				"ext.packaging.carvel.dev/helm-template-values-from-secret-name":      "helm-no-suffix",
				"ext.packaging.carvel.dev/helm-template-values-from-secret-name.4":    "helm-suffix-4",
				"ext.packaging.carvel.dev/helm-template-values-from-secret-name.2":    "helm-suffix-2",
				"ext.packaging.carvel.dev/helm-template-values-from-secret-name.text": "helm-suffix-text",
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
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
					},
				},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &defaultSyncPeriod,
			Template: []kcv1alpha1.AppTemplate{
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "helm-no-suffix",
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "helm-suffix-2",
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "helm-suffix-4",
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "helm-suffix-text",
							},
						},
					},
				}},
				// Second Helm template step is untouched
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},

				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									Name: "ytt-no-suffix",
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									Name: "ytt-suffix-2",
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									Name: "ytt-suffix-4",
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									Name: "ytt-suffix-text",
								},
							},
						},
					},
				}},
				// Second ytt templating step is untouched
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
			},
		},
	}

	// Not interesting in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	if !reflect.DeepEqual(expectedApp, app) {
		bs, _ := yaml.Marshal(app)
		t.Fatalf("App does not match expected app: (actual)\n%s", bs)
	}
}
func TestAppHelmOverlaysFromAnn(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/helm-template-name":      "helm-new-name",
				"ext.packaging.carvel.dev/helm-template-namespace": "helm-new-namespace",
			},
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
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
							Name:      "helm-default-name",
							Namespace: "helm-default-namespace",
						}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
					},
				},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &defaultSyncPeriod,
			Template: []kcv1alpha1.AppTemplate{
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					Name:      "helm-new-name",
					Namespace: "helm-new-namespace",
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "values1",
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "values2",
							},
						},
					},
				}},
				// Ytt template step is untouched
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
			},
		},
	}

	// Not interesting in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	if !reflect.DeepEqual(expectedApp, app) {
		bs, _ := yaml.Marshal(app)
		t.Fatalf("App does not match expected app: (actual)\n%s", bs)
	}
}

func TestAppExtYttDataValuesOverlaysAnn(t *testing.T) {
	ipkg := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/ytt-data-values-overlays": "",
			},
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
					},
				},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &defaultSyncPeriod,
			Template: []kcv1alpha1.AppTemplate{
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									Name: "values1",
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									Name: "values2",
								},
							},
						},
					},
				}},
				// Second ytt templating step is untouched
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
			},
		},
	}

	// Not interesting in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	if !reflect.DeepEqual(expectedApp, app) {
		bs, _ := yaml.Marshal(app)
		t.Fatalf("App does not match expected app: (actual)\n%s", bs)
	}
}

func TestAppYttValues(t *testing.T) {
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

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &defaultSyncPeriod,
			Template: []kcv1alpha1.AppTemplate{
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "values1",
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "values2",
							},
						},
					},
				}},
				// Second ytt templating step is untouched
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
			},
		},
	}

	// Not interesting in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	if !reflect.DeepEqual(expectedApp, app) {
		bs, _ := yaml.Marshal(app)
		t.Fatalf("App does not match expected app: (actual)\n%s", bs)
	}
}

func TestAppHelmTemplateValues(t *testing.T) {
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
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
					},
				},
			},
		},
	}

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			SyncPeriod: &defaultSyncPeriod,
			Template: []kcv1alpha1.AppTemplate{
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "values1",
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								Name: "values2",
							},
						},
					},
				}},
				// Second helm templating step is untouched
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				// Second ytt templating step is untouched
				{Ytt: &kcv1alpha1.AppTemplateYtt{}},
			},
		},
	}

	// Not interesting in metadata in this test
	app.ObjectMeta = metav1.ObjectMeta{}

	if !reflect.DeepEqual(expectedApp, app) {
		bs, _ := yaml.Marshal(app)
		t.Fatalf("App does not match expected app: (actual)\n%s", bs)
	}
}

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

	app, err := packageinstall.NewApp(existingApp, ipkg, pkgVersion)
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

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
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

	app, err := packageinstall.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
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
