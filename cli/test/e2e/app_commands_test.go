// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestAppE2E(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	sas := ServiceAccounts{env.Namespace}

	name := "kctrl-test-app"
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
        inspect: {}
        intoNs: kctrl-test`, name) + sas.ForNamespaceYAML()

	cleanUpApp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpApp()
	defer cleanUpApp()

	logger.Section("deploy", func() {
		out, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, RunOpts{StdinReader: strings.NewReader(appYaml)})

		require.NoError(t, err)
		require.Contains(t, out, "ok: reconcile app/kctrl-test-app")
		require.Contains(t, out, "ok: reconcile role/kappctrl-e2e-ns-role")
		require.Contains(t, out, "ok: reconcile rolebinding/kappctrl-e2e-ns-role-binding")
		require.Contains(t, out, "ok: reconcile serviceaccount/kappctrl-e2e-ns-sa")
	})

	logger.Section("list app", func() {
		out := kappCtrl.Run([]string{"app", "list", "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"age":          "<replaced>",
			"name":         "kctrl-test-app",
			"owner":        "",
			"since_deploy": "<replaced>",
			"status":       "Reconcile succeeded",
		}}
		require.Exactly(t, expectedOutputRows, replaceAgeAndSinceDeployed(output.Tables[0].Rows))
	})

	logger.Section("get app", func() {
		out := kappCtrl.Run([]string{"app", "get", "-a", name, "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":       "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"name":             "kctrl-test-app",
			"namespace":        "kctrl-test",
			"owner_references": "",
			"service_account":  "kappctrl-e2e-ns-sa",
			"status":           "Reconcile succeeded",
		}}
		require.Exactly(t, expectedOutputRows, replaceAgeAndSinceDeployed(output.Tables[0].Rows))
	})

	logger.Section("get app status", func() {
		out := kappCtrl.Run([]string{"app", "status", "-a", name})

		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "App reconciled")
	})

	logger.Section("pause app", func() {
		kappCtrl.Run([]string{"app", "pause", "-a", name})
	})

	logger.Section("get app", func() {
		out := kappCtrl.Run([]string{"app", "get", "-a", name, "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":       "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"name":             "kctrl-test-app",
			"namespace":        "kctrl-test",
			"owner_references": "",
			"service_account":  "kappctrl-e2e-ns-sa",
			"status":           "Canceled/paused",
		}}
		require.Exactly(t, expectedOutputRows, replaceAgeAndSinceDeployed(output.Tables[0].Rows))
	})

	logger.Section("kick app", func() {
		kappCtrl.Run([]string{"app", "kick", "-a", name})
	})

	logger.Section("get app", func() {
		out := kappCtrl.Run([]string{"app", "get", "-a", name, "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":       "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"name":             "kctrl-test-app",
			"namespace":        "kctrl-test",
			"owner_references": "",
			"service_account":  "kappctrl-e2e-ns-sa",
			"status":           "Reconcile succeeded",
		}}
		require.Exactly(t, expectedOutputRows, replaceAgeAndSinceDeployed(output.Tables[0].Rows))
	})

	logger.Section("delete app", func() {
		kappCtrl.Run([]string{"app", "delete", "-a", name})
	})

	logger.Section("list apps", func() {
		out := kappCtrl.Run([]string{"app", "list", "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}
		require.Exactly(t, expectedOutputRows, replaceAgeAndSinceDeployed(output.Tables[0].Rows))
	})

}

func replaceAgeAndSinceDeployed(result []map[string]string) []map[string]string {
	for i, row := range result {
		if len(row["age"]) > 0 {
			row["age"] = "<replaced>"
			row["since_deploy"] = "<replaced>"
		}
		result[i] = row
	}
	return result
}
