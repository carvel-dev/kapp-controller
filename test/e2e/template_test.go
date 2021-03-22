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

func Test_YttTemplate_UsesFileMarks(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	name := "configmap-with-non-yml-ext-file"
	// This App's ConfigMap is in a non yaml file
	// named file. In order for ytt to render the YAML,
	// the file needs a file-mark to denote it is a
	// YAML file.
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
          file: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
  template:
    - ytt:
        paths:
          - file
        fileMarks:
          - file:type=yaml-plain
  deploy:
    - kapp: {}
`, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, RunOpts{StdinReader: strings.NewReader(appYaml)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})
		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := v1alpha1.AppStatus{
			Conditions: []v1alpha1.AppCondition{{
				Type:   v1alpha1.ReconcileSucceeded,
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
			ConsecutiveReconcileSuccesses: 1,
			ObservedGeneration:            1,
			FriendlyDescription:           "Reconcile succeeded",
		}

		{
			// deploy
			cr.Status.Deploy.StartedAt = metav1.Time{}
			cr.Status.Deploy.UpdatedAt = metav1.Time{}
			cr.Status.Deploy.Stdout = ""

			// inspect
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""

			// fetch
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}
			cr.Status.Fetch.Stdout = ""
		}

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, cr.Status)
		}
	})
}
