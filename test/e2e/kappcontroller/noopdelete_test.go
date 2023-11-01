// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"carvel.dev/kapp-controller/test/e2e"
)

func Test_NoopDelete_DeletesAfterServiceAccountDeleted(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}
	name := "instl-pkg-noop-delete"
	cfgMapName := "configmap"

	appYaml := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
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
  - kapp: {}`, name, cfgMapName) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUpConfigMap := func() {
		kubectl.Run([]string{"delete", "configmap", cfgMapName})
	}

	cleanUpApp()
	defer cleanUpApp()
	defer cleanUpConfigMap()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "apps/" + name, "--timeout", "1m"})
	logger.Section("delete Service Account and App", func() {
		kubectl.Run([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa"})
		cleanUpApp()
	})

	logger.Section("check ConfigMap still exists after delete", func() {
		kubectl.Run([]string{"get", "configmap/" + cfgMapName})
	})
}
