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

func TestPackageAvailableGet(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kapp{t, env.Namespace, env.KappCtrlBinaryPath, logger}

	appName := "test-package-name"

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
  valuesSchema:
    openAPIv3:
      properties:
        app_port:
          default: 80
          description: App port
          type: integer
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

	logger.Section("package available get", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "get", "-p", packageMetadataName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"categories":          "",
			"display_name":        "Carvel Test Package",
			"long_description":    "",
			"maintainers":         "",
			"name":                "test-pkg.carvel.dev",
			"provider":            "",
			"short_description":   "Carvel package for testing installation",
			"support_description": "",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)

		expectedOutputRows = []map[string]string{{
			"version":     "1.0.0",
			"released_at": "0001-01-01 00:00:00 +0000 UTC",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[1].Rows)
	})

	logger.Section("package available get value schema", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "get", "-p", fmt.Sprintf("%s/%s", packageMetadataName, packageVersion), "--values-schema", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"default":     "80",
			"description": "App port",
			"key":         "app_port",
			"type":        "integer",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
