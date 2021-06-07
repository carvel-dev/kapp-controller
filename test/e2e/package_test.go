// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test_PackageIsValidated(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	appName := "invalid-pkg-name-test"

	invalidPackageName := "I am invalid"

	invalidPkgYML := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: %s
spec:
  shortDescription: I am invalid
`, invalidPackageName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	defer cleanUp()

	logger.Section("deploy package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName},
			RunOpts{StdinReader: strings.NewReader(invalidPkgYML), AllowError: true})

		if err == nil {
			t.Fatalf("Expected package creation to fail")
		}

		if !strings.Contains(err.Error(), "is invalid: metadata.name") {
			t.Fatalf("Expected package creation error to contain message about invalid name, got: %v", err)
		}
	})
}

func TestOverridePackageDelete(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	k := Kubectl{t, env.Namespace, logger}

	localNS := env.Namespace
	globalNS := env.PackagingGlobalNS

	packagesYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg1.test.carvel.dev
  namespace: %s
spec:
  displayName: "Global Package"
  shortDescription: "Package which is globally available"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg1.test.carvel.dev
  namespace: %s
spec:
  displayName: "Override Package"
  shortDescription: "Package which overrides global package"`, globalNS, localNS)

	cleanup := func() {
		k.RunWithOpts([]string{"delete", "packages/pkg1.test.carvel.dev", "-n", globalNS}, RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", "packages/pkg1.test.carvel.dev", "-n", localNS}, RunOpts{NoNamespace: true, AllowError: true})
	}
	defer cleanup()

	logger.Section("cleanup", cleanup)

	logger.Section("deploy packages", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(packagesYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected package application to succeed, but got: %#v", err)
		}
	})

	logger.Section("attempt to delete the local package", func() {
		timeout := 30 * time.Second
		cancelCh := make(chan struct{})
		defer close(cancelCh)
		go func() {
			time.Sleep(timeout)
			cancelCh <- struct{}{}
		}()

		_, err := k.RunWithOpts([]string{"delete", "packages/pkg1.test.carvel.dev", "-n", localNS}, RunOpts{CancelCh: cancelCh, NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected delete of local package to succeed, but got: %#v", err)
		}
	})
}
