// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func TestServiceAccountNotAllowed(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

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
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1), AllowError: true})
		if err == nil {
			t.Fatalf("Expected err, but was nil")
		}

		expectedErr := `cannot get resource "configmaps" in API group "" in the namespace "kube-system" (reason: Forbidden)`

		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("Expected forbidden error in deploy output, but was: %s", err.Error())
		}
	})

	logger.Section("deploy allowed resources", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml2)})
	})
}
