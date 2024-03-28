// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestGitHttpsPublic(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	yaml1 := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-git-https-public
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - git:
      url: https://github.com/k14s/kapp
      ref: origin/develop
      subPath: examples/gitops/guestbook
  template:
  - ytt: {}
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, env.Namespace) + sas.ForNamespaceYAML()

	name := "test-git-https-public"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1),
				OnErrKubectl: []string{"get", "app", "-oyaml"},
			})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App

		err := yaml.Unmarshal([]byte(out), &cr)
		require.NoError(t, err)

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
			require.Contains(t, cr.Status.Deploy.Stdout, "Wait to:")
			cr.Status.Deploy.StartedAt = metav1.Time{}
			cr.Status.Deploy.UpdatedAt = metav1.Time{}
			cr.Status.Deploy.Stdout = ""
			cr.Status.Deploy.KappDeployStatus = nil

			// fetch
			require.Contains(t, cr.Status.Fetch.Stdout, "kind: LockConfig")
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}
			cr.Status.Fetch.Stdout = ""

			// inspect
			require.Contains(t, cr.Status.Inspect.Stdout, "Resources in app 'test-git-https-public.app'")
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		require.Equal(t, expectedStatus, cr.Status)
	})
}

func TestGitSshPrivate(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	yaml1 := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-git-ssh-private
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - git:
      url: git@git-server.%s.svc.cluster.local:/git-server/repos/myrepo.git
      ref: origin/master
      secretRef:
        name: git-private-key
  template:
  - ytt: {}
  deploy:
  - kapp:
      inspect: {}
      intoNs: %s
`, env.Namespace, env.Namespace) + sas.ForNamespaceYAML()

	name := "test-git-ssh-private"
	gitServerName := "test-git-server"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", gitServerName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy git server", func() {
		kapp.Run([]string{"deploy", "-f", "../assets/git-server.yml", "-a", gitServerName})
	})

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{
				IntoNs:       true,
				StdinReader:  strings.NewReader(yaml1),
				OnErrKubectl: []string{"get", "app/test-git-ssh-private", "-oyaml"}})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App

		err := yaml.Unmarshal([]byte(out), &cr)
		require.NoError(t, err)

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
			require.Contains(t, cr.Status.Deploy.Stdout, "Wait to:")
			cr.Status.Deploy.StartedAt = metav1.Time{}
			cr.Status.Deploy.UpdatedAt = metav1.Time{}
			cr.Status.Deploy.Stdout = ""
			cr.Status.Deploy.KappDeployStatus = nil

			// fetch
			require.Contains(t, cr.Status.Fetch.Stdout, "kind: LockConfig")
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}
			cr.Status.Fetch.Stdout = ""

			// inspect
			require.Contains(t, cr.Status.Inspect.Stdout, "Resources in app 'test-git-ssh-private.app'")
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		require.Equal(t, expectedStatus, cr.Status)
		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("Status is not same: %#v vs %#v", expectedStatus, cr.Status)
		}
	})
}
