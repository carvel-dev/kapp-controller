// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func Test_PackageRepoStatus_Failing(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
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
			UsefulErrorMessage:  "vendir: Error: Syncing directory '0':\n  Syncing directory '.' with imgpkgBundle contents:\n    Imgpkg: exit status 1 (stderr: imgpkg: Error: Checking if image is bundle:\n  Fetching image:\n    GET https://index.docker.io/v2/k8slt/i-dont-exist/manifests/latest:\n      UNAUTHORIZED: authentication required; [map[Action:pull Class: Name:k8slt/i-dont-exist Type:repository]]\n)\n",
		},
		ConsecutiveReconcileFailures: 1,
	}

	cleanup := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanup()
	defer cleanup()

	// deploy failing repo
	logger.Section("deploy failing repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{StdinReader: strings.NewReader(repoYaml), AllowError: true})
	})

	// fetch repo
	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

	var cr v1alpha1.PackageRepository
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	cleanupStatusForAssertion(&cr)

	// assert on expectedStatus
	assert.Equal(t, expectedStatus, cr.Status)
}

func Test_PackageRepoStatus_Success(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
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
			e2e.RunOpts{StdinReader: strings.NewReader(repoYml)})
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

	// fetch repo
	out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

	var cr v1alpha1.PackageRepository
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	cleanupStatusForAssertion(&cr)

	// assert on expectedStatus
	if !reflect.DeepEqual(expectedStatus, cr.Status) {
		t.Fatalf("\nstatus is not same:\nExpected:\n%#v\nGot:\n%#v", expectedStatus, cr.Status)
	}
}

func Test_PackageReposWithSamePackagesButTheyreIdentical(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller-packaging-global", logger}
	kapp := e2e.Kapp{t, env.Namespace, logger}

	name1 := "repo1"
	cleanUp1 := func() {
		kapp.Run([]string{"delete", "-a", name1})
	}
	defer cleanUp1()

	logger.Section("deploy PackageRepository 1", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name1, "-f", "../assets/kc-multi-repo/package-repository1.yml"}, e2e.RunOpts{AllowError: true})
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "packagerepository/e2e-repo1.test.carvel.dev", "-o", "yaml"}, e2e.RunOpts{AllowError: true}))
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "package", "-n", "kapp-controller-packaging-global"}, e2e.RunOpts{AllowError: true}))
	})

	logger.Section("assert packages were installed", func() {
		out := kubectl.Run([]string{"get", "packages"})
		// fmt.Println("kubectl get packages output: ", out)
		require.Contains(t, out, "pkg.test.carvel.dev.1.0.0")
		require.Contains(t, out, "pkg.test.carvel.dev.2.0.0")
		require.Contains(t, out, "pkg.test.carvel.dev.3.0.0-rc.1")
	})

	name2 := "repo2"
	cleanUp2 := func() {
		kapp.Run([]string{"delete", "-a", name2})
	}
	defer cleanUp2()

	logger.Section("deploy PackageRepository 2", func() {
		// kapp.RunWithOpts([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/package-repository2.yml"}, e2e.RunOpts{AllowError: true})
		// fmt.Println(kapp.Run([]string{"inspect", "-a", name2}))
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "packagerepository/e2e-repo2.test.carvel.dev", "-o", "yaml"}, e2e.RunOpts{AllowError: true}))
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "package"}, e2e.RunOpts{AllowError: true}))
		kapp.RunWithOpts([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/package-repository2.yml"}, e2e.RunOpts{AllowError: false})
	})

}

// TODO: add a test that actually tests the "package without a repo creates a conflict" case

func Test_PackageReposWithSamePackagesButTheNewOneHasHigherRev(t *testing.T) {
}
func Test_PackageReposWithSamePackageAsStandalonePackage(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller-packaging-global", logger}
	kapp := e2e.Kapp{t, env.Namespace, logger}

	name1 := "standalonepkg"
	cleanUp1 := func() {
		kapp.Run([]string{"delete", "-a", name1})
	}
	defer cleanUp1()

	logger.Section("deploy Standalone Package", func() {
		kapp.Run([]string{"deploy", "-a", name1, "-f", "../assets/kc-multi-repo/non-repo-package.yml"})
		// fmt.Println(kapp.Run([]string{"inspect", "-a", name1}))
		out := kubectl.Run([]string{"get", "packages"})
		// fmt.Println("kubectl get packages: ", out)
		require.Contains(t, out, "pkg.test.carvel.dev.2.0.0")
	})

	name2 := "pkgr-with-conflicting-pkg"
	cleanUp2 := func() {
		kapp.Run([]string{"delete", "-a", name2})
	}
	defer cleanUp2()

	logger.Section("deploy PackageRepository 2", func() {
		out, err := kapp.RunWithOpts([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/inline-repo2.yml"}, e2e.RunOpts{AllowError: true})
		// fmt.Println("output of kapp deploy of pkgr2: ", out)
		// fmt.Println("error from kapp deploy of pkgr2: ", err)
		assert.Error(t, err)
		assert.Contains(t, out, "Reconcile failed:  (message: Deploying: Error")
	})

	/*

		// TODO: i thought this would be cute to show causality of like, "look if you fake the package being in a repo then it works" but tbh it's very hard to tell why it isn't working.

		logger.Section("annotate package to make it look like it came from a repo", func() {
			kubectl.Run([]string{"annotate", "pkg", "pkg.test.carvel.dev.2.0.0", "packaging.carvel.dev/package-repository-ref=foo/bar.tanzu.carvel.dev"})
		})

		logger.Section("deploy PackageRepository 2 but this time it works", func() {
			out, err := kapp.RunWithOpts([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/inline-repo2.yml"}, e2e.RunOpts{AllowError: true})
			fmt.Println("output of kapp deploy of pkgr2: ", out)
			fmt.Println("error from kapp deploy of pkgr2: ", err)
			fmt.Println(kubectl.Run([]string{"get", "pkgr", "-A", "-o", "yaml"}))
			kapp.Run([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/inline-repo2.yml"})
		})
	*/
}

func Test_PackageReposWithSamePackagesButTheOldOneHasHigherRev(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller-packaging-global", logger}
	kapp := e2e.Kapp{t, env.Namespace, logger}

	name1 := "repo2"
	cleanUp1 := func() {
		kapp.Run([]string{"delete", "-a", name1})
	}
	defer cleanUp1()

	logger.Section("deploy PackageRepository 2", func() {
		kapp.Run([]string{"deploy", "-a", name1, "-f", "../assets/kc-multi-repo/inline-repo2.yml"})
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "package", "-n", "kapp-controller-packaging-global"}, e2e.RunOpts{AllowError: true}))
	})

	logger.Section("assert packages were installed", func() {
		out := kubectl.Run([]string{"get", "packages", "-o", "yaml"})
		// fmt.Println("kubectl get packages output: ", out)
		require.Contains(t, out, "pkg.test.carvel.dev")
		require.Contains(t, out, "this is rev 3 now")
	})

	name2 := "repo1"
	cleanUp2 := func() {
		kapp.Run([]string{"delete", "-a", name2})
	}
	defer cleanUp2()

	logger.Section("deploy PackageRepository 1 should result in keeping the existing package but with no error message", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/inline-repo1.yml"}, e2e.RunOpts{AllowError: true})
		// fmt.Println(kapp.Run([]string{"inspect", "-a", name2}))
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "packagerepository/test-repo1.tanzu.carvel.dev", "-o", "yaml"}, e2e.RunOpts{AllowError: true}))
		// fmt.Println(kubectl.RunWithOpts([]string{"get", "package"}, e2e.RunOpts{AllowError: true}))
		kapp.RunWithOpts([]string{"deploy", "-a", name2, "-f", "../assets/kc-multi-repo/inline-repo1.yml"}, e2e.RunOpts{AllowError: false})

		out := kubectl.Run([]string{"get", "packages", "-o", "yaml"})
		// fmt.Println("kubectl get packages output: ", out)
		require.Contains(t, out, "pkg.test.carvel.dev")
		require.Contains(t, out, "this is rev 3 now")
	})

}

func Test_PackageRepoBundle_PackagesAvailable(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	kapp := e2e.Kapp{t, env.Namespace, logger}
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

	name := "repo-packages-available"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	defer cleanUp()

	logger.Section("deploy PackageRepository", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(yamlRepo)})
	})

	logger.Section("check PackageMetadata/Packages created", func() {
		verify := func(resourceName string) {
			kctlOutput := kubectl.Run([]string{"get", resourceName, "-o", "yaml"})
			assert.Contains(t, kctlOutput, "kapp.k14s.io/identity:")
			assert.Contains(t, kctlOutput, "packaging.carvel.dev/package-repository-ref: kappctrl-test/basic.test.carvel.dev")
		}
		verify("pkgm/pkg.test.carvel.dev")
		verify("pkg/pkg.test.carvel.dev.1.0.0")
		verify("pkg/pkg.test.carvel.dev.2.0.0")
	})
}

func Test_PackageRepoDelete(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kctl := e2e.Kubectl{t, env.Namespace, logger}

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
		kapp.Run([]string{"delete", "-a", "repo"})
	}
	defer cleanUp()

	logger.Section("deploy repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", "repo"},
			e2e.RunOpts{StdinReader: strings.NewReader(repoYaml)})
	})

	logger.Section("check Packages exist", func() {
		kctl.Run([]string{"get", "pkgm/pkg.test.carvel.dev"})
		kctl.Run([]string{"get", "pkg/pkg.test.carvel.dev.1.0.0"})
		kctl.Run([]string{"get", "pkg/pkg.test.carvel.dev.2.0.0"})
	})

	logger.Section("delete repo", func() {
		kapp.Run([]string{"delete", "-a", "repo"})
	})

	logger.Section("check packages are deleted too", func() {
		_, err := kctl.RunWithOpts([]string{"get", "pkgm/pkg.test.carvel.dev"}, e2e.RunOpts{AllowError: true})
		if err == nil || !strings.Contains(err.Error(), "\"pkg.test.carvel.dev\" not found") {
			t.Fatalf("Expected not to find pkgm pkg.test.carvel.dev, but did: %v", err)
		}

		_, err = kctl.RunWithOpts([]string{"get", "package/pkg.test.carvel.dev.1.0.0"}, e2e.RunOpts{AllowError: true})
		if err == nil {
			t.Fatalf("Expected not to find package pkg.test.carvel.dev.1.0.0, but did")
		}
		if !strings.Contains(err.Error(), "Error from server (NotFound)") {
			t.Fatalf("Expected not found error for pkg.test.carvel.dev.1.0.0, but got: %v", err)
		}

		_, err = kctl.RunWithOpts([]string{"get", "package/pkg.test.carvel.dev.2.0.0"}, e2e.RunOpts{AllowError: true})
		if err == nil {
			t.Fatalf("Expected no to find package pkg.test.carvel.dev.2.0.0, but did")
		}

		if !strings.Contains(err.Error(), "Error from server (NotFound)") {
			t.Fatalf("Expected not found error for pkg.test.carvel.dev.2.0.0, but got: %v", err)
		}
	})
}

func Test_PackageRepoStatus_ShowsWithKubectlGet(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
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
			e2e.RunOpts{StdinReader: strings.NewReader(repoYaml), AllowError: true})
	})

	out, err := kubectl.RunWithOpts([]string{"get", "pkgr/" + name}, e2e.RunOpts{AllowError: true})
	if err != nil {
		t.Fatalf("encountered unknown error from kubectl get pkgr: %v", err)
	}

	if !strings.Contains(out, "DESCRIPTION") && !strings.Contains(out, "Reconcile failed") {
		t.Fatalf("output did not contain DESCRIPTION column from kubectl get.\nGot:\n%s", out)
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
