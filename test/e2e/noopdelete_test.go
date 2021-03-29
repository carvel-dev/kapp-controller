// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
)

func Test_NoopDelete_DeletesAfterServiceAccountDeleted(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	sas := ServiceAccounts{env.Namespace}

	appName := "noop-delete-app-with-no-sa"
	cfgMapName := "noop-configmap-with-no-sa"
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
  noopDelete: true
  fetch:
    - inline:
        paths:
          file.yml: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: %s
            data:
              key: value
  template:
    - ytt: {}
  deploy:
    - kapp: {}`, appName, cfgMapName) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	cleanUpConfigMap := func() {
		kubectl.Run([]string{"delete", "configmap", cfgMapName})
	}

	cleanUpApp()
	defer cleanUpApp()
	defer cleanUpConfigMap()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName}, RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("delete Service Account and App", func() {
		kubectl.Run([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa"})
		cleanUpApp()
	})

	logger.Section("check ConfigMap still exists after delete", func() {
		kubectl.Run([]string{"get", "configmap", cfgMapName})
	})
}
