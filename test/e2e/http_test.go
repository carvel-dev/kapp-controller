 // Copyright 2020 VMware, Inc.
 // SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

func TestHTTP(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	yaml1 := fmt.Sprintf(`
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-http
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - http:
      url: https://raw.githubusercontent.com/k14s/kapp/db2cc63e12e988235eb8815af8edc0ca0cfaa79c/examples/simple-app-example/config-1.yml
      sha256: 71b0a2dceb6bb6e7a519d0587abd9400a265cabbe2c084f783df94a92db6d980
  template:
  - ytt: {}
  deploy:
  - kapp:
      intoNs: %s
`, env.Namespace)+sas.ForNamespaceYAML()

	name := "test-http"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App

		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := v1alpha1.AppStatus{
			Conditions: []v1alpha1.AppCondition{{
				Type: v1alpha1.ReconcileSucceeded,
				Status: corev1.ConditionTrue,
			}},
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
			ObservedGeneration: 1,
			FriendlyDescription: "Reconcile succeeded",
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
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}

			// inspect
			if !strings.Contains(cr.Status.Inspect.Stdout, "Resources in app 'test-http-ctrl'") {
				t.Fatalf("Expected non-empty inspect output: '%s'", cr.Status.Inspect.Stdout)
			}
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("Status is not same: %#v vs %#v", expectedStatus, cr.Status)
		}
	})
}
