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

func TestPackageAvailable(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	packageName := "test-pkg.carvel.dev"
	appName := "test-package"

	yaml := `---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: test-pkg.carvel.dev
spec:
  displayName: "Carvel Test Package"
  shortDescription: "Carvel package for testing installation"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: test-pkg.carvel.dev.1.0.0
spec:
  refName: test-pkg.carvel.dev
  version: 1.0.0
  valuesSchema:
    openAPIv3:
      properties:
        app_port:
          default: 80
          description: App port
          type: integer
        app_name:
          description: App Name
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
      - kapp: {}
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: test-pkg.carvel.dev.1.1.0
spec:
  refName: test-pkg.carvel.dev
  version: 1.1.0
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

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("package available list with no package present", func() {
		out := kappCtrl.Run([]string{"package", "available", "list", "--json", "--wide"})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("getting a nonexisting package", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "available", "get", "-p", packageName, "--json"}, RunOpts{
			AllowError: true,
		})
		require.Error(t, err, "Expected to get an error")
	})

	logger.Section("listing packages", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml),
		})

		out := kappCtrl.Run([]string{"package", "available", "list", "--json", "--wide"})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"name":              packageName,
			"display_name":      "Carvel Test Package",
			"short_description": "Carvel package for testing installation",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("listing versions of a package", func() {
		out := kappCtrl.Run([]string{"package", "available", "list", "-p", packageName, "--json"})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"name":        packageName,
				"version":     "1.0.0",
				"released_at": "-",
			},
			{
				"name":        packageName,
				"version":     "1.1.0",
				"released_at": "-",
			},
		}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("getting a package", func() {
		out := kappCtrl.Run([]string{"package", "available", "get", "-p", packageName, "--json"})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"categories":          "",
			"display_name":        "Carvel Test Package",
			"long_description":    "",
			"maintainers":         "",
			"name":                packageName,
			"provider":            "",
			"short_description":   "Carvel package for testing installation",
			"support_description": "",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)

		expectedOutputRows = []map[string]string{
			{
				"version":     "1.0.0",
				"released_at": "-",
			},
			{
				"version":     "1.1.0",
				"released_at": "-",
			},
		}
		require.Exactly(t, expectedOutputRows, output.Tables[1].Rows)
	})

	logger.Section("getting value schema of a package", func() {
		out := kappCtrl.Run([]string{"package", "available", "get", "-p", fmt.Sprintf("%s/%s", packageName, "1.0.0"), "--values-schema", "--json"})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"default":     "",
				"description": "App Name",
				"key":         "app_name",
				"type":        "",
			},
			{
				"default":     "80",
				"description": "App port",
				"key":         "app_port",
				"type":        "integer",
			},
		}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
