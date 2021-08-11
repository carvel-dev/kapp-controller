// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func Test_PackageMetadataIsValidated(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	appName := "invalid-pkg-name-test"

	invalidPackageMetadataName := "I am invalid"

	invalidPkgYML := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: %s
spec:
  shortDescription: I am invalid
`, invalidPackageMetadataName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	defer cleanUp()

	logger.Section("deploy package metadata", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName},
			e2e.RunOpts{StdinReader: strings.NewReader(invalidPkgYML), AllowError: true})

		if err == nil {
			t.Fatalf("Expected package metadata creation to fail")
		}

		if !strings.Contains(err.Error(), "is invalid: metadata.name") {
			t.Fatalf("Expected package metadata creation error to contain message about invalid name, got: %v", err)
		}
	})
}

func TestOverridePackageMetadataDelete(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	k := e2e.Kubectl{t, env.Namespace, logger}

	localNS := env.Namespace
	globalNS := env.PackagingGlobalNS

	packageMetadataYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg1.test.carvel.dev
  namespace: %s
spec:
  displayName: "Global Package"
  shortDescription: "Package which is globally available"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg1.test.carvel.dev
  namespace: %s
spec:
  displayName: "Override Package"
  shortDescription: "Package which overrides global package"`, globalNS, localNS)

	cleanup := func() {
		k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", globalNS}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", localNS}, e2e.RunOpts{NoNamespace: true, AllowError: true})
	}
	defer cleanup()

	logger.Section("cleanup", cleanup)

	logger.Section("deploy package metadata", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(packageMetadataYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected package metadata application to succeed, but got: %v", err)
		}
	})

	logger.Section("attempt to delete the local package metadata", func() {
		timeout := 30 * time.Second

		ctx, _ := context.WithTimeout(context.Background(), timeout)
		_, err := k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", localNS}, e2e.RunOpts{Ctx: ctx, NoNamespace: true, AllowError: true})
		if err != nil {
			t.Fatalf("Expected delete of local package metadata to succeed in %v, but got: %v", timeout, err)
		}
	})
}

func TestOverridePackageMetadataNamespaceDelete(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	k := e2e.Kubectl{t, env.Namespace, logger}

	localNS := "test-ns"
	globalNS := env.PackagingGlobalNS

	packageMetadataYaml := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %[1]s
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg1.test.carvel.dev
  namespace: %[1]s
spec:
  displayName: "Override Package"
  shortDescription: "Package which overrides global package"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg1.test.carvel.dev
  namespace: %[2]s
spec:
  displayName: "Global Package"
  shortDescription: "Package which is globally available"`, localNS, globalNS)

	cleanup := func() {
		k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", globalNS}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", localNS}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", fmt.Sprintf("namespaces/%s", localNS)}, e2e.RunOpts{NoNamespace: true, AllowError: true})
	}
	defer logger.Section("post test cleanup", cleanup)

	logger.Section("pre test cleanup", cleanup)

	logger.Section("deploy package metadata and namespace", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(packageMetadataYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected package metadata application to succeed, but got: %v", err)
		}
	})

	logger.Section("attempt to delete the local namespace", func() {
		timeout := 30 * time.Second

		ctx, _ := context.WithTimeout(context.Background(), timeout)
		_, err := k.RunWithOpts([]string{"delete", fmt.Sprintf("namespaces/%s", localNS)}, e2e.RunOpts{Ctx: ctx, NoNamespace: true, AllowError: true})
		if err != nil {
			if strings.Contains(err.Error(), "signal: interrupt") {
				t.Fatalf("Timed out waiting for delete of namespace '%s'", localNS)
			}
			t.Fatalf("Expected delete of local namespace '%s' to succeed, but got: %v", localNS, err)
		}

		_, err = k.RunWithOpts([]string{"get", fmt.Sprintf("namespaces/%s", localNS)}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		if err == nil {
			t.Fatalf("Expected not to find local namespace '%s', but did", localNS)
		}
	})
}

func TestOverridePackageMetadataCreate(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	k := e2e.Kubectl{t, env.Namespace, logger}

	localNS := "test-ns"
	globalNS := env.PackagingGlobalNS

	packageMetadataYaml := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %[1]s
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg1.test.carvel.dev
  namespace: %[2]s
spec:
  displayName: "Global Package"
  shortDescription: "Package which is globally available"`, localNS, globalNS)

	updatePackageMetadataYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg1.test.carvel.dev
  namespace: %[1]s
spec:
  displayName: "Override Package"
  shortDescription: "Package which overrides global package"`, localNS)

	cleanup := func() {
		k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", globalNS}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", "pkgm/pkg1.test.carvel.dev", "-n", localNS}, e2e.RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", fmt.Sprintf("namespaces/%s", localNS)}, e2e.RunOpts{NoNamespace: true, AllowError: true})
	}
	defer logger.Section("post test cleanup", cleanup)

	logger.Section("pre test cleanup", cleanup)

	logger.Section("deploy package metadata and namespace", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(packageMetadataYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected package metadata application to succeed, but got: %v", err)
		}
	})

	logger.Section("attempt to create an override package metadata", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(updatePackageMetadataYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected creation of local package metadata to succeed, but got: %v", err)
		}
	})
}
