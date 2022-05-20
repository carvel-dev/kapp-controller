// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository_test

import (
	"reflect"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	apppkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
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
		reftracker.NewSecretKey("s5", "default"): struct{}{},
	}

	appWithRefs := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-app",
			Namespace: "default",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{PathsFrom: []v1alpha1.AppFetchInlineSource{{SecretRef: &v1alpha1.AppFetchInlineSourceRef{"", "s"}}}}},
				v1alpha1.AppFetch{Image: &v1alpha1.AppFetchImage{SecretRef: &v1alpha1.AppFetchLocalRef{"s1"}}},
				v1alpha1.AppFetch{HTTP: &v1alpha1.AppFetchHTTP{SecretRef: &v1alpha1.AppFetchLocalRef{"s2"}}},
				v1alpha1.AppFetch{Git: &v1alpha1.AppFetchGit{SecretRef: &v1alpha1.AppFetchLocalRef{"s3"}}},
				v1alpha1.AppFetch{ImgpkgBundle: &v1alpha1.AppFetchImgpkgBundle{SecretRef: &v1alpha1.AppFetchLocalRef{"s5"}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, exec.NewPlainCmdRunner())

	app := apppkg.NewApp(appWithRefs, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log, nil)

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
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, exec.NewPlainCmdRunner())

	app := apppkg.NewApp(appEmpty, apppkg.Hooks{}, fetchFac, tmpFac, deployFac, log, nil)

	out := app.SecretRefs()
	if len(out) != 0 {
		t.Fatalf("\n Expected: %s\nGot: %s\n", "No SecretRefs to be returned", out)
	}
}
