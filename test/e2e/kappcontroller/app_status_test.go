// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func Test_AppStatus_DisplaysUsefulErrorMessage_ForDeploymentFailure(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "useful-err-app-deploy"
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - inline:
        paths:
          file.yml: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
  template:
    - ytt: {}
  deploy:
    - kapp:
        inspect: {}
        intoNs: does-not-exist`, name) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy", func() {
		out, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml), AllowError: true})
		// it's supposed to error, but it's also supposed to create the three supporting resources successfully:
		assert.Error(t, err)
		assert.Contains(t, out, "ok: reconcile role/kappctrl-e2e-ns-role")
		assert.Contains(t, out, "ok: reconcile rolebinding/kappctrl-e2e-ns-role-binding")
		assert.Contains(t, out, "ok: reconcile serviceaccount/kappctrl-e2e-ns-sa")
	})

	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})
	assert.Greater(t, len(out), 1000) // the output yaml should be non-trivial (observed len ~2.7k)

	var cr v1alpha1.App
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	expectedStatus := v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.Condition{{
				Type:    v1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Deploying: Error (see .status.usefulErrorMessage for details)",
			}},
			ObservedGeneration:  1,
			FriendlyDescription: "Reconcile failed: Deploying: Error (see .status.usefulErrorMessage for details)",
			UsefulErrorMessage:  "kapp: Error: Checking existence of resource configmap/configmap (v1) namespace: does-not-exist:\n  API server says: configmaps \"configmap\" is forbidden:\n    User \"system:serviceaccount:" + env.Namespace + ":kappctrl-e2e-ns-sa\" cannot get resource \"configmaps\" in API group \"\" in the namespace \"does-not-exist\" (reason: Forbidden)",
		},
		Deploy: &v1alpha1.AppStatusDeploy{
			ExitCode: 1,
			Finished: true,
			Error:    "Deploying: Error (see .status.usefulErrorMessage for details)",
		},
		Fetch: &v1alpha1.AppStatusFetch{
			ExitCode: 0,
		},
		Inspect: &v1alpha1.AppStatusInspect{
			ExitCode: 0,
		},
		Template: &v1alpha1.AppStatusTemplate{
			ExitCode: 0,
		},
		ConsecutiveReconcileSuccesses: 0,
		ConsecutiveReconcileFailures:  1,
	}

	{
		// deploy
		cr.Status.Deploy.StartedAt = metav1.Time{}
		cr.Status.Deploy.UpdatedAt = metav1.Time{}
		cr.Status.Deploy.Stdout = ""
		cr.Status.Deploy.Stderr = ""
		cr.Status.Deploy.KappDeployStatus = nil

		// fetch
		cr.Status.Fetch.StartedAt = metav1.Time{}
		cr.Status.Fetch.UpdatedAt = metav1.Time{}
		cr.Status.Fetch.Stdout = ""

		// inspect
		cr.Status.Inspect.UpdatedAt = metav1.Time{}
		cr.Status.Inspect.Stdout = ""

		// template
		cr.Status.Template.UpdatedAt = metav1.Time{}
		cr.Status.Template.Stderr = ""

	}

	require.Equal(t, expectedStatus, cr.Status)
}

func Test_AppStatus_Displays_Deploy_Status(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "app-deploy"
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - inline:
        paths:
          file.yml: |
            apiVersion: rbac.authorization.k8s.io/v1
            kind: Role
            metadata:
              name: test-app-status-role
  template:
    - ytt: {}
  deploy:
    - kapp:
        intoNs: %s`, name, env.Namespace) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})
	assert.Greater(t, len(out), 1000) // the output yaml should be non-trivial (observed len ~2.7k)

	var cr v1alpha1.App
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	expectedAppLabel := kapp.Run([]string{"label", "-a", name + ".app", "--tty=false"})

	expectedStatus := v1alpha1.AppStatus{
		Deploy: &v1alpha1.AppStatusDeploy{
			KappDeployStatus: &v1alpha1.KappDeployStatus{
				AssociatedResources: v1alpha1.AssociatedResources{
					Label:      expectedAppLabel,
					Namespaces: []string{"kappctrl-test"},
					GroupKinds: []metav1.GroupKind{
						{Group: "rbac.authorization.k8s.io", Kind: "Role"},
					},
				},
			},
		},
	}

	require.Equal(t, expectedStatus.Deploy.KappDeployStatus, cr.Status.Deploy.KappDeployStatus)

	appYaml = fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - inline:
        paths:
          file.yml: |
            --- 
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
            ---
            apiVersion: batch/v1
            kind: Job
            metadata:
              name: fail
            spec:
              template:
                spec:
                  containers:
                  - name: fail
                    image: busybox
                    command: ["/bin/sh",  "-c", "exit 1"]
                  restartPolicy: Never
              backoffLimit: 1
  template:
    - ytt: {}
  deploy:
    - kapp:
        intoNs: %s`, name, env.Namespace) + sas.ForNamespaceYAML()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml), AllowError: true})
	})

	out = kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})
	assert.Greater(t, len(out), 1000) // the output yaml should be non-trivial (observed len ~2.7k)

	var updateCr v1alpha1.App
	err = yaml.Unmarshal([]byte(out), &updateCr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	expectedStatus = v1alpha1.AppStatus{
		Deploy: &v1alpha1.AppStatusDeploy{
			KappDeployStatus: &v1alpha1.KappDeployStatus{
				AssociatedResources: v1alpha1.AssociatedResources{
					Label:      expectedAppLabel,
					Namespaces: []string{"kappctrl-test"},
					GroupKinds: []metav1.GroupKind{
						{Group: "", Kind: "ConfigMap"},
						{Group: "batch", Kind: "Job"},
						{Group: "rbac.authorization.k8s.io", Kind: "Role"},
					},
				},
			},
		},
	}

	require.Equal(t, expectedStatus.Deploy.KappDeployStatus, updateCr.Status.Deploy.KappDeployStatus)
}
