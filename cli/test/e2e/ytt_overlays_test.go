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

func TestYttOverlays(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	packageName := "test-pkg.carvel.dev"
	packageVersion := "1.0.0"
	appName := "test-package"
	pkgiName := "overlay-test"
	deploymentName := "simple-app"

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
      - kapp: {}`

	overlay1 := `
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind": "Deployment"})
---
metadata:
  #@overlay/match missing_ok=True
  annotations:
    #@overlay/match missing_ok=True
    test: foo
`
	overlay2 := `
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind": "Deployment"})
---
metadata:
  #@overlay/match missing_ok=True
  annotations:
    #@overlay/match missing_ok=True
    test: bar
`

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("add test package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml),
		})
		require.NoError(t, err)
	})

	logger.Section("install with ytt overlays", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "install", "-i", pkgiName, "--package", packageName, "--version", packageVersion, "--ytt-overlay-file", "-"}, RunOpts{
			StdinReader: strings.NewReader(overlay1),
		})
		require.NoError(t, err)
		//Ensure that progress is tailed
		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "Deploy succeeded")

		out = kubectl.Run([]string{"get", "deployment", deploymentName, "-o", "yaml"})
		require.Contains(t, out, "test: foo")
	})

	logger.Section("get package install", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":      "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"namespace":       env.Namespace,
			"name":            pkgiName,
			"overlay_secrets": fmt.Sprintf("- %s-%s-overlays", pkgiName, env.Namespace),
			"package_name":    "test-pkg.carvel.dev",
			"package_version": "1.0.0",
			"status":          "Reconcile succeeded",
		}}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("update ytt overlay", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "update", "-i", pkgiName, "--ytt-overlay-file", "-"}, RunOpts{
			StdinReader: strings.NewReader(overlay2),
		})
		require.NoError(t, err)
		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "Deploy succeeded")

		out = kubectl.Run([]string{"get", "deployment", deploymentName, "-o", "yaml"})
		require.Contains(t, out, "test: bar")
	})

	logger.Section("drop ytt overlays", func() {
		out := kappCtrl.Run([]string{"package", "installed", "update", "-i", pkgiName, "--ytt-overlays=false"})
		require.Contains(t, out, "Fetch succeeded")
		require.Contains(t, out, "Template succeeded")
		require.Contains(t, out, "Deploy succeeded")

		out = kubectl.Run([]string{"get", "deployment", deploymentName, "-o", "yaml"})
		require.NotContains(t, out, "test: foo")
		require.NotContains(t, out, "test: bar")
	})

	logger.Section("get package install", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":      "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"namespace":       env.Namespace,
			"name":            pkgiName,
			"package_name":    "test-pkg.carvel.dev",
			"package_version": "1.0.0",
			"status":          "Reconcile succeeded",
		}}

		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})
}
