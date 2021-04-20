// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

func TestServiceAccountNotAllowed(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	yaml1 := `
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-sa-not-allowed
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        config.yml: |
          kind: ConfigMap
          apiVersion: v1
          metadata:
            name: test-not-allowed
            namespace: kube-system #! <-- not allowed namespace
          data:
            not-allowed: ""
  template:
  - ytt: {}
  deploy:
  - kapp: {}
` + sas.ForNamespaceYAML()

	yaml2 := `
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-sa-not-allowed
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        config.yml: |
          kind: ConfigMap
          apiVersion: v1
          metadata:
            name: test-not-allowed
          data:
            not-allowed: ""
  template:
  - ytt: {}
  deploy:
  - kapp: {}
` + sas.ForNamespaceYAML()

	name := "test-service-account-not-allowed"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy forbidden resource", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1), AllowError: true})
		if err == nil {
			t.Fatalf("Expected err, but was nil")
		}

		if !strings.Contains(err.Error(), "Reconcile failed:  (message: Deploying: Error (see .status.usefulErrorMessage for details))") {
			t.Fatalf("Expected err to contain service account failure, but was: %s", err)
		}

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=App"})

		var cr v1alpha1.App

		err = yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedErr := `cannot get resource "configmaps" in API group "" in the namespace "kube-system" (reason: Forbidden)`

		if !strings.Contains(cr.Status.Deploy.Stderr, expectedErr) {
			t.Fatalf("Expected forbidden error in deploy output, but was: %#v", cr.Status.Deploy)
		}
	})

	logger.Section("deploy allowed resources", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml2)})
	})
}

func Test_AppDeletes_WhenServiceAccountDoesNotExist_AndNoAppResourcesDeployed(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}

	name := "test-sa-not-exist"
	appYamlNoSA := fmt.Sprintf(`
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
        config.yml: |
          kind: ConfigMap
          apiVersion: v1
          metadata:
            name: test-no-sa
  template:
  - ytt: {}
  deploy:
  - kapp: {}
`, name)

	cleanUp := func() {
		// Since no error is expected, this serves
		// as final assertion for test
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy App with non-existent serviceaccount", func() {
		stdout, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(appYamlNoSA), AllowError: true})
		if err == nil {
			t.Fatalf("Expected err, but was nil.\nStdout: %s", stdout)
		}

		if !strings.Contains(err.Error(), "Preparing kapp: Getting service account: serviceaccounts \"kappctrl-e2e-ns-sa\" not found") {
			t.Fatalf("Expected err to contain service account failure, but was: %s", err)
		}
	})
}
