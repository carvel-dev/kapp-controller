// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_AppStatus_DisplaysUsefulErrorMessage_ForDeploymentFailure(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

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
        intoNs: does-not-exist`, name) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, RunOpts{StdinReader: strings.NewReader(appYaml), AllowError: true})
	})

	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

	var cr v1alpha1.App
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %s", err)
	}

	expectedStatus := v1alpha1.AppStatus{
		Conditions: []v1alpha1.AppCondition{{
			Type:   v1alpha1.ReconcileFailed,
			Status: corev1.ConditionTrue,
			Message: "Deploying: exit status 1",
		}},
		Deploy: &v1alpha1.AppStatusDeploy{
			ExitCode: 1,
			Finished: true,
			Error: "Deploying: exit status 1",
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
		ObservedGeneration:            1,
		FriendlyDescription:           "Reconcile failed: Deploying: exit status 1",
		UsefulErrorMessage:            "kapp: Error: Checking existance of resource configmap/configmap (v1) namespace: does-not-exist: configmaps \"configmap\" is forbidden:\n  User \"system:serviceaccount:" + env.Namespace + ":kappctrl-e2e-ns-sa\" cannot get resource \"configmaps\" in API group \"\" in the namespace \"does-not-exist\" (reason: Forbidden)",
	}

	{
		// deploy
		cr.Status.Deploy.StartedAt = metav1.Time{}
		cr.Status.Deploy.UpdatedAt = metav1.Time{}
		cr.Status.Deploy.Stdout = ""
		cr.Status.Deploy.Stderr = ""

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

	if !reflect.DeepEqual(expectedStatus, cr.Status) {
		t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, cr.Status)
	}
}
