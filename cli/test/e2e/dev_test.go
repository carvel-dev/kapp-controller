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
	gitAppName := "dev-gitapp-test"
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

	gitAppYaml := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  namespace: kctrl-test
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - git:
      url: https://github.com/k14s/k8s-simple-app-example
      ref: origin/develop
      subPath: config-step-2-template
  template:
  - ytt: {}
  deploy:
  - kapp:
      intoNs: kctrl-test
`, gitAppName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", fmt.Sprintf("%s.app", appName)})
		kapp.Run([]string{"delete", "-a", fmt.Sprintf("%s.app", gitAppName)})
		kapp.Run([]string{"delete", "-a", saAppName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("create service accounts", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", saAppName, "-f", "-", "--debug"}, RunOpts{StdinReader: strings.NewReader(sa)})
	})

	logger.Section("dev deploy app", func() {
		kappCtrl.RunWithOpts([]string{"dev", "-f", "-", "-l"}, RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("inspect app resources", func() {
		out := kapp.Run([]string{"inspect", "-a", fmt.Sprintf("%s.app", appName), "--json"})
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

	logger.Section("dev deploy git app", func() {
		out, err := kappCtrl.RunWithOpts([]string{"dev", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(gitAppYaml)})
		fmt.Printf("\n\n Out: %s \n err: %+v ", out, err)
	})

	logger.Section("inspect gitApp app resources", func() {
		out := kapp.Run([]string{"inspect", "-a", fmt.Sprintf("%s.app", gitAppName), "--json"})
		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"age":             "<replaced>",
				"kind":            "Deployment",
				"name":            "simple-app",
				"namespace":       "kctrl-test",
				"owner":           "kapp",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
			{
				"age":             "<replaced>",
				"kind":            "Endpoints",
				"name":            "simple-app",
				"namespace":       "kctrl-test",
				"owner":           "cluster",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
			{
				"age":             "<replaced>",
				"kind":            "Service",
				"name":            "simple-app",
				"namespace":       "kctrl-test",
				"owner":           "kapp",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
			{
				"age":             "<replaced>",
				"kind":            "ReplicaSet",
				"name":            "simple-app",
				"namespace":       "kctrl-test",
				"owner":           "cluster",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
			{
				"age":             "<replaced>",
				"kind":            "Pod",
				"name":            "simple-app",
				"namespace":       "kctrl-test",
				"owner":           "cluster",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
			{
				"age":             "<replaced>",
				"kind":            "EndpointSlice",
				"name":            "simple-app",
				"namespace":       "kctrl-test",
				"owner":           "cluster",
				"reconcile_info":  "",
				"reconcile_state": "ok",
			},
		}

		require.Len(t, output.Tables[0].Rows, 6)

		deploymentItem := filterByKeyValuePair(output.Tables[0].Rows, "kind", "Deployment")
		require.Exactly(t, expectedOutputRows[0], replaceAgeAndSinceDeployed(deploymentItem)[0])

		endpointItem := filterByKeyValuePair(output.Tables[0].Rows, "kind", "Endpoints")
		require.Exactly(t, expectedOutputRows[1], replaceAgeAndSinceDeployed(endpointItem)[0])

		serviceItem := filterByKeyValuePair(output.Tables[0].Rows, "kind", "Service")
		require.Exactly(t, expectedOutputRows[2], replaceAgeAndSinceDeployed(serviceItem)[0])

		replicaSetItem := filterByKeyValuePair(output.Tables[0].Rows, "kind", "ReplicaSet")
		replicaSetItem[0]["name"] = "simple-app"
		require.Exactly(t, expectedOutputRows[3], replaceAgeAndSinceDeployed(replicaSetItem)[0])

		podItem := filterByKeyValuePair(output.Tables[0].Rows, "kind", "Pod")
		podItem[0]["name"] = "simple-app"
		require.Exactly(t, expectedOutputRows[4], replaceAgeAndSinceDeployed(podItem)[0])

		endpointSliceItem := filterByKeyValuePair(output.Tables[0].Rows, "kind", "EndpointSlice")
		endpointSliceItem[0]["name"] = "simple-app"
		require.Exactly(t, expectedOutputRows[5], replaceAgeAndSinceDeployed(endpointSliceItem)[0])
	})
}

func filterByKeyValuePair(slice []map[string]string, key, value string) []map[string]string {
	var filteredSlice []map[string]string

	for _, item := range slice {
		if item[key] == value {
			filteredSlice = append(filteredSlice, item)
		}
	}
	return filteredSlice
}
