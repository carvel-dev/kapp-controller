// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"strings"
	"testing"

	"carvel.dev/kapp-controller/test/e2e"
)

func TestDeleteCancelsDeploys(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	appYaml := `---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: delete-while-deploying
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        file.yml: |
          apiVersion: v1
          kind: Pod
          metadata:
            name: sleep
          spec:
            containers:
             - name: nginx
               image: nginx
               readinessProbe:
                 httpGet:
                   port: 8080
            terminationGracePeriodSeconds: 0
  template:
  - ytt: {}
  deploy:
  - kapp: {}
` + sas.ForNamespaceYAML()

	name := "delete-while-deploying"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	t.Cleanup(cleanUp)

	logger.Section("begin deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--wait=false"},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(appYaml)})

		waitForDeployToStart(name, kapp, t)
	})

	logger.Section("delete", func() {
		kapp.RunWithOpts([]string{"delete", "-a", name, "--filter-kind", "App"}, e2e.RunOpts{})
	})
}
