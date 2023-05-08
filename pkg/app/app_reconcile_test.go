// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/k14s/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func Test_NoInspectReconcile_IfNoDeployAttempted(t *testing.T) {
	log := logf.Log.WithName("kc")
	var appMetrics = metrics.NewAppMetrics()

	// The url under fetch is invalid, which will cause this
	// app to fail before deploy.
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			UID:       uuid.NewUUID(),
			Name:      "simple-app",
			Namespace: "pkg-standalone",
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

	crdApp := NewCRDApp(&app, log, appMetrics, kappcs, fetchFac, tmpFac, deployFac, FakeComponentInfo{}, Opts{MinimumSyncPeriod: 30 * time.Second})
	_, err := crdApp.Reconcile(false)
	assert.Nil(t, err, "unexpected error with reconciling", err)

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
			ExitCode: 1,
		},
		ConsecutiveReconcileFailures: 1,
	}

	crdApp.app.Status().Fetch.StartedAt = metav1.Time{}
	crdApp.app.Status().Fetch.UpdatedAt = metav1.Time{}
	// No need to assert on stderr as its captured elsewhere
	crdApp.app.Status().Fetch.Stderr = ""

	assert.Equal(t, expectedStatus, crdApp.app.Status())
}

func Test_NoInspectReconcile_IfInspectNotEnabled(t *testing.T) {
	log := logf.Log.WithName("kc")
	var appMetrics = metrics.NewAppMetrics()

	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			UID:       uuid.NewUUID(),
			Name:      "simple-app",
			Namespace: "pkg-standalone",
		},
		Spec: v1alpha1.AppSpec{
			Fetch: []v1alpha1.AppFetch{
				v1alpha1.AppFetch{Inline: &v1alpha1.AppFetchInline{Paths: map[string]string{"file.yml": `|
                apiVersion: v1
                kind: ConfigMap
                metadata:
                  name: configmap
                data:
                  key: value`}}},
			},
			Template: []v1alpha1.AppTemplate{
				v1alpha1.AppTemplate{Ytt: &v1alpha1.AppTemplateYtt{}},
			},
			Deploy: []v1alpha1.AppDeploy{
				v1alpha1.AppDeploy{Kapp: &v1alpha1.AppDeployKapp{}},
			},
		},
	}

	k8scs := k8sfake.NewSimpleClientset()
	kappcs := fake.NewSimpleClientset()
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	crdApp := NewCRDApp(&app, log, appMetrics, kappcs, fetchFac, tmpFac, deployFac, FakeComponentInfo{}, Opts{MinimumSyncPeriod: 30 * time.Second})
	_, err := crdApp.Reconcile(false)
	assert.Nil(t, err, "unexpected error with reconciling", err)

	// Expected app status has no inspect on status
	// since it's not enabled
	expectedStatus := v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.Condition{{
				Type:    v1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Blocking for deploy: Updating app: apps.kappctrl.k14s.io \"simple-app\" not found",
			}},
			ObservedGeneration:  0,
			FriendlyDescription: "Reconcile failed: Blocking for deploy: Updating app: apps.kappctrl.k14s.io \"simple-app\" not found",
			UsefulErrorMessage:  "Blocking for deploy: Updating app: apps.kappctrl.k14s.io \"simple-app\" not found",
		},
		Fetch: &v1alpha1.AppStatusFetch{
			Stdout:   "apiVersion: vendir.k14s.io/v1alpha1\ndirectories:\n- contents:\n  - inline: {}\n    path: .\n  path: \"0\"\nkind: LockConfig\n",
			ExitCode: 0,
		},
		ConsecutiveReconcileFailures: 1,
		Template: &v1alpha1.AppStatusTemplate{
			ExitCode: 0,
		},
		Deploy: &v1alpha1.AppStatusDeploy{
			Finished: true,
			ExitCode: -1,
			Error:    "Blocking for deploy: Updating app: apps.kappctrl.k14s.io \"simple-app\" not found",
		},
	}

	crdApp.app.Status().Fetch.StartedAt = metav1.Time{}
	crdApp.app.Status().Fetch.UpdatedAt = metav1.Time{}
	crdApp.app.Status().Deploy.StartedAt = metav1.Time{}
	crdApp.app.Status().Deploy.UpdatedAt = metav1.Time{}
	crdApp.app.Status().Template.UpdatedAt = metav1.Time{}
	// No need to assert on stderr as its captured elsewhere
	crdApp.app.Status().Fetch.Stderr = ""

	assert.Equal(t, expectedStatus, crdApp.app.Status())
}

func Test_TemplateError_DisplayedInStatus_UsefulErrorMessageProperty(t *testing.T) {
	log := logf.Log.WithName("kc")
	var appMetrics = metrics.NewAppMetrics()

	fetchInline := map[string]string{
		"file.yml": `foo: #@ data.values.nothere`,
	}
	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			UID:       uuid.NewUUID(),
			Name:      "simple-app",
			Namespace: "pkg-standalone",
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
	fetchFac := fetch.NewFactory(k8scs, fetch.VendirOpts{}, exec.NewPlainCmdRunner())
	tmpFac := template.NewFactory(k8scs, fetchFac, false, exec.NewPlainCmdRunner())
	deployFac := deploy.NewFactory(k8scs, kubeconfig.NewKubeconfig(k8scs, log), nil, exec.NewPlainCmdRunner(), log)

	crdApp := NewCRDApp(&app, log, appMetrics, kappcs, fetchFac, tmpFac, deployFac, FakeComponentInfo{}, Opts{MinimumSyncPeriod: 30 * time.Second})
	_, err := crdApp.Reconcile(false)
	assert.Nil(t, err, "Unexpected error with reconciling", err)

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
		ConsecutiveReconcileFailures: 1,
	}

	// Unset time for assertions
	crdApp.app.Status().Fetch.StartedAt = metav1.Time{}
	crdApp.app.Status().Fetch.UpdatedAt = metav1.Time{}
	crdApp.app.Status().Template.UpdatedAt = metav1.Time{}

	crdApp.app.Status().Fetch.Stdout = ""
	// No need to assert on stderr as its captured elsewhere
	crdApp.app.Status().Template.Stderr = ""

	assert.True(t,
		reflect.DeepEqual(expectedStatus, crdApp.app.Status()),
		fmt.Sprintf("Status is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, crdApp.app.Status()))
}

type FakeComponentInfo struct {
	KCVersion       semver.Version
	KCVersionCount  *int
	K8sVersion      semver.Version
	K8sVersionCount *int
	K8sAPIs         []string
	K8sAPIsCount    *int
}

func (f FakeComponentInfo) KubernetesAPIs() ([]string, error) {
	*f.K8sAPIsCount++
	return f.K8sAPIs, nil
}

func (f FakeComponentInfo) KappControllerVersion() (semver.Version, error) {
	*f.KCVersionCount++
	return f.KCVersion, nil
}

func (f FakeComponentInfo) KubernetesVersion(_ string, _ *v1alpha1.AppCluster, _ *metav1.ObjectMeta) (semver.Version, error) {
	*f.K8sVersionCount++
	return f.K8sVersion, nil
}
