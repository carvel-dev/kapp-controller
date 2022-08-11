// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestDev(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, Logger{}}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}

	appName := "dev-test"
	saAppName := "dev-test"

	sa := ServiceAccounts{env.Namespace}.ForNamespaceYAML()
	appYaml := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
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
        intoNs: kctrl-test
`, appName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", fmt.Sprintf("%s-ctrl", appName)})
		kapp.Run([]string{"delete", "-a", fmt.Sprintf("%s-ctrl", saAppName)})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("create service accounts", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", saAppName, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(sa)})
	})

	logger.Section("dev deploy app", func() {
		kappCtrl.RunWithOpts([]string{"dev", "-f", "-", "-l"}, RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("inspect app resources", func() {
		out := kapp.Run([]string{"inspect", "-a", fmt.Sprintf("%s-ctrl", appName), "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"age":             "<replaced>",
				"kind":            "ConfigMap",
				"name":            "configmap",
				"namespace":       "kctrl-test",
				"owner":           "kapp",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
		}
		require.Exactly(t, expectedOutputRows, replaceAgeAndSinceDeployed(output.Tables[0].Rows))
	})
}
