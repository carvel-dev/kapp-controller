// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app_test

import (
	"reflect"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	apppkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_SecretRefs_RetrievesAllSecretRefs(t *testing.T) {
	log := logf.Log.WithName("kc")
	expected := map[reftracker.RefKey]struct{}{
		reftracker.NewSecretKey("s", "default"):  struct{}{},
		reftracker.NewSecretKey("s1", "default"): struct{}{},
		reftracker.NewSecretKey("s2", "default"): struct{}{},
		reftracker.NewSecretKey("s3", "default"): struct{}{},
		reftracker.NewSecretKey("s4", "default"): struct{}{},
		reftracker.NewSecretKey("s5", "default"): struct{}{},
		reftracker.NewSecretKey("s6", "default"): struct{}{},
		reftracker.NewSecretKey("s7", "default"): struct{}{},
	}

	appWithRefs := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{SecretRef: &v1alpha1.AppFetchInlineSourceRef{"", v1.LocalObjectReference{Name: "s"}}}}}},
				v1alpha1.AppFetch{Image: &v1alpha1.AppFetchImage{SecretRef: &v1alpha1.AppFetchLocalRef{v1.LocalObjectReference{Name: "s1"}}}},
				v1alpha1.AppFetch{HTTP: &v1alpha1.AppFetchHTTP{SecretRef: &v1alpha1.AppFetchLocalRef{v1.LocalObjectReference{Name: "s2"}}}},
				v1alpha1.AppFetch{Git: &v1alpha1.AppFetchGit{SecretRef: &v1alpha1.AppFetchLocalRef{v1.LocalObjectReference{Name: "s3"}}}},
				v1alpha1.AppFetch{HelmChart: &v1alpha1.AppFetchHelmChart{Repository: &v1alpha1.AppFetchHelmChartRepo{SecretRef: &v1alpha1.AppFetchLocalRef{v1.LocalObjectReference{Name: "s4"}}}}},
				v1alpha1.AppFetch{ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{SecretRef: &v1alpha1.AppFetchLocalRef{v1.LocalObjectReference{Name: "s5"}}}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{SecretRef: &v1alpha1.AppFetchInlineSourceRef{"", v1.LocalObjectReference{Name: "s6"}}}}}}},
				v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{ValuesFrom: []v1alpha1.AppTemplateHelmTemplateValuesSource{{SecretRef: &v1alpha1.AppTemplateHelmTemplateValuesSourceRef{v1.LocalObjectReference{Name: "s7"}}}}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, nil)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	app := apppkg.NewApp(appWithRefs, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log)

	out := app.SecretRefs()
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\n Expected: %s\nGot: %s\n", expected, out)
	}
}

func Test_SecretRefs_RetrievesNoSecretRefs_WhenNonePresent(t *testing.T) {
	log := logf.Log.WithName("kc")

	appEmpty := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch:    []v1alpha1.AppFetch{},
			Template: []v1alpha1.AppTemplate{},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, nil)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	app := apppkg.NewApp(appEmpty, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log)

	out := app.SecretRefs()
	if len(out) != 0 {
		t.Fatalf("\n Expected: %s\nGot: %s\n", "No SecretRefs to be returned", out)
	}
}

func Test_ConfigMapRefs_RetrievesAllConfigMapRefs(t *testing.T) {
	log := logf.Log.WithName("kc")

	expected := map[reftracker.RefKey]struct{}{
		reftracker.NewConfigMapKey("c", "default"):  struct{}{},
		reftracker.NewConfigMapKey("c1", "default"): struct{}{},
		reftracker.NewConfigMapKey("c2", "default"): struct{}{},
	}

	appWithRefs := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{ConfigMapRef: &v1alpha1.AppFetchInlineSourceRef{"", v1.LocalObjectReference{Name: "c"}}}}}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{ConfigMapRef: &v1alpha1.AppFetchInlineSourceRef{"", v1.LocalObjectReference{Name: "c1"}}}}}}},
				v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{ValuesFrom: []v1alpha1.AppTemplateHelmTemplateValuesSource{{ConfigMapRef: &v1alpha1.AppTemplateHelmTemplateValuesSourceRef{v1.LocalObjectReference{Name: "c2"}}}}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, nil)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	app := apppkg.NewApp(appWithRefs, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log)

	out := app.ConfigMapRefs()
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\n Expected: %s\nGot: %s\n", expected, out)
	}
}

func Test_ConfigMapRefs_RetrievesNoConfigMapRefs_WhenNonePresent(t *testing.T) {
	log := logf.Log.WithName("kc")

	appEmpty := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch:    []v1alpha1.AppFetch{},
			Template: []v1alpha1.AppTemplate{},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, nil)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	app := apppkg.NewApp(appEmpty, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log)

	out := app.ConfigMapRefs()
	if len(out) != 0 {
		t.Fatalf("\n Expected: %s\nGot: %s\n", "No ConfigMapRefs to be returned", out)
	}
}
