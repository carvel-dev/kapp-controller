// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
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

	cleanup := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanup()
	defer cleanup()

	logger.Section("deploy failing repo", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{StdinReader: strings.NewReader(repoYaml), AllowError: true})
	})

	logger.Section("check against expected error status", func() {
		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

		var cr v1alpha1.PackageRepository
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}

		expectedStatus := v1alpha1.PackageRepositoryStatus{
			Fetch: &kcv1alpha1.AppStatusFetch{
				ExitCode: 1,
				Error:    "Fetching resources: Error (see .status.usefulErrorMessage for details)",
			},
			GenericStatus: kcv1alpha1.GenericStatus{
				Conditions: []kcv1alpha1.Condition{{
					Type:    kcv1alpha1.ReconcileFailed,
					Status:  corev1.ConditionTrue,
					Message: "Fetching resources: Error (see .status.usefulErrorMessage for details)",
				}},
				ObservedGeneration:  1,
				FriendlyDescription: "Reconcile failed: Fetching resources: Error (see .status.usefulErrorMessage for details)",
				UsefulErrorMessage:  "vendir: Error: Syncing directory '0':\n  Syncing directory '.' with imgpkgBundle contents:\n    Imgpkg: exit status 1 (stderr: imgpkg: Error: Fetching image:\n  GET https://index.docker.io/v2/k8slt/i-dont-exist/manifests/latest:\n    UNAUTHORIZED: authentication required; [map[Action:pull Class: Name:k8slt/i-dont-exist Type:repository]]\n)\n",
			},
			ConsecutiveReconcileFailures: 1,
		}

		cleanupStatusForAssertion(&cr)
		assert.Equal(t, expectedStatus, cr.Status)
	})
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
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(repoYml),
			OnErrKubectl: []string{"get", "pkgr", "-oyaml"},
		})
	})

	logger.Section("check against expected successful status", func() {
		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

		var cr v1alpha1.PackageRepository
		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}

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
				Conditions: []kcv1alpha1.Condition{{
					Type:    kcv1alpha1.ReconcileSucceeded,
					Status:  corev1.ConditionTrue,
					Message: "",
				}},
				ObservedGeneration:  1,
				FriendlyDescription: "Reconcile succeeded",
				UsefulErrorMessage:  "",
			},
		}

		cleanupStatusForAssertion(&cr)
		assert.Equal(t, expectedStatus, cr.Status)
	})

	logger.Section("force a second reconcile and see if it all still works", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(repoYml + "\n  syncPeriod: 30s\n"),
			OnErrKubectl: []string{"get", "pkgr", "-oyaml"},
		})
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

	verifyPkg := func(resourceName, imgRef string) {
		out := kubectl.Run([]string{"get", resourceName, "-o", "yaml"})
		assert.Contains(t, out, "packaging.carvel.dev/package-repository-ref: kappctrl-test/basic.test.carvel.dev")
		assert.Contains(t, out, "image: "+imgRef+"\n")
	}

	logger.Section("deploy pkg repository", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(yamlRepo)})

		out := kubectl.Run([]string{"get", "pkgm/pkg.test.carvel.dev", "-o", "yaml"})
		assert.Contains(t, out, "packaging.carvel.dev/package-repository-ref: kappctrl-test/basic.test.carvel.dev")

		verifyPkg("pkg/pkg.test.carvel.dev.1.0.0", "index.docker.io/k8slt/kctrl-example-pkg@sha256:8ffa7f9352149dba1d539d0006b38eda357917edcdd39b82497a61dab2c27b75")
		verifyPkg("pkg/pkg.test.carvel.dev.2.0.0", "index.docker.io/k8slt/kctrl-example-pkg@sha256:73713d922b5f561c0db2a7ea5f4f6384f7d2d6289886f8400a8aaf5e8fdf134a")
	})

	logger.Section("deploy pkg repository with same content from a different location (where bundle was copied)", func() {
		// Change location of the bundle to be different
		// (imgpkg copy was used to copy original repo to a new location)
		updatedRepo := `---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: basic.test.carvel.dev
spec:
  fetch:
    imgpkgBundle:
      image: index.docker.io/k8slt/kc-e2e-test-repo-copied@sha256:ddd93b67b97c1460580ca1afd04326d16900dc716c4357cade85b83deab76f1c`

		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, e2e.RunOpts{StdinReader: strings.NewReader(updatedRepo)})

		out := kubectl.Run([]string{"get", "pkgm/pkg.test.carvel.dev", "-o", "yaml"})
		assert.Contains(t, out, "packaging.carvel.dev/package-repository-ref: kappctrl-test/basic.test.carvel.dev")

		// Note that location of Packages has also changed to k8slt/kc-e2e-test-repo-copied
		// since kbld is being applied with imgpkg's images.yml relocation data
		verifyPkg("pkg/pkg.test.carvel.dev.1.0.0", "index.docker.io/k8slt/kc-e2e-test-repo-copied@sha256:8ffa7f9352149dba1d539d0006b38eda357917edcdd39b82497a61dab2c27b75")
		verifyPkg("pkg/pkg.test.carvel.dev.2.0.0", "index.docker.io/k8slt/kc-e2e-test-repo-copied@sha256:73713d922b5f561c0db2a7ea5f4f6384f7d2d6289886f8400a8aaf5e8fdf134a")
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

func Test_PackageReposWithOverlappingPackages_identicalPackagesWithUpdates(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	pkgr1Name := "repo1.tankyu.carvel.dev"
	pkgr2Name := "repo2.tankyu.carvel.dev"
	pkgName := "pkg0.test.carvel.dev.0.0.0"

	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[1]s
spec:
  fetch:
    inline:
      paths:
        packages/pkg.test.carvel.dev/pkg0.test.carvel.dev.0.0.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: %[2]s.%[3]s
          spec:
            refName: %[2]s
            version: %[3]s
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
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, "pkg0.test.carvel.dev", "0.0.0")
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, "pkg0.test.carvel.dev", "0.0.0")

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})
		kapp.Run([]string{"delete", "-a", pkgr2Name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)
	})

	logger.Section("deploy pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr2Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr2),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)
	})

	pkgrUpdateTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[1]s
spec:
  fetch:
    inline:
      paths:
        packages/pkg.test.carvel.dev/pkg0.test.carvel.dev.0.0.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: %[2]s.%[3]s
          spec:
            refName: %[2]s
            version: %[3]s
            releaseNotes: "this update just adds release notes without bumping the revision."
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
	pkgr1u := fmt.Sprintf(pkgrUpdateTemplate, pkgr1Name, "pkg0.test.carvel.dev", "0.0.0")
	pkgr2u := fmt.Sprintf(pkgrUpdateTemplate, pkgr2Name, "pkg0.test.carvel.dev", "0.0.0")

	logger.Section("updated pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1u),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)

		assert.Contains(t, kubectl.Run([]string{"get", "pkg", pkgName, "-oyaml"}), "releaseNotes")
	})

	logger.Section("update pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr2Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr2u),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)
	})
}

func Test_PackageReposWithOverlappingPackages_identicalPackagesReconcile(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	pkgr1Name := "repo1.tankyu.carvel.dev"
	pkgr2Name := "repo2.tankyu.carvel.dev"
	pkgr3Name := "repo3.tankyu.carvel.dev"
	pkgr4Name := "repo4.tankyu.carvel.dev"
	pkgr5Name := "repo5.tankyu.carvel.dev"
	pkgName := "pkg0.test.carvel.dev.0.0.0"

	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s
spec:
  fetch:
    inline:
      paths:

        packages/pkg.test.carvel.dev/pkg0.test.carvel.dev.0.0.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: %s
          spec:
            refName: pkg0.test.carvel.dev
            version: 0.0.0
            template:
              spec: {}
`
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, pkgName)
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, pkgName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})
		kapp.Run([]string{"delete", "-a", pkgr2Name})
		kapp.Run([]string{"delete", "-a", pkgr3Name})
		kapp.Run([]string{"delete", "-a", pkgr4Name})
		kapp.Run([]string{"delete", "-a", pkgr5Name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)
	})

	logger.Section("deploy pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr2Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr2),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)
	})

	// test cases where the two packages aren't quite identical so there's still an error

	logger.Section("deploy pkgr3 but it fails because the annotations are different", func() {
		pkgNameAndAnn := fmt.Sprintf("%s\n            annotations: {some.co.internal.ann: value}", pkgName)
		pkgr3 := fmt.Sprintf(pkgrTemplate, pkgr3Name, pkgNameAndAnn)

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr3Name},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr3), AllowError: true})
		require.Error(t, err)

		out := string(kubectl.Run([]string{"get", "pkgr", pkgr3Name, "-oyaml"}))
		assert.Contains(t, out, "is already present but not identical (mismatch in metadata.annotations)")
	})

	logger.Section("deploy pkgr4 but it fails because the labels are different", func() {
		pkgNameAndLabel := fmt.Sprintf("%s\n            labels: {some.co.internal.label: label-value}", pkgName)
		pkgr4 := fmt.Sprintf(pkgrTemplate, pkgr4Name, pkgNameAndLabel)

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr4Name},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr4), AllowError: true})
		require.Error(t, err)

		out := string(kubectl.Run([]string{"get", "pkgr", pkgr4Name, "-oyaml"}))
		assert.Contains(t, out, "is already present but not identical (mismatch in metadata.labels)")
	})

	logger.Section("deploy pkgr5 but it fails because the package specs are different", func() {
		pkgr5 := strings.Replace(fmt.Sprintf(pkgrTemplate, pkgr5Name, pkgName),
			`
              spec: {}
`,
			`
              spec:
                fetch:
                - imgpkgBundle:
                    image: k8slt/kctrl-example-pkg:v1.0.0
`,
			-1)

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr5Name},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr5), AllowError: true})
		require.Error(t, err)

		out := string(kubectl.Run([]string{"get", "pkgr", pkgr5Name, "-oyaml"}))
		assert.Contains(t, out, "is already present but not identical (mismatch in spec.template)")
	})
}

func Test_PackageReposWithOverlappingPackages_localAndGlobalNS(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	pkgName := "pkg0.test.carvel.dev.0.0.0"
	pkgr1Name := "repo1.tankyu.carvel.dev"
	pkgr2Name := "repo2.tankyu.carvel.dev"

	pkgr1LocalNS := env.Namespace
	pkgr2GlobalNS := env.PackagingGlobalNS

	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s
  namespace: %s
spec:
  fetch:
    inline:
      paths:

        packages/pkg.test.carvel.dev/pkg0.test.carvel.dev.0.0.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: %s
          spec:
            refName: pkg0.test.carvel.dev
            version: 0.0.0
            template:
              spec: {}
`
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, pkgr1LocalNS, pkgName)
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, pkgr2GlobalNS, pkgName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})
		kapp.Run([]string{"delete", "-a", pkgr2Name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1 into local namespace", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, pkgr1LocalNS)
	})

	logger.Section("deploy pkgr2 into global namespace successfully, but pkg is still owned by pkgr1 in local namespace", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr2Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr2),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwnedWithNs(t, kubectl, pkgName, pkgr1LocalNS, pkgr1Name, pkgr1LocalNS)
		assertPkgOwnedWithNs(t, kubectl, pkgName, pkgr2GlobalNS, pkgr2Name, pkgr2GlobalNS)
	})

	logger.Section("redeploy pkgr1 after pkgr exists in global namespace", func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})

		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwnedWithNs(t, kubectl, pkgName, pkgr2GlobalNS, pkgr2Name, pkgr2GlobalNS)
		// package in local namespace will be owned by local pkgr because
		// kapp deploy (within pkgr) will issue a create call which
		// will be accepted by the packaging API server. create call
		// is done optimistically (no get is performed first).
		assertPkgOwnedWithNs(t, kubectl, pkgName, pkgr1LocalNS, pkgr1Name, pkgr1LocalNS)
	})
}

func Test_PackageReposWithOverlappingPackages_packagesHaveDifferentRevisions(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	pkgName := "pkg0.test.carvel.dev.0.0.0"
	pkgr1Name := "repo1.tankyu.carvel.dev"
	pkgr2Name := "repo2.tankyu.carvel.dev"
	pkgr3Name := "repo3.tankyu.carvel.dev"
	pkgr4Name := "repo4.tankyu.carvel.dev"
	pkgr5Name := "repo5.tankyu.carvel.dev"

	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s
spec:
  fetch:
    inline:
      paths:

        packages/pkg.test.carvel.dev/pkg0.test.carvel.dev.0.0.0.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: %s
            annotations:
              packaging.carvel.dev/revision: "%s"
          spec:
            refName: pkg0.test.carvel.dev
            version: 0.0.0
            template:
              spec: {}
`
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, pkgName, "1")

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})
		kapp.Run([]string{"delete", "-a", pkgr2Name})
		kapp.Run([]string{"delete", "-a", pkgr3Name})
		kapp.Run([]string{"delete", "-a", pkgr4Name})
		kapp.Run([]string{"delete", "-a", pkgr5Name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr1Name, env.Namespace)
	})

	logger.Section("deploy pkgr2 successfully, and it overrides bc it has higher revision", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr2Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(fmt.Sprintf(pkgrTemplate, pkgr2Name, pkgName, "2")),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr2Name, env.Namespace)
	})

	logger.Section("uninstall and reinstall pkgr1 but it never takes ownership of the pkg from pkgr2 bc it has lower revision", func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})

		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr2Name, env.Namespace)
	})

	logger.Section("install pkgr with higher revision using .0 suffix (2.0 > 2)", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr3Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(fmt.Sprintf(pkgrTemplate, pkgr3Name, pkgName, "2.0")),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr3Name, env.Namespace)
	})

	logger.Section("install pkgr with lower revision (2.0 > 1.6.8)", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr4Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(fmt.Sprintf(pkgrTemplate, pkgr4Name, pkgName, "1.6.8")),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr3Name, env.Namespace)
	})

	logger.Section("install pkgr with higher revision (2.1.0 > 2.0)", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr5Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(fmt.Sprintf(pkgrTemplate, pkgr5Name, pkgName, "2.1.0")),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})
		assertPkgOwned(t, kubectl, pkgName, pkgr5Name, env.Namespace)
	})
}

func Test_PackageReposWithOverlappingPackages_NonTrivialPackages(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	pkgr1Name := "repo-1.tankyu.carvel.dev"
	pkgr2Name := "repo-2.tankyu.carvel.dev"
	pkgr3Name := "repo-3.tankyu.carvel.dev"
	pkgr4Name := "repo-4.tankyu.carvel.dev"

	pkgrPreamble := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: repo-%d.tankyu.carvel.dev
spec:
  fetch:
    inline:
      paths:
`
	pkgTemplate := `
        packages/pkg.test.carvel.dev/%[1]s.%[2]s.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: Package
          metadata:
            name: %[1]s.%[2]s
          spec:
            refName: %[1]s
            version: %[2]s
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
	pkgMetadataTemplate := `
        packages/pkg.test.carvel.dev/metadata.%[1]s.yml: |
          ---
          apiVersion: data.packaging.carvel.dev/v1alpha1
          kind: PackageMetadata
          metadata:
            name: %[1]s
          spec:
            displayName: %[2]s
            shortDescription: %[3]s
            providerName: myCorp`
	pkgr1 := fmt.Sprintf(
		"%s%s%s%s%s%s",
		fmt.Sprintf(pkgrPreamble, 1),
		fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.5.5"),
		fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.6.0"),
		fmt.Sprintf(pkgTemplate, "coredino.co.uk", "32.76.7"),
		fmt.Sprintf(pkgMetadataTemplate, "shirt-mgr.co.uk", "shirt manager", "get your shirts to work for you"),
		fmt.Sprintf(pkgMetadataTemplate, "coredino.co.uk", "Core Dino", "i dunno its a coreDNS joke"),
	)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", pkgr1Name})
		kapp.Run([]string{"delete", "-a", pkgr2Name})
		kapp.Run([]string{"delete", "-a", pkgr3Name})
		kapp.Run([]string{"delete", "-a", pkgr4Name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr1Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr1),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})

		out := kubectl.Run([]string{"get", "packages"})
		require.Contains(t, out, "shirt-mgr.co.uk.5.5.5")
		require.Contains(t, out, "shirt-mgr.co.uk.5.6.0")
		require.Contains(t, out, "coredino.co.uk.32.76.7")

		assertPkgOwned(t, kubectl, "shirt-mgr.co.uk.5.5.5", pkgr1Name, env.Namespace)
		assertPkgOwned(t, kubectl, "shirt-mgr.co.uk.5.6.0", pkgr1Name, env.Namespace)
		assertPkgOwned(t, kubectl, "coredino.co.uk.32.76.7", pkgr1Name, env.Namespace)
	})

	logger.Section("deploy pkgr2 with partial package overlap", func() {
		pkgr2 := fmt.Sprintf(
			"%s%s%s",
			fmt.Sprintf(pkgrPreamble, 2),
			fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.5.5"),
			fmt.Sprintf(pkgTemplate, "contooor.co.uk", "0.22.0"))

		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr2Name}, e2e.RunOpts{
			StdinReader:  strings.NewReader(pkgr2),
			OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"},
		})

		assertPkgOwned(t, kubectl, "shirt-mgr.co.uk.5.5.5", pkgr1Name, env.Namespace)
		assertPkgOwned(t, kubectl, "contooor.co.uk.0.22.0", pkgr2Name, env.Namespace)

		// non-overlapping havent been affected
		assertPkgOwned(t, kubectl, "shirt-mgr.co.uk.5.6.0", pkgr1Name, env.Namespace)
		assertPkgOwned(t, kubectl, "coredino.co.uk.32.76.7", pkgr1Name, env.Namespace)
	})

	logger.Section("pkgr3 should fail to reconcile because the contooor package will conflict in the spec", func() {
		pkgr3 := fmt.Sprintf(
			"%s%s%s",
			fmt.Sprintf(pkgrPreamble, 3),
			fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.5.5"),
			strings.Replace(fmt.Sprintf(pkgTemplate, "contooor.co.uk", "0.22.0"), "k8slt/kctrl-example-pkg:v1.0.0", "k8slt/some-other-image:latest", -1))

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr3Name},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr3), AllowError: true})
		assert.Error(t, err)

		out := string(kubectl.Run([]string{"get", "pkgr", pkgr3Name, "-oyaml"}))
		assert.Contains(t, out, "Conflicting resources: Package/contooor.co.uk.0.22.0 is already present but not identical (mismatch in spec.template)")

		assertPkgOwned(t, kubectl, "shirt-mgr.co.uk.5.5.5", pkgr1Name, env.Namespace)
		assertPkgOwned(t, kubectl, "contooor.co.uk.0.22.0", pkgr2Name, env.Namespace)
	})

	logger.Section("pkgr4 will fail because of PackageMetadatas conflict", func() {
		pkgr4 := fmt.Sprintf(
			"%s%s%s",
			fmt.Sprintf(pkgrPreamble, 4),
			fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.6.0"),
			fmt.Sprintf(pkgMetadataTemplate, "shirt-mgr.co.uk", "shirt manager", "now with dress shirts"))

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", pkgr4Name},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr4), AllowError: true})
		assert.Error(t, err)

		out := string(kubectl.Run([]string{"get", "pkgr", pkgr4Name, "-oyaml"}))
		assert.Contains(t, out, " Conflicting resources: PackageMetadata/shirt-mgr.co.uk is already present but not identical (mismatch in spec.shortDescription)")
	})
}

func assertPkgOwned(t *testing.T, kubectl e2e.Kubectl, pkgName, pkgrName, pkgrNs string) {
	out := kubectl.Run([]string{"get", "package", pkgName, "-oyaml"})

	expectedOwner := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s\n", pkgrNs, pkgrName)
	assert.Contains(t, out, expectedOwner, "package is not owned by expected pkgr")
}

func assertPkgOwnedWithNs(t *testing.T, kubectl e2e.Kubectl, pkgName, pkgNs, pkgrName, pkgrNs string) {
	out, _ := kubectl.RunWithOpts([]string{"get", "package", pkgName, "-oyaml", "-n", pkgNs}, e2e.RunOpts{NoNamespace: true})

	expectedOwner := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s\n", pkgrNs, pkgrName)
	assert.Contains(t, out, expectedOwner, "package is not owned by expected pkgr")
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
