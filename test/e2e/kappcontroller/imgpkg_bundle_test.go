// Copyright 2021 VMware, Inc.
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

func Test_FetchAndDeployImgpkgBundle_Successfully(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	// contents for kappctrl-e2e-bundle
	// available in test/e2e/assets/kappctrl-e2e-bundle
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-bundle-app
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - imgpkgBundle:
        image: k8slt/kappctrl-e2e-bundle
  template:
  - kbld:
      paths:
      - .imgpkg/images.yml
      - config.yml
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	name := "test-bundle-app"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(appYaml)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := v1alpha1.AppStatus{
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

		{
			// deploy
			cr.Status.Deploy.StartedAt = metav1.Time{}
			cr.Status.Deploy.UpdatedAt = metav1.Time{}
			cr.Status.Deploy.Stdout = ""
			cr.Status.Deploy.KappDeployStatus = nil

			// fetch
			if !strings.Contains(cr.Status.Fetch.Stdout, "- imgpkgBundle") {
				t.Fatalf("Expected to find imgpkgBundle contents in fetch stdout but got:\n%s", cr.Status.Fetch.Stdout)
			}
			if !strings.Contains(cr.Status.Fetch.Stdout, "image: index.docker.io/k8slt/kappctrl-e2e-bundle@sha256:83f86234f68a980490ec66f2d347ad4c8148713073c0993760b8eaaef3eb48d7") {
				t.Fatalf("Expected to find imgpkgBundle contents in fetch stdout but got:\n%s", cr.Status.Fetch.Stdout)
			}
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}
			cr.Status.Fetch.Stdout = ""

			// inspect
			if !strings.Contains(cr.Status.Inspect.Stdout, "simple-app") && !strings.Contains(cr.Status.Inspect.Stdout, "Succeeded") {
				t.Fatalf("Expected to find simple-app resources created but got:\n%s", cr.Status.Inspect.Stdout)
			}
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("\nStatus is not same:\nExpected:\n%#v\nGot:\n%#v\n", expectedStatus, cr.Status)
		}
	})
}
