// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app_test

import (
	"reflect"
	"testing"

	"github.com/k14s/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	apppkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
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
		reftracker.NewSecretKey("s8", "default"): struct{}{},
	}

	appWithRefs := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{SecretRef: &v1alpha1.AppFetchInlineSourceRef{"", "s"}}}}},
				v1alpha1.AppFetch{Image: &v1alpha1.AppFetchImage{SecretRef: &v1alpha1.AppFetchLocalRef{Name: "s1"}}},
				v1alpha1.AppFetch{HTTP: &v1alpha1.AppFetchHTTP{SecretRef: &v1alpha1.AppFetchLocalRef{Name: "s2"}}},
				v1alpha1.AppFetch{Git: &v1alpha1.AppFetchGit{SecretRef: &v1alpha1.AppFetchLocalRef{Name: "s3"}}},
				v1alpha1.AppFetch{HelmChart: &v1alpha1.AppFetchHelmChart{Repository: &v1alpha1.AppFetchHelmChartRepo{SecretRef: &v1alpha1.AppFetchLocalRef{Name: "s4"}}}},
				v1alpha1.AppFetch{ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{SecretRef: &v1alpha1.AppFetchLocalRef{Name: "s5"}}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{SecretRef: &v1alpha1.AppFetchInlineSourceRef{"", "s6"}}}}}},
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{SecretRef: &v1alpha1.AppTemplateValuesSourceRef{Name: "s8"}}}}},
				v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{SecretRef: &v1alpha1.AppTemplateValuesSourceRef{Name: "s7"}}}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	app := apppkg.NewApp(appWithRefs, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log, apppkg.Opts{}, nil, FakeComponentInfo{})

	out := app.SecretRefs()
	assert.Truef(t, reflect.DeepEqual(out, expected), "Expected: %s\nGot: %s\n", expected, out)
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
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	app := apppkg.NewApp(appEmpty, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log, apppkg.Opts{}, nil, FakeComponentInfo{})

	out := app.SecretRefs()
	assert.Equal(t, 0, len(out), "No SecretRefs to be returned")
}

func Test_ConfigMapRefs_RetrievesAllConfigMapRefs(t *testing.T) {
	log := logf.Log.WithName("kc")

	expected := map[reftracker.RefKey]struct{}{
		reftracker.NewConfigMapKey("c", "default"):  struct{}{},
		reftracker.NewConfigMapKey("c1", "default"): struct{}{},
		reftracker.NewConfigMapKey("c2", "default"): struct{}{},
		reftracker.NewConfigMapKey("c3", "default"): struct{}{},
	}

	appWithRefs := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{ConfigMapRef: &v1alpha1.AppFetchInlineSourceRef{"", "c"}}}}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{ConfigMapRef: &v1alpha1.AppFetchInlineSourceRef{"", "c1"}}}}}},
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{ConfigMapRef: &v1alpha1.AppTemplateValuesSourceRef{Name: "c3"}}}}},
				v1alpha1.AppTemplate{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{ConfigMapRef: &v1alpha1.AppTemplateValuesSourceRef{"c2"}}}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	app := apppkg.NewApp(appWithRefs, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log, apppkg.Opts{}, nil, FakeComponentInfo{})

	out := app.ConfigMapRefs()
	assert.Truef(t, reflect.DeepEqual(out, expected), "Expected: %s\nGot: %s\n", expected, out)
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
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	app := apppkg.NewApp(appEmpty, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log, apppkg.Opts{}, nil, FakeComponentInfo{})

	out := app.ConfigMapRefs()
	assert.Lenf(t, out, 0, "Expected: %s\nGot: %s\n", "No ConfigMapRefs to be returned", out)
}

type FakeComponentInfo struct {
	KCVersion  semver.Version
	K8sVersion semver.Version
	K8sAPIs    []string
}

func (f FakeComponentInfo) KubernetesAPIs() ([]string, error) {
	return f.K8sAPIs, nil
}

func (f FakeComponentInfo) KappControllerVersion() (semver.Version, error) {
	return f.KCVersion, nil
}

func (f FakeComponentInfo) KubernetesVersion(_ string, _ *v1alpha1.AppCluster, _ *metav1.ObjectMeta) (semver.Version, error) {
	return f.K8sVersion, nil
}
