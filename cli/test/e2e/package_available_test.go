// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"os"
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
        required_property_without_default:
          type: integer
          description: This is a required property without a default value
      required:
      - required_property_without_default
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

	logger.Section("listing packages with column names", func() {
		out := kappCtrl.Run([]string{"package", "available", "list", "--json", "--column=name,namespace"})
		output := uitest.JSONUIFromBytes(t, []byte(out))
		expectedOutputRows := []map[string]string{{
			"name":      packageName,
			"namespace": env.Namespace,
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("listing packages with non existing column names", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "available", "list", "--json", "--column=name,invalid,namespace,ns"}, RunOpts{
			AllowError: true,
		})
		expectedError := "kctrl: Error: invalid column names: invalid,ns"
		require.ErrorContains(t, err, expectedError)
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

	logger.Section("getting a package with invalid column names", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "available", "get", "-p", packageName, "--column=name,invalid,namespace,ns"}, RunOpts{
			AllowError: true,
		})
		expectedError := "kctrl: Error: invalid column names: invalid,namespace,ns"
		require.ErrorContains(t, err, expectedError)
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
			{
				"default":     "",
				"description": "This is a required property without a default value",
				"key":         "required_property_without_default",
				"type":        "integer",
			},
		}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("getting default-values-file-output of a package", func() {
		tmpDir, err := os.MkdirTemp("", "*")
		require.NoError(t, err, "Error creating temp directory")
		defer os.RemoveAll(tmpDir)

		kappCtrl.Run([]string{"package", "available", "get", "-p", fmt.Sprintf("%s/%s", packageName, "1.0.0"), "--default-values-file-output", tmpDir + "/default-values.yaml", "--json"})

		// properties with no default value should be included with just the key
		expectedDefaultValuesFileOutout := `# app_port: 80
# required_property_without_default:
`
		out, err := os.ReadFile(tmpDir + "/default-values.yaml")
		require.NoError(t, err, "Error reading default values file output")

		require.Equal(t, expectedDefaultValuesFileOutout, string(out), "default values file output does not match")
	})
}

func TestPackageAvailableGet_WithEmptyValuesSchema(t *testing.T) {
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
      title: Empty values schema
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

	logger.Section("getting value schema of a package", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml),
		})

		_, err := kappCtrl.RunWithOpts([]string{"package", "available", "get", "-p", fmt.Sprintf("%s/%s", packageName, "1.0.0"), "--values-schema", "--json"},
			RunOpts{AllowError: true})

		require.Error(t, err)
		require.ErrorContains(t, err, "hint: the valuesSchema might not have any properties")
	})

	logger.Section("getting default-values-file-output of a package", func() {
		out := kappCtrl.Run([]string{"package", "available", "get", "-p", fmt.Sprintf("%s/%s", packageName, "1.0.0"), "--default-values-file-output", "default-values.yaml", "--json"})
		require.Contains(t, out, "does not have any user configurable values")
	})
}
