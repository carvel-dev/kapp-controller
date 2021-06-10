// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installedpkg_test

import (
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/installedpkg"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAppExtYttPathsFromSecretNameAnn(t *testing.T) {
	ipkg := &pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name":      "no-suffix",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.4":    "suffix-4",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.2":    "suffix-2",
				"ext.packaging.carvel.dev/ytt-paths-from-secret-name.text": "suffix-text",
			},
		},
	}

	pkgVersion := datapkgingv1alpha1.PackageVersion{
		Spec: datapkgingv1alpha1.PackageVersionSpec{
			PackageName: "expec-pkg",
			Version:     "1.5.0",
			Template: datapkgingv1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{
					Template: []kcv1alpha1.AppTemplate{
						{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
						{Ytt: &kcv1alpha1.AppTemplateYtt{}},
					},
				},
			},
		},
	}

	app, err := installedpkg.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			Template: []kcv1alpha1.AppTemplate{
				// Helm template step is untouched
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}},
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{Name: "no-suffix"},
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{Name: "suffix-2"},
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{Name: "suffix-4"},
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{Name: "suffix-text"},
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

func TestAppExtYttDataValuesOverlaysAnn(t *testing.T) {
	ipkg := &pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"ext.packaging.carvel.dev/ytt-data-values-overlays": "",
			},
		},
		Spec: pkgingv1alpha1.InstalledPackageSpec{
			Values: []pkgingv1alpha1.InstalledPackageValues{
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values1"}},
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values2"}},
			},
		},
	}

	pkgVersion := datapkgingv1alpha1.PackageVersion{
		Spec: datapkgingv1alpha1.PackageVersionSpec{
			PackageName: "expec-pkg",
			Version:     "1.5.0",
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

	app, err := installedpkg.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			Template: []kcv1alpha1.AppTemplate{
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					Inline: &kcv1alpha1.AppFetchInline{
						PathsFrom: []kcv1alpha1.AppFetchInlineSource{
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{Name: "values1"},
								},
							},
							kcv1alpha1.AppFetchInlineSource{
								SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
									LocalObjectReference: corev1.LocalObjectReference{Name: "values2"},
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
	ipkg := &pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.InstalledPackageSpec{
			Values: []pkgingv1alpha1.InstalledPackageValues{
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values1"}},
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values2"}},
			},
		},
	}

	pkgVersion := datapkgingv1alpha1.PackageVersion{
		Spec: datapkgingv1alpha1.PackageVersionSpec{
			PackageName: "expec-pkg",
			Version:     "1.5.0",
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

	app, err := installedpkg.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			Template: []kcv1alpha1.AppTemplate{
				{Ytt: &kcv1alpha1.AppTemplateYtt{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								LocalObjectReference: corev1.LocalObjectReference{Name: "values1"},
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								LocalObjectReference: corev1.LocalObjectReference{Name: "values2"},
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
	ipkg := &pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.InstalledPackageSpec{
			Values: []pkgingv1alpha1.InstalledPackageValues{
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values1"}},
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values2"}},
			},
		},
	}

	pkgVersion := datapkgingv1alpha1.PackageVersion{
		Spec: datapkgingv1alpha1.PackageVersionSpec{
			PackageName: "expec-pkg",
			Version:     "1.5.0",
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

	app, err := installedpkg.NewApp(&kcv1alpha1.App{}, ipkg, pkgVersion)
	if err != nil {
		t.Fatalf("Expected no err, but was: %s", err)
	}

	expectedApp := &kcv1alpha1.App{
		Spec: kcv1alpha1.AppSpec{
			Template: []kcv1alpha1.AppTemplate{
				{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
					ValuesFrom: []kcv1alpha1.AppTemplateValuesSource{
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								LocalObjectReference: corev1.LocalObjectReference{Name: "values1"},
							},
						},
						kcv1alpha1.AppTemplateValuesSource{
							SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
								LocalObjectReference: corev1.LocalObjectReference{Name: "values2"},
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

	ipkg := &pkgingv1alpha1.InstalledPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
		Spec: pkgingv1alpha1.InstalledPackageSpec{
			Values: []pkgingv1alpha1.InstalledPackageValues{
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values1"}},
				{SecretRef: &pkgingv1alpha1.InstalledPackageValuesSecretRef{Name: "values2"}},
			},
		},
	}

	pkgVersion := datapkgingv1alpha1.PackageVersion{
		Spec: datapkgingv1alpha1.PackageVersionSpec{
			PackageName: "expec-pkg",
			Version:     "1.5.0",
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

	app, err := installedpkg.NewApp(existingApp, ipkg, pkgVersion)
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
