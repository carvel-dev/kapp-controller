// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"reflect"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
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
	fetchFac := fetch.NewFactory(k8scs, nil)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	crdApp := NewCRDApp(&app, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	// Expected app status has no inspect on status
	// since the app deployment was not attempted
	expectedStatus := v1alpha1.AppStatus{
		Conditions: []v1alpha1.AppCondition{{
			Type:    v1alpha1.ReconcileFailed,
			Status:  corev1.ConditionTrue,
			Message: "Fetching resources: exit status 1",
		}},
		Fetch: &v1alpha1.AppStatusFetch{
			Error:    "Fetching resources: exit status 1",
			ExitCode: 1,
		},
		ConsecutiveReconcileFailures: 1,
		ObservedGeneration:           0,
		FriendlyDescription:          "Reconcile failed: Fetching resources: exit status 1",
		UsefulErrorMessage:           "Error: Syncing directory '0': Syncing directory '.' with HTTP contents: Downloading URL: Initiating URL download: Get i-dont-exist: unsupported protocol scheme \"\"\n",
	}

	// Unset time for assertions
	crdApp.app.Status().Fetch.StartedAt = metav1.Time{}
	crdApp.app.Status().Fetch.UpdatedAt = metav1.Time{}
	// No need to assert on stderr as its captured elsewhere
	crdApp.app.Status().Fetch.Stderr = ""

	if !reflect.DeepEqual(expectedStatus, crdApp.app.Status()) {
		t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, crdApp.app.Status())
	}
}

func Test_TemplateError_DisplayedInStatus_UsefulErrorMessageProperty(t *testing.T) {
	log := logf.Log.WithName("kc")

	fetchInline := map[string]string{
		"file.yml": `# comment
foo: bar`,
	}
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simple-app",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{Paths: fetchInline}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{Paths: []string{"file.yml"}}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, nil)
	tmpFac := template.NewFactory(k8scs, fetchFac)
	deployFac := deploy.NewFactory(k8scs)

	crdApp := NewCRDApp(&app, log, kappcs, fetchFac, tmpFac, deployFac)
	_, err := crdApp.Reconcile(false)
	if err != nil {
		t.Fatalf("Unexpected error with reconciling: %v", err)
	}

	// Expected app status has no inspect on status
	// since the app deployment was not attempted
	expectedStatus := v1alpha1.AppStatus{
		Conditions: []v1alpha1.AppCondition{{
			Type:    v1alpha1.ReconcileFailed,
			Status:  corev1.ConditionTrue,
			Message: "Templating dir: exit status 1",
		}},
		Fetch: &v1alpha1.AppStatusFetch{
			ExitCode: 0,
		},
		Template: &v1alpha1.AppStatusTemplate{
			Error:    "Templating dir: exit status 1",
			ExitCode: 1,
		},
		ConsecutiveReconcileFailures: 1,
		ObservedGeneration:           0,
		FriendlyDescription:          "Reconcile failed: Templating dir: exit status 1",
		UsefulErrorMessage:           "ytt: Error: Non-ytt comment at line file.yml:1: '# comment':\n  Unrecognized comment type (expected '#@' or '#!'). (hint: if this is plain YAML — not a template — consider `--file-mark '<filename>:type=yaml-plain'`)\n",
	}

	// Unset time for assertions
	crdApp.app.Status().Fetch.StartedAt = metav1.Time{}
	crdApp.app.Status().Fetch.UpdatedAt = metav1.Time{}
	crdApp.app.Status().Template.UpdatedAt = metav1.Time{}

	crdApp.app.Status().Fetch.Stdout = ""
	// No need to assert on stderr as its captured elsewhere
	crdApp.app.Status().Template.Stderr = ""

	if !reflect.DeepEqual(expectedStatus, crdApp.app.Status()) {
		t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, crdApp.app.Status())
	}
}
