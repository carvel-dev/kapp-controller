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
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Right now the implementation of the package repo reconciler needs improvement
// hopefully after which, these tests can be cleaned up to remove retries and time related
// falkeyness

func Test_PackageRepoStatus_Failing(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	name := "repo"

	repoYaml := `apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: test-repo
spec:
  fetch:
    imgpkgBundle:
      image: k8slt/i-dont-exist`

	expectedStatus := v1alpha1.PackageRepositoryStatus{
		Fetch: &kcv1alpha1.AppStatusFetch{
			ExitCode: 1,
			Error:    "Fetching resources: Error (see .status.usefulErrorMessage for details)",
		},
		GenericStatus: kcv1alpha1.GenericStatus{
			Conditions: []kcv1alpha1.AppCondition{{
				Type:    kcv1alpha1.ReconcileFailed,
				Status:  corev1.ConditionTrue,
				Message: "Fetching resources: Error (see .status.usefulErrorMessage for details)",
			}},
			ObservedGeneration:  1,
			FriendlyDescription: "Reconcile failed: Fetching resources: Error (see .status.usefulErrorMessage for details)",
			UsefulErrorMessage:  "Error: Syncing directory '0': Syncing directory '.' with imgpkgBundle contents: Imgpkg: exit status 1 (stderr: Error: Checking if image is bundle: Collecting images: Working with index.docker.io/k8slt/i-dont-exist:latest: GET https://index.docker.io/v2/k8slt/i-dont-exist/manifests/latest: UNAUTHORIZED: authentication required; [map[Action:pull Class: Name:k8slt/i-dont-exist Type:repository]]\n)\n",
		},
		ConsecutiveReconcileFailures: 4,
	}

	cleanup := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanup()
	defer cleanup()

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

		cleanupStatusForAssertion(&cr)

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
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c`

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
		Fetch: &kcv1alpha1.AppStatusFetch{
			ExitCode: 0,
		},
		Template: &kcv1alpha1.AppStatusTemplate{
			ExitCode: 0,
		},
		Deploy: &kcv1alpha1.AppStatusDeploy{
			ExitCode: 0,
			Finished: true,
		},
		ConsecutiveReconcileSuccesses: 1,
		ConsecutiveReconcileFailures:  0,
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

		cleanupStatusForAssertion(&cr)

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
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c`

	cleanUp := func() {
		kubectl.Run([]string{"delete", "pkgr/basic.test.carvel.dev"})
	}
	defer cleanUp()

	kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(yamlRepo)})

	retry(t, 30*time.Second, func() error {
		_, err := kubectl.RunWithOpts([]string{"get", "pkgm/pkg.test.carvel.dev"}, RunOpts{AllowError: true})
		if err != nil {
			return fmt.Errorf("Expected to find pkgm pkg.test.carvel.dev, but couldn't: %v", err)
		}

		_, err = kubectl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true})
		if err != nil {
			return fmt.Errorf("Expected to find pkg (pkg.test.carvel.dev.1.0.0, pkg.test.carvel.dev.2.0.0) but couldn't: %v", err)
		}
		_, err = kubectl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true})
		if err != nil {
			return fmt.Errorf("Expected to find pkg (pkg.test.carvel.dev.1.0.0, pkg.test.carvel.dev.2.0.0) but couldn't: %v", err)
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
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.delete.test.carvel.dev
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c`

	cleanUp := func() {
		kctl.RunWithOpts([]string{"delete", "package/pkg.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true})
		kctl.RunWithOpts([]string{"delete", "package/pkg.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true})
		kctl.RunWithOpts([]string{"delete", "packagemetadata/pkg.test.carvel.dev"}, RunOpts{AllowError: true})
		kctl.RunWithOpts([]string{"delete", "pkgr/basic.test.carvel.dev"}, RunOpts{AllowError: true})
		// kapp delete still needed in event checking if packages existing fails
		kapp.Run([]string{"delete", "-a","repo"})
	}
	defer cleanUp()

	logger.Section("deploy repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", "repo"},
			RunOpts{StdinReader: strings.NewReader(repoYaml)})
	})

	logger.Section("check packages exist", func() {
		retry(t, 30*time.Second, func() error {
			_, err := kctl.RunWithOpts([]string{"get", "pkgm/pkg.test.carvel.dev"}, RunOpts{AllowError: true})
			if err != nil {
				return fmt.Errorf("Expected to find pkgm pkg.test.carvel.dev, but couldn't: %v", err)
			}

			_, err = kctl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true})
			if err != nil {
				return fmt.Errorf("Expected to find package pkg.test.carvel.dev.1.0.0: %v", err)
			}

			_, err = kctl.RunWithOpts([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true})
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
			_, err := kctl.RunWithOpts([]string{"get", "pkgm/pkg.test.carvel.dev"}, RunOpts{AllowError: true})
			if err == nil || !strings.Contains(err.Error(), "\"pkg.test.carvel.dev\" not found") {
				return fmt.Errorf("Expected not to find pkgm pkg.test.carvel.dev, but did: %v", err)
			}

			_, err = kctl.RunWithOpts([]string{"get", "package/pkg.test.carvel.dev.1.0.0"}, RunOpts{AllowError: true})
			if err == nil {
				return fmt.Errorf("Expected not to find package pkg.test.carvel.dev.1.0.0, but did")
			}

			if !strings.Contains(err.Error(), "Error from server (NotFound)") {
				return fmt.Errorf("Expected not found error for pkg.test.carvel.dev.1.0.0, but got: %v", err)
			}

			_, err = kctl.RunWithOpts([]string{"get", "package/pkg.test.carvel.dev.2.0.0"}, RunOpts{AllowError: true})
			if err == nil {
				return fmt.Errorf("Expected no to find package pkg.test.carvel.dev.2.0.0, but did")
			}

			if !strings.Contains(err.Error(), "Error from server (NotFound)") {
				return fmt.Errorf("Expected not found error for pkg.test.carvel.dev.2.0.0, but got: %v", err)
			}
			return nil
		})
	})
}

func Test_PackageRepoStatus_ShowsWithKubectlGet(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	name := "repo-status"

	repoYaml := fmt.Sprintf(`apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s
spec:
  fetch:
    imgpkgBundle:
      image: k8slt/i-dont-exist`, name)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	// deploy failing repo
	logger.Section("deploy failing repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{StdinReader: strings.NewReader(repoYaml)})
	})

	out, err := kubectl.RunWithOpts([]string{"get", "pkgr/" + name}, RunOpts{AllowError: true})
	if err != nil {
		t.Fatalf("encountered unknown error from kubectl get pkgr: %v", err)
	}

	if !strings.Contains(out, "DESCRIPTION") && !strings.Contains(out, "Reconcile failed") {
		t.Fatalf("output did not contain DESCRIPTION column from kubectl get.\nGot:\n%s", out)
	}
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

func cleanupStatusForAssertion(pkgr *v1alpha1.PackageRepository) {
	// fetch
	if pkgr.Status.Fetch != nil {
		pkgr.Status.Fetch.StartedAt = metav1.Time{}
		pkgr.Status.Fetch.UpdatedAt = metav1.Time{}
		pkgr.Status.Fetch.Stdout = ""
		pkgr.Status.Fetch.Stderr = ""
	}

	// template
	if pkgr.Status.Template != nil {
		pkgr.Status.Template.UpdatedAt = metav1.Time{}
		pkgr.Status.Template.Stderr = ""
	}

	// deploy
	if pkgr.Status.Deploy != nil {
		pkgr.Status.Deploy.StartedAt = metav1.Time{}
		pkgr.Status.Deploy.UpdatedAt = metav1.Time{}
		pkgr.Status.Deploy.Stdout = ""
		pkgr.Status.Deploy.Stderr = ""
	}
}
