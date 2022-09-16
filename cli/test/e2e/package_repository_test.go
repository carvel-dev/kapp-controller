// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"testing"

	uitest "github.com/cppforlife/go-cli-ui/ui/test"
	"github.com/stretchr/testify/require"
)

func TestPackageRepository(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	pkgrName := "test-package-repository"
	pkgrURL := `index.docker.io/k8slt/kc-e2e-test-repo:latest`

	kind := "PackageRepository"

	cleanUp := func() {
		RemoveClusterResource(t, kind, pkgrName, env.Namespace, kubectl)
	}

	cleanUp()
	defer cleanUp()

	logger.Section("listing repositories with no repository present", func() {
		out, _ := kappCtrl.RunWithOpts([]string{"package", "repository", "list", "--json"}, RunOpts{})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("getting a nonexistent repository", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "repository", "get", "-r", pkgrName, "--json"}, RunOpts{
			AllowError: true,
		})
		require.Error(t, err)
	})

	logger.Section("deletion of nonexistent repository", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "repository", "delete", "-r", pkgrName}, RunOpts{
			AllowError: true,
		})
		require.Error(t, err)
	})

	logger.Section("adding a repository", func() {
		kappCtrl.Run([]string{"package", "repository", "add", "-r", pkgrName, "--url", pkgrURL})

		kubectl.Run([]string{"get", kind, pkgrName})
		kubectl.Run([]string{"get", "pkgm/pkg.test.carvel.dev"})
		kubectl.Run([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"})
		kubectl.Run([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"})
	})

	logger.Section("kicking a repository", func() {
		out := kappCtrl.Run([]string{"package", "repository", "kick", "-r", pkgrName})

		require.Contains(t, out, "Deploy succeeded")
	})

	logger.Section("listing repositories with one repository being present", func() {
		out, _ := kappCtrl.RunWithOpts([]string{"package", "repository", "list", "--json"}, RunOpts{})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"name":   pkgrName,
			"source": fmt.Sprintf("(imgpkg) %s", pkgrURL),
			"status": "Reconcile succeeded",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("getting a repository", func() {
		out, _ := kappCtrl.RunWithOpts([]string{"package", "repository", "get", "-r", pkgrName, "--json"}, RunOpts{})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":           "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"status":               "Reconcile succeeded",
			"namespace":            env.Namespace,
			"name":                 pkgrName,
			"source":               fmt.Sprintf("(imgpkg) %s", pkgrURL),
			"useful_error_message": "",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

	logger.Section("deletion of an existing repository", func() {
		kappCtrl.RunWithOpts([]string{"package", "repository", "delete", "-r", pkgrName}, RunOpts{})

		_, err := kubectl.RunWithOpts([]string{"get", "packagerepository", pkgrName}, RunOpts{AllowError: true})
		require.Contains(t, err.Error(), "not found")
	})

	logger.Section("updating a repository", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "repository", "add", "-r", pkgrName, "--url", "https://carvel.dev"}, RunOpts{
			AllowError: true})
		require.Error(t, err)

		kubectl.Run([]string{"get", kind, pkgrName})

		kappCtrl.Run([]string{"package", "repository", "update", "-r", pkgrName, "--url", pkgrURL})

		out := kappCtrl.Run([]string{"package", "repository", "get", "-r", pkgrName, "--json"})

		output := uitest.JSONUIFromBytes(t, []byte(out))

		expectedOutputRows := []map[string]string{{
			"conditions":           "- type: ReconcileSucceeded\n  status: \"True\"\n  reason: \"\"\n  message: \"\"",
			"status":               "Reconcile succeeded",
			"namespace":            env.Namespace,
			"name":                 pkgrName,
			"source":               fmt.Sprintf("(imgpkg) %s", pkgrURL),
			"useful_error_message": "",
		}}
		require.Exactly(t, expectedOutputRows, output.Tables[0].Rows)
	})

}
