// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package secretgencontroller

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func Test_PlaceholderSecrets_DeletedWhenPackageInstallDeleted(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "placeholder-garbage-collection"
	sas := e2e.ServiceAccounts{env.Namespace}

	pkgiYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: "package"
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt: {}
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting package"
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Create PackageInstall", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgiYaml)})

		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "1m"})
	})

	logger.Section("Check placeholder secret created", func() {
		kubectl.Run([]string{"get", "secret", name + "-fetch-0"})
	})

	logger.Section("Check placeholder secret deleted after PackageInstall deleted", func() {
		cleanUp()

		// Sleep to give PackageInstall time to delete
		time.Sleep(1 * time.Second)
		out, err := kubectl.RunWithOpts([]string{"get", "secret", name + "-fetch-0"}, e2e.RunOpts{AllowError: true})

		require.NotNil(t, err, "expected error from not finding placeholder secret.\nGot: "+out)
		assert.True(t, strings.Contains("NotFound", out), "expected error to be is not found but got: %s", err)
	})
}

func Test_PackageInstallAndRepo_CanAuthenticateToPrivateRepository_UsingPlaceholderSecret(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	// If this changes, the skip-tls-verify domain must be updated to match
	pkgiName := "placeholder-private-auth-pkgi"
	pkgrName := "placeholder-private-auth-pkgr"
	registryNamespace := "registry"
	registryName := "test-registry"
	configName := "test-registry-ca-cert-config"

	secretYaml := fmt.Sprintf(`
---
apiVersion: v1
kind: Secret
metadata:
  name: regcred
type: kubernetes.io/dockerconfigjson
stringData:
  .dockerconfigjson: |
    {
      "auths": {
        "registry-svc.%s.svc.cluster.local:443": {
          "username": "testuser",
          "password": "testpassword",
          "auth": ""
        }
      }
    }
---
apiVersion: secretgen.carvel.dev/v1alpha1
kind: SecretExport
metadata:
  name: regcred
spec:
  toNamespaces:
  - %s
`, registryNamespace, env.Namespace)

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + pkgiName}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", pkgiName})
		kapp.Run([]string{"delete", "-a", pkgrName})
		kapp.Run([]string{"delete", "-a", configName})
		kapp.Run([]string{"delete", "-a", "secret-export"})
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy controller config to skip registry TLS verify", func() {
		config := `
apiVersion: v1
kind: Secret
metadata:
  name: kapp-controller-config
  namespace: kapp-controller
stringData:
  dangerousSkipTLSVerify: registry-svc.registry.svc.cluster.local
`
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", configName},
			e2e.RunOpts{StdinReader: strings.NewReader(config)})

		// Since config propagation is async, just wait a little bit
		time.Sleep(2 * time.Second)
	})

	logger.Section("deploy registry with self signed certs", func() {
		kapp.Run([]string{
			"deploy", "-a", registryName,
			"-f", "../assets/registry/registry2.yml",
			"-f", "../assets/registry/certs-for-skip-tls.yml",
			"-f", "../assets/registry/htpasswd-auth",
			"-f", "../assets/registry/registry-contents.yml",
		})
	})

	logger.Section("create exported registry secret", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", "secret-export", "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(secretYaml)})
	})

	logger.Section("deploy PackageInstall that auths to registry", func() {
		pkgiYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: "package"
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: registry-svc.%[3]s.svc.cluster.local:443/secret-test/test-bundle
      template:
      - ytt: {}
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
    kapp.k14s.io/change-rule: "upsert after upserting package"
spec:
  syncPeriod: 30s
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, pkgiName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgiName, "-f", "-"}, e2e.RunOpts{
			StdinReader: strings.NewReader(pkgiYaml),
			OnErrKubectl: []string{"get", "app", "placeholder-private-auth", "-oyaml"},
		})

		kubectl.Run([]string{"get", "configmap", "e2e-test-map"})

		// Clean up befor trying to deploy package repository
		kapp.Run([]string{"delete", "-a", pkgiName})
	})

	logger.Section("deploy PackageRepository that auths to registry", func() {
		pkgrYaml := fmt.Sprintf(`---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[2]s
  namespace: %[1]s
spec:
  fetch:
    image:
      url: registry-svc.%[3]s.svc.cluster.local:443/secret-test/test-repo
`, env.Namespace, pkgrName, registryNamespace) + sas.ForNamespaceYAML()

		kapp.RunWithOpts([]string{"deploy", "-a", pkgrName, "-f", "-"}, e2e.RunOpts{
			StdinReader: strings.NewReader(pkgrYaml),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})

		kubectl.Run([]string{"get", "package", "pkg.test.carvel.dev.1.0.0"})
	})
}
