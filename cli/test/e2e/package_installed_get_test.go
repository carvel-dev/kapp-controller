// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageInstalledGet(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kapp{t, env.Namespace, env.KappCtrlBinaryPath, logger}

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

	packageName := "test-pkg.carvel.dev.1.0.0"
	packageVersion := "1.0.0"

	packageCR := fmt.Sprintf(`---
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
      - kapp: {}`, packageName, packageVersion)

	yaml := packageMetadata + "\n" + packageCR

	cleanUp := func() {
		kappCtrl.Run([]string{"package", "installed", "delete", "--package-install", pkgiName})
		kapp.Run([]string{"delete", "-a", appName})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("package installed get", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml),
		})

		kappCtrl.RunWithOpts([]string{
			"package", "installed", "create",
			"--package-install", pkgiName,
			"--package-name", packageMetadataName,
			"--version", packageVersion,
		}, RunOpts{})

		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":           "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"namespace":            env.Namespace,
			"name":                 "testpkgi",
			"package_name":         "test-pkg.carvel.dev",
			"package_version":      "1.0.0",
			"description":          "Reconcile succeeded",
			"useful_error_message": "",
		}}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
