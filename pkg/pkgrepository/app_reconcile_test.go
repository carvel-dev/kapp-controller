// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_NoInspectReconcile_IfNoDeployAttempted(t *testing.T) {
	log := logf.Log.WithName("kc")

	// The url under fetch is invalid, which will cause this
	// app to fail before deploy.
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{HTTP: &v1alpha1.AppFetchHTTP{URL: "i-dont-exist"}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)
	pkgr := v1alpha12.PackageRepository{}

	crdApp := NewCRDApp(&app, &pkgr, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	// Expected app status has no inspect on status
	// since the app deployment was not attempted
	expectedStatus := v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.Condition{{
				Type:    v1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Fetching resources: Error (see .status.usefulErrorMessage for details)",
			}},
			ObservedGeneration:  0,
			FriendlyDescription: "Reconcile failed: Fetching resources: Error (see .status.usefulErrorMessage for details)",
			UsefulErrorMessage:  "vendir: Error: Syncing directory '0':\n  Syncing directory '.' with HTTP contents:\n    Downloading URL:\n      Initiating URL download:\n        Get \"i-dont-exist\": unsupported protocol scheme \"\"\n",
		},
		Fetch: &v1alpha1.AppStatusFetch{
			Error:    "Fetching resources: Error (see .status.usefulErrorMessage for details)",
			Stderr:   "vendir: Error: Syncing directory '0':\n  Syncing directory '.' with HTTP contents:\n    Downloading URL:\n      Initiating URL download:\n        Get \"i-dont-exist\": unsupported protocol scheme \"\"\n",
			ExitCode: 1,
		},
		ConsecutiveReconcileFailures: 1,
	}

	assertStatusEqual(t, crdApp, expectedStatus)
}

func Test_TemplateError_DisplayedInStatus_UsefulErrorMessageProperty(t *testing.T) {
	log := logf.Log.WithName("kc")

	fetchInline := map[string]string{
		"packages/file.yml": `foo: #@ data.values.nothere`,
	}
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				{Inline: &v1alpha1.AppFetchInline{Paths: fetchInline}},
			},
			// Note: PKGR Template phase is hardcoded in app_template.go
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)
	pkgr := v1alpha12.PackageRepository{}

	crdApp := NewCRDApp(&app, &pkgr, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	// Expected app status has no inspect on status
	// since the app deployment was not attempted
	expectedStatus := v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.Condition{{
				Type:    v1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Templating dir: Error (see .status.usefulErrorMessage for details)",
			}},
			ObservedGeneration:  0,
			FriendlyDescription: "Reconcile failed: Templating dir: Error (see .status.usefulErrorMessage for details)",
			UsefulErrorMessage:  "ytt: Error: \n- undefined: data\n    file.yml:1 | foo: #@ data.values.nothere\n",
		},
		Fetch: &v1alpha1.AppStatusFetch{
			ExitCode: 0,
		},
		Template: &v1alpha1.AppStatusTemplate{
			Error:    "Templating dir: Error (see .status.usefulErrorMessage for details)",
			ExitCode: 1,
		},
	}

	assertStatusEqual(t, crdApp, expectedStatus)
}

func TestInvalidPackageRepositoryFormat(t *testing.T) {
	log := logf.Log.WithName("kc")

	fetchInline := map[string]string{
		"file.yml": `foo: hi`,
	}
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				{Inline: &v1alpha1.AppFetchInline{Paths: fetchInline}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)
	pkgr := v1alpha12.PackageRepository{}

	crdApp := NewCRDApp(&app, &pkgr, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	// Expected app status has no inspect on status
	// since the app deployment was not attempted
	expectedStatus := v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.Condition{{
				Type:    v1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Invalid package repository content: must contain 'packages/' directory but did not",
			}},
			ObservedGeneration:  0,
			FriendlyDescription: "Reconcile failed: Invalid package repository content: must contain 'packages/' directory but did not",
			UsefulErrorMessage:  "Invalid package repository content: must contain 'packages/' directory but did not",
		},
		Fetch: &v1alpha1.AppStatusFetch{
			ExitCode: 0,
		},
		Template: &v1alpha1.AppStatusTemplate{
			Error:    "Invalid package repository content: must contain 'packages/' directory but did not",
			ExitCode: -1,
		},
	}

	assertStatusEqual(t, crdApp, expectedStatus)
}

func assertStatusEqual(t *testing.T, crdApp *CRDApp, expectedStatus v1alpha1.AppStatus) {
	// Unset time for assertions
	// No need to assert on stderr as its captured elsewhere
	if crdApp.app.Status().Fetch != nil {
		crdApp.app.Status().Fetch.StartedAt = metav1.Time{}
		crdApp.app.Status().Fetch.UpdatedAt = metav1.Time{}
		crdApp.app.Status().Fetch.Stdout = ""
	}
	if crdApp.app.Status().Template != nil {
		crdApp.app.Status().Template.UpdatedAt = metav1.Time{}
		crdApp.app.Status().Template.Stderr = ""
	}

	gs := crdApp.app.Status().GenericStatus
	assert.Equal(t, expectedStatus.GenericStatus, gs)

	assert.Equal(t, expectedStatus.Fetch, crdApp.app.Status().Fetch)
	assert.Equal(t, expectedStatus.Template, crdApp.app.Status().Template)
}
