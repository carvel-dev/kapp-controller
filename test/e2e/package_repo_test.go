// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// TODO: Right now the implementation of the package repo reconciler needs improvement
// hopefully after which, these tests can be cleaned up to remove retries and time related
// falkeyness

func Test_PackageRepoStatus_Failing(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	name := "repo"

	repoYaml := `apiVersion: install.package.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: test-repo
spec:
  fetch:
    imgpkgBundle:
      image: k8slt/i-dont-exist`

	expectedStatus := v1alpha1.PackageRepositoryStatus{
		GenericStatus: kcv1alpha1.GenericStatus{
			Conditions: []kcv1alpha1.AppCondition{{
				Type:    kcv1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Syncing packages: (see .status.usefulErrorMessage for details)",
			}},
			ObservedGeneration:  1,
			FriendlyDescription: "Reconcile failed: Syncing packages: (see .status.usefulErrorMessage for details)",
			UsefulErrorMessage:  "Error: Syncing directory '0': Syncing directory '.' with imgpkgBundle contents: Imgpkg: exit status 1 (stderr: Error: Checking if image is bundle: Collecting images: Working with index.docker.io/k8slt/i-dont-exist:latest: GET https://index.docker.io/v2/k8slt/i-dont-exist/manifests/latest: UNAUTHORIZED: authentication required; [map[Action:pull Class: Name:k8slt/i-dont-exist Type:repository]]\n)\n",
		},
	}

	// deploy failing repo
	logger.Section("deploy failing repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{StdinReader: strings.NewReader(repoYaml)})
	})

	retry(t, 30*time.Second, func() error {
		// fetch repo
		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

		var cr v1alpha1.PackageRepository
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}

		// assert on expectedStatus
		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			return fmt.Errorf("\nstatus is not same:\nExpected:\n%#v\nGot:\n%#v", expectedStatus, cr.Status)
		}
		return nil
	})
}

func Test_PackageRepoStatus_Success(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	name := "test-repo-status-success"

	repoYml := `---
apiVersion: install.package.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
  # cluster scoped
spec:
  fetch:
    imgpkgBundle:
      image: k8slt/kappctrl-e2e-repo-bundle`

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{StdinReader: strings.NewReader(repoYml)})
	})

	expectedStatus := v1alpha1.PackageRepositoryStatus{
		GenericStatus: kcv1alpha1.GenericStatus{
			Conditions: []kcv1alpha1.AppCondition{{
				Type:    kcv1alpha1.ReconcileSucceeded,
				Status:  corev1.ConditionTrue,
				Message: "",
			}},
			ObservedGeneration:  1,
			FriendlyDescription: "Reconcile succeeded",
			UsefulErrorMessage:  "",
		},
	}

	retry(t, 30*time.Second, func() error {
		// fetch repo
		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

		var cr v1alpha1.PackageRepository
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}

		// assert on expectedStatus
		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			return fmt.Errorf("\nstatus is not same:\nExpected:\n%#v\nGot:\n%#v", expectedStatus, cr.Status)
		}
		return nil
	})
}

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
    imgpkgBundle:
      image: k8slt/kappctrl-e2e-repo-bundle`

	cleanUp := func() {
		kubectl.RunWithOpts([]string{"delete", "pkgr/basic.test.carvel.dev"}, RunOpts{NoNamespace: true})
	}
	defer cleanUp()

	kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(yamlRepo)})

	retry(t, 10*time.Second, func() error {
		_, err := kubectl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"}, RunOpts{NoNamespace: true, AllowError: true})
		if err != nil {
			return fmt.Errorf("Expected to find pkgs (pkg.test.carvel.dev.1.0.0, pkg.test.carvel.dev.2.0.0) but couldn't: %v", err)
		}
		_, err = kubectl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"}, RunOpts{NoNamespace: true, AllowError: true})
		if err != nil {
			return fmt.Errorf("Expected to find pkgs (pkg.test.carvel.dev.1.0.0, pkg.test.carvel.dev.2.0.0) but couldn't: %v", err)
		}
		return nil
	})
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
    imgpkgBundle:
      image: k8slt/kappctrl-e2e-repo-bundle`

	packageNames := []string{"pkg.test.carvel.dev.1.0.0", "pkg.test.carvel.dev.2.0.0"}

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
		retry(t, 20*time.Second, func() error {
			_, err := kctl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err != nil {
				return fmt.Errorf("Expected to find package pkg.test.carvel.dev.1.0.0: %v", err)
			}

			_, err = kctl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err != nil {
				return fmt.Errorf("Expected to find package pkg.test.carvel.dev.2.0.0: %v", err)
			}
			return nil
		})
	})

	logger.Section("delete repo", func() {
		kapp.Run([]string{"delete", "-a", "repo"})
	})

	logger.Section("check packages are deleted too", func() {
		retry(t, 10*time.Second, func() error {
			_, err := kctl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err == nil || !strings.Contains(err.Error(), "\"pkg.test.carvel.dev.1.0.0\" not found") {
				return fmt.Errorf("Expected not to find package pkg.test.carvel.dev.1.0.0, but did")
			}

			_, err = kctl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true, NoNamespace: true})
			if err == nil || !strings.Contains(err.Error(), "\"pkg.test.carvel.dev.2.0.0\" not found") {
				return fmt.Errorf("Expected no to find package pkg.test.carvel.dev.2.0.0, but did")
			}
			return nil
		})
	})
}

func retry(t *testing.T, timeout time.Duration, f func() error) {
	var err error
	stopTime := time.Now().Add(timeout)
	for {
		err = f()
		if err == nil {
			return
		}
		if time.Now().After(stopTime) {
			t.Fatalf("retry timed out after %s: %v", timeout.String(), err)
		}
		time.Sleep(1 * time.Second)
	}
}
