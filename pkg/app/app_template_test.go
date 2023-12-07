// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"os"
	"testing"

	"github.com/k14s/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_BuildAdditionalDownwardAPIValues_MemoizedCallCount(t *testing.T) {
	log := logf.Log.WithName("kc")

	appEmpty := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{},
			Template: []v1alpha1.AppTemplate{
				{Ytt: &v1alpha1.AppTemplateYtt{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{KubernetesVersion: &v1alpha1.Version{}},
					{KappControllerVersion: &v1alpha1.Version{}},
					{KubernetesAPIs: &v1alpha1.KubernetesAPIs{}}},
				}}}}},
				{Ytt: &v1alpha1.AppTemplateYtt{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{KubernetesVersion: &v1alpha1.Version{}},
					{KappControllerVersion: &v1alpha1.Version{}},
					{KubernetesAPIs: &v1alpha1.KubernetesAPIs{}}},
				}}}}},
				{HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{
					KubernetesVersion: &v1alpha1.Version{},
					KubernetesAPIs:    &v1alpha1.KubernetesAPIs{},
				}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	k8sVersionCallCount, kcVersionCallCount, k8sAPIsCallCount := 0, 0, 0
	fakeInfo := FakeComponentInfo{K8sVersion: semver.MustParse("1.1.0"), KCVersion: semver.MustParse("2.0.0"), K8sAPIs: []string{"test", "test2"},
		K8sVersionCount: &k8sVersionCallCount,
		K8sAPIsCount:    &k8sAPIsCallCount,
		KCVersionCount:  &kcVersionCallCount,
	}
	app := NewApp(appEmpty, Hooks{}, fetchFac, tmpFac, deployFac, log, Opts{}, metrics.NewCountMetrics(), metrics.NewReconcileTimeMetrics(), fakeInfo)

	dir, err := os.MkdirTemp("", "temp")
	assert.NoError(t, err)

	app.template(dir)

	assert.Equal(t, 1, k8sVersionCallCount)
	assert.Equal(t, 2, kcVersionCallCount)
	assert.Equal(t, 1, k8sAPIsCallCount)
}
