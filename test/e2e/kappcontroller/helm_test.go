// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestHelm(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	expectedStatus := &v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.Condition{{
				Type:   v1alpha1.ReconcileSucceeded,
				Status: corev1.ConditionTrue,
			}},
			ObservedGeneration:  1,
			FriendlyDescription: "Reconcile succeeded",
		},
		Deploy: &v1alpha1.AppStatusDeploy{
			ExitCode: 0,
			Finished: true,
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
		ConsecutiveReconcileSuccesses: 1,
	}

	helmV3YAML := fmt.Sprintf(`
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-helm
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - helmChart:
      name: redis
      version: "12.10.1"
      repository:
        url: https://charts.bitnami.com/bitnami
  template:
  - helmTemplate:
      valuesFrom:
      - secretRef:
          name: test-helm-values
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
      delete:
        rawOptions: ["--apply-ignored=true"]
---
apiVersion: v1
kind: Secret
metadata:
  name: test-helm-values
stringData:
  data.yml: |
    password: "1234567891234"
`, env.Namespace) + sas.ForNamespaceYAML()

	tests := []struct {
		desc           string
		appCRName      string
		deploymentYAML string
		expectedStatus *v1alpha1.AppStatus
	}{
		{
			"Helm v3 (chart spec v2) deployment",
			"test-helm",
			helmV3YAML,
			expectedStatus,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			cleanUp := func() {
				kapp.Run([]string{"delete", "-a", tc.appCRName})
			}
			cleanUp()
			defer cleanUp()

			logger.Section("deploy", func() {
				kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", tc.appCRName},
					e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(tc.deploymentYAML)})

				out := kapp.Run([]string{"inspect", "-a", tc.appCRName, "--raw", "--tty=false", "--filter-kind=App"})

				var cr v1alpha1.App

				err := yaml.Unmarshal([]byte(out), &cr)
				if err != nil {
					t.Fatalf("Failed to unmarshal: %s", err)
				}

				{
					// deploy
					if !strings.Contains(cr.Status.Deploy.Stdout, "Wait to:") {
						t.Fatalf("Expected non-empty deploy output: '%s'", cr.Status.Deploy.Stdout)
					}
					cr.Status.Deploy.StartedAt = metav1.Time{}
					cr.Status.Deploy.UpdatedAt = metav1.Time{}
					cr.Status.Deploy.Stdout = ""

					// fetch
					if !strings.Contains(cr.Status.Fetch.Stdout, "kind: LockConfig") {
						t.Fatalf("Expected non-empty fetch output: '%s'", cr.Status.Fetch.Stdout)
					}
					cr.Status.Fetch.StartedAt = metav1.Time{}
					cr.Status.Fetch.UpdatedAt = metav1.Time{}
					cr.Status.Fetch.Stdout = ""

					// inspect
					if !strings.Contains(cr.Status.Inspect.Stdout, fmt.Sprintf("Resources in app '%s-ctrl'", tc.appCRName)) {
						t.Fatalf("Expected non-empty inspect output: '%s'", cr.Status.Inspect.Stdout)
					}
					cr.Status.Inspect.UpdatedAt = metav1.Time{}
					cr.Status.Inspect.Stdout = ""

					// template
					cr.Status.Template.UpdatedAt = metav1.Time{}
					cr.Status.Template.Stderr = ""
				}

				if !reflect.DeepEqual(*tc.expectedStatus, cr.Status) {
					t.Fatalf("Status is not same: %#v vs %#v", expectedStatus, cr.Status)
				}
			})
		})
	}
}

func TestHelmValuesOverStdin(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "test-helm-values-stdin"
	config := fmt.Sprintf(`
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-helm-values-stdin
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - helmChart:
      name: redis
      version: "12.10.1"
      repository:
        url: https://charts.bitnami.com/bitnami
  - inline:
      paths:
        data.yml: |
          global:
            redis:
              password: my-secret-password
  template:
  - ytt:
      paths: ["1/data.yml"]
  - helmTemplate:
      path: "0/"
      valuesFrom:
      - path: "-"
  - ytt:
      inline:
        paths:
          check.yml: |
            #@ load("@ytt:overlay", "overlay")
            #@ load("@ytt:base64", "base64")
            #@ load("@ytt:assert", "assert")

            #@ def check_password(l,_):
            #@   actual_pass = base64.decode(l)
            #@   if actual_pass != "my-secret-password":
            #@     assert.fail("Expected password '{}' to == 'my-secret-password'".format(actual_pass))
            #@   end
            #@ end

            #@overlay/match by=overlay.subset({"kind":"Secret"})
            ---
            data:
              #@overlay/assert via=check_password
              redis-password:

            #! to speed up deploy just remove everything
            #@overlay/match by=overlay.not_op(overlay.subset({"kind":"Secret"})),expects="1+"
            #@overlay/remove
            ---
  deploy:
  - kapp:
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// App CR will fail if ytt assertion fails
	kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
		e2e.RunOpts{StdinReader: strings.NewReader(config)})
}
