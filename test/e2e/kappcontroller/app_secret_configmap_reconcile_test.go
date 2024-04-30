// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"carvel.dev/kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func Test_AppReconcileOccurs_WhenSecretUpdated(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "configmap-with-secret"
	// syncPeriod set to 1 hour so that test
	// won't pass because of reconcile from time sync.
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  syncPeriod: 1h
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - inline:
        paths:
          file.yml: |
            #@ load("@ytt:data", "data")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              hello_msg: #@ data.values.hello_msg
  template:
  - ytt:
      inline:
        pathsFrom:
          - secretRef:
              name: simple-app-values
      paths:
        - file.yml
  deploy:
    - kapp: {}
---
apiVersion: v1
kind: Secret
metadata:
  name: simple-app-values
stringData:
  values.yml: |
    #@data/values
    ---
    hello_msg: original`, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("update secret", func() {
		updatedSecret := `
---
apiVersion: v1
kind: Secret
metadata:
  name: simple-app-values
stringData:
  values.yml: |
    #@data/values
    ---
    hello_msg: updated`

		// Update secret
		kubectl.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(updatedSecret)})
	})

	logger.Section("check App uses new secret", func() {
		retry(t, 10, func() error {
			out := kubectl.Run([]string{"get", "configmap/configmap", "-o", "yaml"})

			var cm corev1.ConfigMap
			err := yaml.Unmarshal([]byte(out), &cm)
			if err != nil {
				return fmt.Errorf("failed to unmarshal: %s", err)
			}

			if cm.Data["hello_msg"] != "updated" {
				return fmt.Errorf("secret message was not updated to \"updated\"\nGot: %s", cm.Data["hello_msg"])
			}
			return nil
		})
	})
}

func Test_AppReconcileOccurs_WhenConfigMapUpdated(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	name := "configmap-with-configmap"
	// syncPeriod set to 1 hour so that test
	// won't pass because of reconcile from time sync.
	appYaml := fmt.Sprintf(`
---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  syncPeriod: 1h
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
    - inline:
        paths:
          file.yml: |
            #@ load("@ytt:data", "data")
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: configmap
            data:
              hello_msg: #@ data.values.hello_msg
  template:
  - ytt:
      inline:
        pathsFrom:
          - configMapRef:
              name: simple-app-values
      paths:
        - file.yml
  deploy:
    - kapp: {}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: simple-app-values
data:
  values.yml: |
    #@data/values
    ---
    hello_msg: original`, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("update configmap", func() {
		updatedConfigMap := `
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: simple-app-values
data:
  values.yml: |
    #@data/values
    ---
    hello_msg: updated`

		// Update configmap
		kubectl.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(updatedConfigMap)})
	})

	logger.Section("check App uses new configmap", func() {
		retry(t, 10, func() error {
			out := kubectl.Run([]string{"get", "configmap/configmap", "-o", "yaml"})

			var cm corev1.ConfigMap
			err := yaml.Unmarshal([]byte(out), &cm)
			if err != nil {
				return fmt.Errorf("failed to unmarshal: %s", err)
			}

			if cm.Data["hello_msg"] != "updated" {
				return fmt.Errorf("configmap message was not updated to \"updated\"\nGot: %s", cm.Data["hello_msg"])
			}
			return nil
		})
	})
}

func retry(t *testing.T, maxRetries int, f func() error) {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = f()
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("retry failed after %d attempts: %v", maxRetries, err)
}
