// Copyright 2020 VMware, Inc.
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

func TestHelm(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	expectedStatus := &v1alpha1.AppStatus{
		GenericStatus: v1alpha1.GenericStatus{
			Conditions: []v1alpha1.AppCondition{{
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

	helmV2YAML := fmt.Sprintf(`
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-helm-v2
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - helmChart:
      name: redis
      # Chart version v1, DEPRECATED 
      version: "11.3.4"
      repository:
        url: https://charts.bitnami.com/bitnami
  template:
  - helmTemplate:
      valuesFrom:
      - secretRef:
          name: test-helm-values
  deploy:
  - kapp:
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
			"Helm v2 (chart spec v1) deployment",
			"test-helm-v2",
			helmV2YAML,
			expectedStatus,
		},
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
					RunOpts{IntoNs: true, StdinReader: strings.NewReader(tc.deploymentYAML)})

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
