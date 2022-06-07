// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageInstallValues(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	appName := "test-package-name"
	pkgiName := "testpkgi"
	packageMetadataName := "test-pkg.carvel.dev"
	packageVersion := "1.0.0"

	packageMetadata := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: %s
spec:
  displayName: "Carvel Test Package"
  shortDescription: "Carvel package for testing installation"`, packageMetadataName)

	packageCR := `---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: test-pkg.carvel.dev.1.0.0
spec:
  refName: test-pkg.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt:
          paths:
          - config/
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}`

	valuesFile := `
---
foo: bar
`

	yaml := packageMetadata + "\n" + packageCR

	cleanUp := func() {
		// TODO: Check for error while uninstalling in cleanup?
		kappCtrl.Run([]string{"package", "installed", "delete", "--package-install", pkgiName})
		kapp.Run([]string{"delete", "-a", appName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("Adding test package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml), AllowError: true,
		})
		require.NoError(t, err)
	})

	logger.Section("Installing test package with values config", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "create", "--package-install", pkgiName, "-p", packageMetadataName,
			"--version", packageVersion, "--values-file", "-"}, RunOpts{StdinReader: strings.NewReader(valuesFile)})
		require.NoError(t, err)

		// Check for owned value secret
		secretName := fmt.Sprintf("%s-%s-values", pkgiName, env.Namespace)
		out, err := kubectl.RunWithOpts([]string{"get", "secret", secretName, "-o", "yaml"}, RunOpts{})
		require.NoError(t, err)
		require.Contains(t, out, pkgiName+"-"+"kctrl-test")

		// Check for reference to values secret
		out, err = kubectl.RunWithOpts([]string{"get", "pkgi", pkgiName, "-o", "yaml"}, RunOpts{})
		require.NoError(t, err)
		require.Contains(t, out, secretName)
	})

	// TODO: Add check for ensuring that we wait for reconciliation when secrets are updated
	// When https://github.com/vmware-tanzu/carvel-kapp-controller/issues/670 is resolved

	logger.Section("Updating values config for test package", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "update", "--package-install", pkgiName, "--values-file", "-"}, RunOpts{StdinReader: strings.NewReader(valuesFile)})
		require.NoError(t, err)
		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "Deploy succeeded")

		// Check that ownership annotations are intact
		secretName := fmt.Sprintf("%s-%s-values", pkgiName, env.Namespace)
		out, err = kubectl.RunWithOpts([]string{"get", "secret", secretName, "-o", "yaml"}, RunOpts{})
		require.NoError(t, err)
		require.Contains(t, out, pkgiName+"-"+"kctrl-test")
	})

	logger.Section("Dropping consumed values file", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "update", "--package-install", pkgiName, "--values=false"}, RunOpts{})
		require.NoError(t, err)

		// Check for owned value secret
		secretName := fmt.Sprintf("%s-%s-values", pkgiName, env.Namespace)
		_, err = kubectl.RunWithOpts([]string{"get", "secret", secretName, "-o", "yaml"}, RunOpts{AllowError: true})
		require.Error(t, err)
		require.Contains(t, err.Error(), "NotFound")

		// Check for reference to values secret
		out, err := kubectl.RunWithOpts([]string{"get", "pkgi", pkgiName, "-o", "yaml"}, RunOpts{})
		require.NoError(t, err)
		require.NotContains(t, out, secretName)
	})

	// Check for successful deletion post dropping values file. kctrl should not try to delete secret
	logger.Section("package install delete", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "delete", "--package-install", pkgiName}, RunOpts{})
		require.NoError(t, err)
	})
}
