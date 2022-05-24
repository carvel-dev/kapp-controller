// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"testing"
	"strings"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	uitest "github.com/cppforlife/go-cli-ui/ui/test"
)

func TestConfig_HTTPProxy(t *testing.T) {
	assert := assert.New(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller", logger}

	// Proxy configured in config-test/secret-config.yml
	logger.Section("inspect controller logs for propagation of proxy env vars", func() {
		// app name must match the app name being deployed in hack/deploy-test.sh
		out := kubectl.Run([]string{"logs", "deployment/kapp-controller"})

		assert.Contains(out, "http_proxy is enabled.", "expected log line detailing http_proxy is enabled")
		assert.Contains(out, "no_proxy is enabled.", "expected log line detailing no_proxy is enabled")
	})
}

func TestConfig_KappDeployRawOptions(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "global-kapp-deploy-raw-opts"
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
          file.yml: |
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              key: value
  template:
    - ytt: {}
  deploy:
    - kapp:
        rawOptions: ["--labels=local-lbl=local-lbl-val"]`, name) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpApp()
	defer cleanUpApp()

	// Global label is configured in config-test/secret-config.yml
	logger.Section("deploy and check that kc-test label is set", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})

		// Exactly one app should have global and local label combination
		out := kapp.Run([]string{"ls", "-A", "--filter-labels", "kc-test=kc-test-val,local-lbl=local-lbl-val", "--json"})
		resp := uitest.JSONUIFromBytes(t, []byte(out))
		assert.Equal(t, len(resp.Tables[0].Rows), 1)
	})
}
