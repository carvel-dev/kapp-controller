// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func Test_AppDefaultNamespace(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "app-default-namespace"
	defaultNamespace := "default-namespace"
	defaultNamespaceApp := "default-namespace-app"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", defaultNamespaceApp})
	}
	cleanUp()
	defer cleanUp()

	namespaceYAML := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %s`, defaultNamespace)

	sas := e2e.ServiceAccounts{env.Namespace}
	appYAML := fmt.Sprintf(`---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  namespace: %s
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  defaultNamespace: %s
  fetch:
  - inline:
      paths:
        file.yml: |
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: my-cm
          data:
            key: value
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, env.Namespace, defaultNamespace) + sas.ForClusterYAML()

	kapp.RunWithOpts([]string{"deploy", "-a", defaultNamespaceApp, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(namespaceYAML)})

	_, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true, StdinReader: strings.NewReader(appYAML)})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// Assert that app resources are in defaultNamespace
	kubectl.RunWithOpts([]string{"get", "configmap", "my-cm", "-n", defaultNamespace}, e2e.RunOpts{NoNamespace: true})

	// Assert that kapp metaconfigmap is in app namespace
	kubectl.RunWithOpts([]string{"get", "configmap", name + ".app", "-n", env.Namespace}, e2e.RunOpts{NoNamespace: true})
}

func Test_AppDefaultNamespace_WithTargetCluster(t *testing.T) {
	targetClusterKubeconfig := os.Getenv("TEST_E2E_TARGET_CLUSTER_KUBECONFIG")
	if targetClusterKubeconfig == "" {
		t.Skip("Skipping test as target cluster kubeconfig is not set")
	}

	kubeconfigFile, err := os.CreateTemp("", "e2e-kubeconfig-*")
	assert.NoError(t, err)
	defer os.Remove(kubeconfigFile.Name())

	_, err = kubeconfigFile.Write([]byte(targetClusterKubeconfig))
	assert.NoError(t, err)

	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	name := "app-default-namespace-target-cluster"
	defaultNamespace := "e2e-default-namespace"
	clusterNamespace := "e2e-cluster-namespace"
	namespaceApp := "namespace-app"
	secretName := "e2e-kubeconfig-secret"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})

	}
	cleanUpTargetCluster := func() {
		kapp.RunWithOpts([]string{"delete", "-a", namespaceApp, "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})
	}

	cleanUp()
	cleanUpTargetCluster()
	defer cleanUp()
	defer cleanUpTargetCluster()

	namespaceYAML := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Namespace
metadata:
  name: %s`, defaultNamespace, clusterNamespace)

	secret := e2e.Secrets{secretName, env.Namespace, targetClusterKubeconfig}
	appYAML := `---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: %s
  namespace: %s
  annotations:
    kapp.k14s.io/change-group: "kappctrl-e2e.k14s.io/apps"
spec:
  cluster:
    namespace: %s
    kubeconfigSecretRef:
      name: %s
  defaultNamespace: %s
  fetch:
  - inline:
      paths:
        file.yml: |
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: my-cm
          data:
            key: value
  template:
  - ytt: {}
  deploy:
  - kapp: {}`

	// create test namespaces on target cluster
	kapp.RunWithOpts([]string{"deploy", "-a", namespaceApp, "-f", "-", "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true, StdinReader: strings.NewReader(namespaceYAML)})

	// Provide both _defaultNamespace_ and _cluster.namespace_
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true,
		StdinReader: strings.NewReader(fmt.Sprintf(appYAML, name, env.Namespace, clusterNamespace, secretName, defaultNamespace) + secret.ForTargetCluster())})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// Assert that app resources are in defaultNamespace
	kubectl.RunWithOpts([]string{"get", "configmap", "my-cm", "-n", defaultNamespace, "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	// Assert that kapp metaconfigmap is in cluster.namespace
	kubectl.RunWithOpts([]string{"get", "configmap", name + ".app", "-n", clusterNamespace, "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	cleanUp()

	// Provide only _cluster.namespace_
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true,
		StdinReader: strings.NewReader(fmt.Sprintf(appYAML, name, env.Namespace, clusterNamespace, secretName, "") + secret.ForTargetCluster())})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// Assert that app resources are in cluster.namespace
	kubectl.RunWithOpts([]string{"get", "configmap", "my-cm", "-n", clusterNamespace, "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	// Assert that kapp metaconfigmap is in cluster.namespace
	kubectl.RunWithOpts([]string{"get", "configmap", name + ".app", "-n", clusterNamespace, "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	cleanUp()

	// Provide only _defaultNamespace_
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true,
		StdinReader: strings.NewReader(fmt.Sprintf(appYAML, name, env.Namespace, "", secretName, defaultNamespace) + secret.ForTargetCluster())})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// Assert that app resources are in defaultNamespace
	kubectl.RunWithOpts([]string{"get", "configmap", "my-cm", "-n", defaultNamespace, "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	// Assert that kapp metaconfigmap is in kubeconfig preferred namespace (default)
	kubectl.RunWithOpts([]string{"get", "configmap", name + ".app", "-n", "default", "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	cleanUp()

	// Do not provide _defaultNamespace_ and _cluster.namespace_
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true,
		StdinReader: strings.NewReader(fmt.Sprintf(appYAML, name, env.Namespace, "", secretName, "") + secret.ForTargetCluster())})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// Assert that app resources are in kubeconfig preferred namespace (default)
	kubectl.RunWithOpts([]string{"get", "configmap", "my-cm", "-n", "default", "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})

	// Assert that kapp metaconfigmap is in kubeconfig preferred namespace (default)
	kubectl.RunWithOpts([]string{"get", "configmap", name + ".app", "-n", "default", "--kubeconfig", kubeconfigFile.Name()}, e2e.RunOpts{NoNamespace: true})
}

func Test_PackageInstall_DefaultNamespace(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "app-default-namespace"
	defaultNamespace := "default-namespace"
	pkgiNamespace := "pkgi-namespace"
	pkgiName := "test-pkgi"
	appName := "test-app"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	sas := e2e.ServiceAccounts{env.Namespace}

	namespacesYAML := fmt.Sprintf(`
---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
`, defaultNamespace, pkgiNamespace)

	rbacYAML := sas.ForDefaultNamespaceYAML(defaultNamespace, pkgiNamespace)

	installPkgYAML := fmt.Sprintf(`
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg.test.carvel.dev
  namespace: %[1]s
spec:
  # This is the name we want to reference in resources such as PackageInstall.
  displayName: "Test PackageMetadata in repo"
  shortDescription: "PackageMetadata used for testing"
  longDescription: "A longer, more detailed description of what the package contains and what it is for"
  providerName: Carvel
  maintainers:
  - name: carvel
  categories:
  - testing
  supportDescription: "Description of support provided for the package"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  licenses:
  - Apache 2.0
  capactiyRequirementsDescription: "cpu: 1,RAM: 2, Disk: 3"
  releaseNotes: |
    - Introduce simple-app package
  releasedAt: 2021-05-05T18:57:06Z
  template:
    spec:
      fetch:
      - inline:
          paths:
            file.yml: |
              apiVersion: v1
              kind: ConfigMap
              metadata:
                name: my-cm
              data:
                key: value
      template:
      - ytt: {}
      deploy:
      - kapp:
          rawOptions: ["--app-changes-max-to-keep=0"]
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
spec:
  defaultNamespace: %[3]s
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0`, pkgiNamespace, pkgiName, defaultNamespace)

	// deploy app with workspace and pkgi
	_, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(appWithInlineResources(appName, namespacesYAML+rbacYAML+installPkgYAML) + sas.ForClusterYAML()), AllowError: true})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// assert that pkgi resources are present in defaultNamespace
	kubectl.RunWithOpts([]string{"get", "configmap", "my-cm", "-n", defaultNamespace}, e2e.RunOpts{NoNamespace: true})

	// assert that kapp configmap is present in the pkgi namespace
	kubectl.RunWithOpts([]string{"get", "configmap", pkgiName + ".app", "-n", pkgiNamespace}, e2e.RunOpts{NoNamespace: true})

	// deploy app with workspace only, i.e. delete pkgi
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true, StdinReader: strings.NewReader(appWithInlineResources(appName, namespacesYAML) + sas.ForClusterYAML())})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	cleanUp()

	// deploy app with workspace and pkgi again
	_, err = kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{AllowError: true, StdinReader: strings.NewReader(appWithInlineResources(appName, namespacesYAML+rbacYAML+installPkgYAML) + sas.ForClusterYAML())})
	assert.NoError(t, err, "Expected app deploy to succeed, it did not")

	// delete app, i.e complete workspace along with pkgis
	_, err = kapp.RunWithOpts([]string{"delete", "-a", name}, e2e.RunOpts{AllowError: true})
	assert.NoError(t, err, "Expected app delete to succeed gracefully, it did not")
}

func appWithInlineResources(name, resources string) string {
	// Add indentation for the resources as these would be added inline to the App CR
	indentedResourceYAML := ""
	for _, line := range strings.Split(resources, "\n") {
		if line != "" {
			indentedResourceYAML += "          " + line + "\n"
		}
	}
	return fmt.Sprintf(`---
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
%s
  template:
  - ytt: {}
  deploy:
  - kapp: {}`, name, indentedResourceYAML)
}
