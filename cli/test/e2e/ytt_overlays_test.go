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
	cmName := "kctrl-overlay-test"

	yaml := fmt.Sprintf(`---
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
      - inline:
          paths:
            file.yml: |
              apiVersion: v1
              kind: ConfigMap
              metadata:
                name: %s
              data:
                key: value
      template:
      - ytt:
          paths:
          - "."
      deploy:
      - kapp: {}`, cmName)

	overlay1 := `
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind": "ConfigMap"})
---
metadata:
  #@overlay/match missing_ok=True
  annotations:
    #@overlay/match missing_ok=True
    test: foo
`
	overlay2 := `
#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.subset({"kind": "ConfigMap"})
---
metadata:
  #@overlay/match missing_ok=True
  annotations:
    #@overlay/match missing_ok=True
    test: bar
`

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
		kappCtrl.Run([]string{"package", "installed", "delete", "-i", pkgiName, "--yes"})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("add test package", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml),
		})
	})

	logger.Section("install with ytt overlays", func() {
		kappCtrl.RunWithOpts([]string{"package", "install", "-i", pkgiName, "--package", packageName, "--version", packageVersion, "--ytt-overlay-file", "-"}, RunOpts{
			StdinReader: strings.NewReader(overlay1),
		})

		out := kubectl.Run([]string{"get", "cm", cmName, "-o", "yaml"})
		require.Contains(t, out, "test: foo")
	})

	logger.Section("get package install", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))
		require.Exactly(t, fmt.Sprintf("- %s-%s-overlays", pkgiName, env.Namespace), output.Tables[0].Rows[0]["overlay_secrets"])
	})

	logger.Section("update ytt overlay", func() {
		kappCtrl.RunWithOpts([]string{"package", "installed", "update", "-i", pkgiName, "--ytt-overlay-file", "-"}, RunOpts{
			StdinReader: strings.NewReader(overlay2),
		})

		out := kubectl.Run([]string{"get", "cm", cmName, "-o", "yaml"})
		require.Contains(t, out, "test: bar")
	})

	logger.Section("drop ytt overlays", func() {
		kappCtrl.Run([]string{"package", "installed", "update", "-i", pkgiName, "--ytt-overlays=false"})

		out := kubectl.Run([]string{"get", "cm", cmName, "-o", "yaml"})
		require.NotContains(t, out, "test: foo")
		require.NotContains(t, out, "test: bar")
	})

	logger.Section("get package install", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "installed", "get", "--package-install", pkgiName, "--json"}, RunOpts{})
		require.NoError(t, err)

		output := uitest.JSONUIFromBytes(t, []byte(out))
		_, found := output.Tables[0].Rows[0]["overlay_secrets"]

		require.Exactly(t, false, found)
	})
}

func TestOverlaySecretCleanup(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	pkgiName := "overlay-test"
	appName := "test-secret"

	secretYaml := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
  name: %s-%s-overlays
  annotations:
    packaging.carvel.dev/package: %s-%s
type: Opaque
data:
  username: YWRtaW4=
`, pkgiName, env.Namespace, pkgiName, env.Namespace)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("testing overlay secret cleanup", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(secretYaml)})

		out := kappCtrl.Run([]string{"p", "i", "delete", "-i", pkgiName})
		require.Contains(t, out, fmt.Sprintf("Deleting 'secrets': %s-%s-overlays", pkgiName, env.Namespace))
	})
}
