// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"carvel.dev/kapp-controller/test/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NamespaceDelete_AppWithResourcesInSameNamespace(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	nsName := "ns-delete"
	name := "resources-in-same-namespace"

	namespaceYAML := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %v`, nsName)

	appYAML := fmt.Sprintf(`---
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
           namespace: %s
          data:
           key: value
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, nsName) + e2e.ServiceAccounts{nsName}.ForNamespaceYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", nsName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("create namespace and deploy App", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", nsName, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(namespaceYAML)})
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--into-ns", nsName}, e2e.RunOpts{StdinReader: strings.NewReader(appYAML)})
	})

	logger.Section("delete namespace", func() {
		kubectl.Run([]string{"delete", "ns", nsName, "--timeout=1m"})
	})
}

func Test_NamespaceDelete_AppWithResourcesInDifferentTerminatingNamespaces(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	nsName1 := "ns-delete-1"
	nsName2 := "ns-delete-2"
	nsApp := "testnamespaces"
	name := "resources-in-different-namespaces"

	namespaceTemplate := `---
apiVersion: v1
kind: Namespace
metadata:
  name: %v`
	namespaceYAML := fmt.Sprintf(namespaceTemplate, nsName1) + "\n" + fmt.Sprintf(namespaceTemplate, nsName2)

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
           namespace: %s
          data:
           key: value
          ---
          apiVersion: v1
          kind: ConfigMap
          metadata:
           name: configmap
           namespace: %s
          data:
           key: value
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, nsName1, nsName2) + e2e.ServiceAccounts{nsName1}.ForClusterYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpTestNamespace := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", nsApp})
	}

	cleanUp()
	defer cleanUp()
	defer cleanUpTestNamespace()

	logger.Section("create namespace and deploy App", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", nsApp, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(namespaceYAML)})
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--into-ns", nsName1}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("delete namespace", func() {
		// delete SA first to reduce flakiness, sometimes SA deletion happens after app is deleted
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa", "-n", nsName1},
			e2e.RunOpts{NoNamespace: true})
		kubectl.Run([]string{"delete", "ns", nsName1, nsName2, "--timeout=1m"})
	})
}

func Test_NamespaceDelete_AppWithClusterScopedResources(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	nsName := "ns-delete-1"
	name := "app-with-cluster-scoped-resource"
	inAppNsName := "delete-test-ns"

	namespaceYAML := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %v`, nsName)

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
          kind: Namespace
          metadata:
            name: %s
          ---
          apiVersion: v1
          kind: ConfigMap
          metadata:
           name: configmap
           namespace: %s
          data:
           key: value
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, inAppNsName, nsName) + e2e.ServiceAccounts{nsName}.ForClusterYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpTestNamespace := func() {
		kubectl.Run([]string{"delete", "ns", inAppNsName})
		kubectl.RunWithOpts([]string{"patch", "App", name, "--type=json", "--patch", `[{ "op": "replace", "path": "/spec/noopDelete", "value": true}]`,
			"-n", nsName}, e2e.RunOpts{NoNamespace: true})
		kapp.Run([]string{"delete", "-a", nsName})
	}

	cleanUp()
	defer cleanUp()
	defer cleanUpTestNamespace()

	logger.Section("create namespace and deploy App", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", nsName, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(namespaceYAML)})
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--into-ns", nsName}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("delete namespace", func() {
		// delete SA first to reduce flakiness, sometimes SA deletion happens after app is deleted
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa", "-n", nsName},
			e2e.RunOpts{NoNamespace: true})
		_, err := kubectl.RunWithOpts([]string{"delete", "ns", nsName, "--timeout=30s"},
			e2e.RunOpts{AllowError: true})
		assert.Error(t, err, "Expected to get time out error, but did not")
	})
}

func Test_NamespaceDelete_AppWithWithOneNonTerminatingAffectedNamespace(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	nsName := "ns-delete"
	name := "resources-in-different-namespaces"

	namespaceYAML := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %v`, nsName)

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
           namespace: %s
          data:
           key: value
          ---
          apiVersion: v1
          kind: ConfigMap
          metadata:
           name: configmap
           namespace: %s
          data:
           key: value
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, nsName, env.Namespace) + e2e.ServiceAccounts{nsName}.ForClusterYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpTestNamespace := func() {
		kubectl.Run([]string{"delete", "configmap", "configmap"})
		kubectl.RunWithOpts([]string{"patch", "App", name, "--type=json", "--patch", `[{ "op": "replace", "path": "/spec/noopDelete", "value": true}]`,
			"-n", nsName}, e2e.RunOpts{NoNamespace: true})
		kapp.Run([]string{"delete", "-a", nsName})
	}

	cleanUp()
	defer cleanUp()
	defer cleanUpTestNamespace()

	logger.Section("create namespace and deploy App", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", nsName, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(namespaceYAML)})
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--into-ns", nsName}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("delete namespace", func() {
		// delete SA first to reduce flakiness, sometimes SA deletion happens after app is deleted
		kubectl.RunWithOpts([]string{"delete", "serviceaccount", "kappctrl-e2e-ns-sa", "-n", nsName},
			e2e.RunOpts{NoNamespace: true})
		_, err := kubectl.RunWithOpts([]string{"delete", "ns", nsName, "--timeout=30s"},
			e2e.RunOpts{AllowError: true})
		assert.Error(t, err, "Expected to get time out error, but did not")
	})
}

func Test_NamespaceDelete_AppWithResourcesInOutOfOrderTerminatingNamespaces(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	nsName1 := "ns-delete-1"
	nsName2 := "ns-delete-2"
	nsApp := "testnamespaces"
	name := "resources-in-different-namespaces"

	namespaceTemplate := `---
apiVersion: v1
kind: Namespace
metadata:
  name: %v`
	namespaceYAML := fmt.Sprintf(namespaceTemplate, nsName1) + "\n" + fmt.Sprintf(namespaceTemplate, nsName2)

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
           namespace: %s
          data:
           key: value
          ---
          apiVersion: v1
          kind: ConfigMap
          metadata:
           finalizers:
           - kubernetes
           name: configmap
           namespace: %s
          data:
           key: value
  template:
  - ytt: {}
  deploy:
  - kapp:
      delete:
        rawOptions: ["--wait-timeout=2s"]`, name, nsName1, nsName2) + e2e.ServiceAccounts{nsName1}.ForClusterYAML()

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUpTestNamespace := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", nsApp})
	}

	cleanUp()
	defer cleanUp()
	defer cleanUpTestNamespace()

	logger.Section("create namespace and deploy App", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", nsApp, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(namespaceYAML)})
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name, "--into-ns", nsName1}, e2e.RunOpts{StdinReader: strings.NewReader(appYaml)})
	})

	logger.Section("delete namespace", func() {
		// Deleting the namespace directly could delete the service account, so let's delete the App first
		_, err := kubectl.RunWithOpts([]string{"delete", "app", name, "-n", nsName1, "--timeout=10s"}, e2e.RunOpts{AllowError: true, NoNamespace: true})
		require.Error(t, err, "Expected to timeout")

		// Deletion should timeout as test-ns-2 is not being terminated and the cm present in it has a finalizer
		_, err = kubectl.RunWithOpts([]string{"delete", "ns", nsName1, "--timeout=10s"}, e2e.RunOpts{AllowError: true, NoNamespace: true})
		require.Error(t, err, "Expected to timeout")

		// At this point the service account should be deleted,
		// so removing the finalizer on cm in test-ns-2 wouldn't help
		kubectl.RunWithOpts([]string{"patch", "cm", "configmap", "--type=json", "--patch", `[{ "op": "remove", "path": "/metadata/finalizers"}]`,
			"-n", nsName2}, e2e.RunOpts{NoNamespace: true})

		kubectl.RunWithOpts([]string{"get", "App", name, "-n", nsName1}, e2e.RunOpts{NoNamespace: true})

		// Deleting the second namespace now should do a noop delete for the App
		kubectl.RunWithOpts([]string{"delete", "ns", nsName2, "--timeout=1m"}, e2e.RunOpts{NoNamespace: true})

		_, err = kubectl.RunWithOpts([]string{"get", "App", name, "-n", nsName1}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		require.Error(t, err, "Expected not found error")
	})
}
