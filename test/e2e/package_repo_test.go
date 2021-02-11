// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TODO: Right now the implementation of the package repo reconciler needs improvement
// hopefully after which, these tests can be cleaned up to remove retries and time related
// falkeyness

func Test_PackageRepoBundle_PackagesAvailable(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kubectl := Kubectl{t, env.Namespace, logger}
	// contents of this bundle (k8slt/k8slt/kappctrl-e2e-repo-bundle)
	// under examples/packaging-demo/repo-bundle
	yamlRepo := `---
apiVersion: install.package.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
  # cluster scoped
spec:
  fetch:
    bundle:
      image: k8slt/kappctrl-e2e-repo-bundle`

	cleanUp := func() {
		kubectl.RunWithOpts([]string{"delete", "pkgr/basic.test.carvel.dev"}, RunOpts{NoNamespace: true})
	}
	defer cleanUp()

	kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(yamlRepo)})

	retryFunc := func() error {
		_, err := kubectl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.1.0.0"}, RunOpts{NoNamespace: true, AllowError: true})
		if err != nil {
			return err
		}
		_, err = kubectl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.2.0.0"}, RunOpts{NoNamespace: true, AllowError: true})
		if err != nil {
			return err
		}
		return nil
	}

	err := retry(10*time.Second, retryFunc)
	if err != nil {
		t.Fatalf("Expected to find pkgs (pkg2.test.carvel.dev.1.0.0, pkg2.test.carvel.dev.2.0.0) but couldn't: %v", err)
	}
}

func Test_PackageRepoDelete(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kctl := Kubectl{t, env.Namespace, logger}

	repoYaml := `---
apiVersion: install.package.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.delete.test.carvel.dev
spec:
  fetch:
    bundle:
      image: k8slt/kappctrl-e2e-repo-bundle`

	packageNames := []string{"pkg2.test.carvel.dev.1.0.0", "pkg2.test.carvel.dev.2.0.0"}

	cleanUp := func() {
		kctl.RunWithOpts([]string{"delete", "pkgr/basic.test.carvel.dev"}, RunOpts{NoNamespace: true, AllowError: true})
		for _, name := range packageNames {
			kctl.RunWithOpts([]string{"delete", fmt.Sprintf("package/%s", name)}, RunOpts{NoNamespace: true, AllowError: true})
		}
	}
	defer cleanUp()

	logger.Section("deploy repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", "repo"},
			RunOpts{StdinReader: strings.NewReader(repoYaml)})
	})

	logger.Section("check packages exist", func() {
		err := retry(20*time.Second, func() error {
			_, err := kctl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err != nil {
				return err
			}

			_, err = kctl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			t.Fatalf("Expected to find pkgs '%s' but couldn't: %v", strings.Join(packageNames, ", "), err)
		}
	})

	logger.Section("delete repo", func() {
		kapp.Run([]string{"delete", "-a", "repo"})
	})

	logger.Section("check packages are deleted too", func() {
		err := retry(10*time.Second, func() error {
			_, err := kctl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err == nil || !strings.Contains(err.Error(), "\"pkg2.test.carvel.dev.1.0.0\" not found") {
				return fmt.Errorf("found package pkg2.test.carvel.dev.1.0.0")
			}

			_, err = kctl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err == nil || !strings.Contains(err.Error(), "\"pkg2.test.carvel.dev.2.0.0\" not found") {
				return fmt.Errorf("found package pkg2.test.carvel.dev.2.0.0")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Expected not to find pkgs '%s' but %s", strings.Join(packageNames, ", "), err)
		}
	})
}

func retry(timeout time.Duration, f func() error) error {
	var err error
	stopTime := time.Now().Add(timeout)
	for {
		err = f()
		if err == nil {
			return nil
		}
		if time.Now().After(stopTime) {
			return fmt.Errorf("retry timed out after %s: %v", timeout.String(), err)
		}
		time.Sleep(1 * time.Second)
	}
}
