// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"
	"fmt"
	"testing"
	"time"

	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	fakekappctrl "carvel.dev/kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stesting "k8s.io/client-go/testing"
)

func TestNewAppPreservesPackageTemplate(t *testing.T) {
	for _, template := range []string{"ytt", "helm"} {
		t.Run(template, func(t *testing.T) {
			install, pkg := syntheticTemplateIsolationFixture(template)
			original := pkg.DeepCopy()
			var first *kcv1alpha1.App
			for i := 0; i < 6; i++ {
				app, err := NewApp(&kcv1alpha1.App{}, install, pkg, Opts{DefaultSyncPeriod: time.Minute})
				require.NoError(t, err)
				var values []kcv1alpha1.AppTemplateValuesSource
				if template == "ytt" {
					values = app.Spec.Template[0].Ytt.ValuesFrom
				} else {
					values = app.Spec.Template[0].HelmTemplate.ValuesFrom
				}
				require.Equal(t, []kcv1alpha1.AppTemplateValuesSource{{
					SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "synthetic-values"},
				}}, values)
				require.Equal(t, original, &pkg, "App construction must not mutate the package")
				if first == nil {
					first = app.DeepCopy()
				} else {
					require.Equal(t, first, app, "repeated construction must produce the same App")
				}
			}
		})
	}
}

func TestReconcileAppValuesRemainStableAcrossConflicts(t *testing.T) {
	for _, conflicts := range []int{0, 3} {
		t.Run(fmt.Sprintf("conflicts-%d", conflicts), func(t *testing.T) {
			install, pkg := syntheticTemplateIsolationFixture("ytt")
			original := pkg.DeepCopy()
			existing := &kcv1alpha1.App{ObjectMeta: metav1.ObjectMeta{
				Name: install.Name, Namespace: install.Namespace,
			}}
			client := fakekappctrl.NewSimpleClientset(existing)
			attempts := 0
			client.PrependReactor("update", "apps", func(_ k8stesting.Action) (bool, runtime.Object, error) {
				attempts++
				if attempts <= conflicts {
					return true, nil, errors.NewConflict(schema.GroupResource{
						Group: "kappctrl.k14s.io", Resource: "apps",
					}, install.Name, fmt.Errorf("synthetic concurrent update"))
				}
				return false, nil, nil
			})
			pi := &PackageInstallCR{
				model: install, kcclient: client, opts: Opts{DefaultSyncPeriod: time.Minute},
			}
			_, err := pi.reconcileAppWithPackage(existing, pkg, func(*kcv1alpha1.App) {})
			require.NoError(t, err)
			require.Equal(t, conflicts+1, attempts)
			persisted, err := client.KappctrlV1alpha1().Apps(install.Namespace).Get(
				context.Background(), install.Name, metav1.GetOptions{})
			require.NoError(t, err)
			require.Equal(t, []kcv1alpha1.AppTemplateValuesSource{{
				SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{Name: "synthetic-values"},
			}}, persisted.Spec.Template[0].Ytt.ValuesFrom)
			require.Equal(t, original, &pkg)

			// A subsequent reconciliation must converge without another App update.
			_, err = pi.reconcileAppWithPackage(persisted, pkg, func(*kcv1alpha1.App) {})
			require.NoError(t, err)
			require.Equal(t, conflicts+1, attempts)
		})
	}
}

func syntheticTemplateIsolationFixture(template string) (*pkgingv1alpha1.PackageInstall, datapkgingv1alpha1.Package) {
	install := &pkgingv1alpha1.PackageInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name: "synthetic-install", Namespace: "synthetic-namespace", UID: "synthetic-install-uid",
		},
		Spec: pkgingv1alpha1.PackageInstallSpec{
			Values: []pkgingv1alpha1.PackageInstallValues{{
				SecretRef: &pkgingv1alpha1.PackageInstallValuesSecretRef{Name: "synthetic-values"},
			}},
		},
	}
	step := kcv1alpha1.AppTemplate{Ytt: &kcv1alpha1.AppTemplateYtt{}}
	if template == "helm" {
		step = kcv1alpha1.AppTemplate{HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{}}
	}
	pkg := datapkgingv1alpha1.Package{Spec: datapkgingv1alpha1.PackageSpec{
		RefName: "synthetic-package.example.invalid", Version: "1.0.0",
		Template: datapkgingv1alpha1.AppTemplateSpec{Spec: &kcv1alpha1.AppSpec{
			Template: []kcv1alpha1.AppTemplate{step},
		}},
	}}
	return install, pkg
}
