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

func Test_NoKappInspect_IfNoDeployAttempted_AfterFetchFailure(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	// The url under fetch is invalid, which will cause this
	// app to never be deployed.
	yaml1 := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: simple-app
  namespace: %s
spec:
  serviceAccountName: default-ns-sa
  fetch:
  - http:
      url: i-dont-exist
  template:
  - ytt: {}
  deploy:
  - kapp: {}
`, env.Namespace) + sas.ForNamespaceYAML()

	name := "simple-app"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1), AllowError: true})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
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
			ObservedGeneration:           1,
			FriendlyDescription:          "Reconcile failed: Fetching resources: exit status 1",
		}

		cr.Status.Fetch.StartedAt = metav1.Time{}
		cr.Status.Fetch.UpdatedAt = metav1.Time{}
		cr.Status.Fetch.Stderr = ""

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("Status is not same:\n%#v\nvs\n%#v\n", expectedStatus, cr.Status)
		}

		// Assert deletion is successful after failed fetch
		kapp.RunWithOpts([]string{"delete", "-a", name}, RunOpts{})
	})
}
