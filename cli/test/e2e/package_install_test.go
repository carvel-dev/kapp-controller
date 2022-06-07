// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageInstalled(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	appName := "test-package-name"
	pkgiName := "testpkgi"
	packageMetadataName := "test-pkg.carvel.dev"

	packageMetadata := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: %s
spec:
  displayName: "Carvel Test Package"
  shortDescription: "Carvel package for testing installation"`, packageMetadataName)

	packageName1 := "test-pkg.carvel.dev.1.0.0"
	packageName2 := "test-pkg.carvel.dev.2.0.0"
	packageVersion1 := "1.0.0"
	packageVersion2 := "2.0.0"

	packageCR := `---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: %s
spec:
  refName: test-pkg.carvel.dev
  version: %s
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

	packageCR1 := fmt.Sprintf(packageCR, packageName1, packageVersion1)

	packageCR2 := fmt.Sprintf(packageCR, packageName2, packageVersion2)

	yaml := packageMetadata + "\n" + packageCR1 + "\n" + packageCR2

	cleanUp := func() {
		// TODO: Check for error while uninstalling in cleanup?
		kappCtrl.Run([]string{"package", "installed", "delete", "--package-install", pkgiName})
		kapp.Run([]string{"delete", "-a", appName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("package installed list with no package present", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("Adding test package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml), AllowError: true,
		})
		require.NoError(t, err)
	})

	logger.Section("Installing test package", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "create", "--package-install", pkgiName, "-p", packageMetadataName, "--version", packageVersion1}, RunOpts{})
		require.NoError(t, err)
	})

	logger.Section("Check for created resources", func() {
		svcAccountName := fmt.Sprintf("%s-%s-sa", pkgiName, env.Namespace)
		_, err := kubectl.RunWithOpts([]string{"get", "sa", svcAccountName}, RunOpts{})
		require.NoError(t, err)

		clusterRoleName := fmt.Sprintf("%s-%s-cluster-role", pkgiName, env.Namespace)
		_, err = kubectl.RunWithOpts([]string{"get", "clusterroles", clusterRoleName}, RunOpts{})
		require.NoError(t, err)

		clusterRoleBindingName := fmt.Sprintf("%s-%s-cluster-rolebinding", pkgiName, env.Namespace)
		_, err = kubectl.RunWithOpts([]string{"get", "clusterrolebindings", clusterRoleBindingName}, RunOpts{})
		require.NoError(t, err)
	})

	logger.Section("package installed list with one package installed", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"name":            "testpkgi",
			"package_name":    "test-pkg.carvel.dev",
			"package_version": "1.0.0",
			"status":          "Reconcile succeeded",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("package installed get", func() {

		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":      "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"namespace":       env.Namespace,
			"name":            "testpkgi",
			"package_name":    "test-pkg.carvel.dev",
			"package_version": "1.0.0",
			"status":          "Reconcile succeeded",
		}}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("package installed status", func() {
		out := kappCtrl.Run([]string{"package", "installed", "status", "-i", pkgiName})

		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "Deploy succeeded")
	})

	logger.Section("package installed update", func() {

		_, err := kappCtrl.RunWithOpts([]string{
			"package", "installed", "update",
			"--package-install", pkgiName,
			"-p", packageMetadataName,
			"--version", packageVersion2,
		}, RunOpts{})
		require.NoError(t, err)

		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":      "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"namespace":       env.Namespace,
			"name":            "testpkgi",
			"package_name":    "test-pkg.carvel.dev",
			"package_version": "2.0.0",
			"status":          "Reconcile succeeded",
		}}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("package installed pause", func() {

		_, err := kappCtrl.RunWithOpts([]string{
			"package", "installed", "pause",
			"--package-install", pkgiName,
		}, RunOpts{})
		require.NoError(t, err)

		// Sleep for 1 second as it takes time for the underlying app cr to be paused
		time.Sleep(1 * time.Second)

		out, err := kubectl.RunWithOpts([]string{"get", "app", pkgiName}, RunOpts{})
		require.NoError(t, err)

		require.Contains(t, out, "Canceled/paused")
	})

	logger.Section("package installed kick", func() {

		out, err := kappCtrl.RunWithOpts([]string{
			"package", "installed", "kick",
			"--package-install", pkgiName,
		}, RunOpts{})
		require.NoError(t, err)

		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "Deploy succeeded")
	})

	logger.Section("package install delete", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "installed", "delete", "--package-install", pkgiName}, RunOpts{})
		require.NoError(t, err)

		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "list", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("Check for deletion of created resources", func() {
		svcAccountName := fmt.Sprintf("%s-%s-sa", pkgiName, env.Namespace)
		_, err := kubectl.RunWithOpts([]string{"get", "sa", svcAccountName}, RunOpts{
			AllowError: true,
		})
		require.Contains(t, err.Error(), "not found")

		clusterRoleName := fmt.Sprintf("%s-%s-cluster-role", pkgiName, env.Namespace)
		_, err = kubectl.RunWithOpts([]string{"get", "clusterroles", clusterRoleName}, RunOpts{
			AllowError: true,
		})
		require.Contains(t, err.Error(), "not found")

		clusterRoleBindingName := fmt.Sprintf("%s-%s-cluster-rolebinding", pkgiName, env.Namespace)
		_, err = kubectl.RunWithOpts([]string{"get", "clusterrolebindings", clusterRoleBindingName}, RunOpts{
			AllowError: true,
		})
		require.Contains(t, err.Error(), "not found")
	})
}
