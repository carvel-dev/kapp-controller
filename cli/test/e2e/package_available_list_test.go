// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"strings"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageAvailableList(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	appName := "test-package-name"

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
	defer cleanUp()

	logger.Section("package available list with no package present", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list", "--json", "--wide"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("package available list with one package available", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml),
		})

		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list", "--json", "--wide"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"name":              "test-pkg.carvel.dev",
				"display_name":      "Carvel Test Package",
				"short_description": "Carvel package for testing installation",
			},
		}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("package available list versions of a package", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list", "-p", "test-pkg.carvel.dev", "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{
			{
				"name":        "test-pkg.carvel.dev",
				"version":     "1.0.0",
				"released_at": "0001-01-01 00:00:00 +0000 UTC",
			},
			{
				"name":        "test-pkg.carvel.dev",
				"version":     "1.1.0",
				"released_at": "0001-01-01 00:00:00 +0000 UTC",
			},
		}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
