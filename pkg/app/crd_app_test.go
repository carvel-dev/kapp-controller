// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"reflect"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

func Test_GetSecretRefs_RetrievesAllSecretRefs(t *testing.T) {
	log := logf.Log.WithName("kc")
	expected := map[string]struct{}{
		"s":  struct{}{},
		"s1": struct{}{},
		"s2": struct{}{},
		"s3": struct{}{},
		"s4": struct{}{},
		"s5": struct{}{},
		"s6": struct{}{},
		"s7": struct{}{},
	}
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
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
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	crdApp := NewCRDApp(&app, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	out := crdApp.GetSecretRefs()
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\n Expected: %s\nGot: %s\n", expected, out)
	}
}

func Test_GetSecretRefs_RetrievesNoSecretRefs_WhenNonePresent(t *testing.T) {
	log := logf.Log.WithName("kc")
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch:    []v1alpha1.AppFetch{},
			Template: []v1alpha1.AppTemplate{},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	crdApp := NewCRDApp(&app, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	out := crdApp.GetSecretRefs()
	if len(out) != 0 {
		t.Fatalf("\n Expected: %s\nGot: %s\n", "No SecretRefs to be returned", out)
	}
}

func Test_GetConfigMapRefs_RetrievesAllConfigMapRefs(t *testing.T) {
	log := logf.Log.WithName("kc")
	expected := map[string]struct{}{
		"s":  struct{}{},
		"s1": struct{}{},
		"s2": struct{}{},
	}
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{ConfigMapRef: &v1alpha1.AppFetchInlineSourceRef{"", v1.LocalObjectReference{Name: "s"}}}}}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{ConfigMapRef: &v1alpha1.AppFetchInlineSourceRef{"", v1.LocalObjectReference{Name: "s1"}}}}}}},
				v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{ValuesFrom: []v1alpha1.AppTemplateHelmTemplateValuesSource{{ConfigMapRef: &v1alpha1.AppTemplateHelmTemplateValuesSourceRef{v1.LocalObjectReference{Name: "s2"}}}}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	crdApp := NewCRDApp(&app, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	out := crdApp.GetConfigMapRefs()
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\n Expected: %s\nGot: %s\n", expected, out)
	}
}

func Test_GetConfigMapRefs_RetrievesNoConfigMapRefs_WhenNonePresent(t *testing.T) {
	log := logf.Log.WithName("kc")
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch:    []v1alpha1.AppFetch{},
			Template: []v1alpha1.AppTemplate{},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	crdApp := NewCRDApp(&app, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	out := crdApp.GetConfigMapRefs()
	if len(out) != 0 {
		t.Fatalf("\n Expected: %s\nGot: %s\n", "No SecretRefs to be returned", out)
	}
}
