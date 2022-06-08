// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

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
			Conditions: []kcv1alpha1.Condition{{
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
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
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
			e2e.RunOpts{StdinReader: strings.NewReader(repoYml), OnErrKubectl: []string{"get", "pkgr", "-oyaml"}})
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

	var cr v1alpha1.PackageRepository

	populateCrStatus := func() {
		// fetch repo
		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageRepository"})

		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}
		cleanupStatusForAssertion(&cr)
	}

	populateCrStatus()
	// assert on expectedStatus
	if !reflect.DeepEqual(expectedStatus, cr.Status) {
		t.Fatalf("\nstatus is not same:\nExpected:\n%#v\nGot:\n%#v", expectedStatus, cr.Status)
	}

	logger.Section("force a second reconcile and see if it all still works", func() {
		resourceVersion1 := kubectl.Run([]string{"get", "pkgr", "-o=jsonpath='{.items[0].metadata.resourceVersion}'"})
		kubectl.Run([]string{"annotate", "pkgr", "basic.test.carvel.dev", "foo=value"})

		var out string
		resourceVersion2 := resourceVersion1
		for i := 1; i < 5; i++ {
			out = kubectl.Run([]string{"get", "pkgr", "basic.test.carvel.dev", "-oyaml"})

			// if you can tell that the resourceVersion has incremented, that means the reconciler ran and we can exit the loop
			resourceVersion2 = kubectl.Run([]string{"get", "pkgr", "-o=jsonpath='{.items[0].metadata.resourceVersion}'"})
			if resourceVersion2 != resourceVersion1 {
				break
			}
			time.Sleep(time.Duration(i) * time.Second)
		}
		assert.NotEqual(t, resourceVersion1, resourceVersion2, "failed to observe a second reconciliation of the PKGR")
		assert.Contains(t, out, "Reconcile succeeded")
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

func Test_PackageReposWithOverlappingPackages_identicalPackagesWithUpdates(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	pkgName := "pkg0.test.carvel.dev.0.0.0"
	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[1]s.tankyu.carvel.dev
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

	pkgr1Name := "repo1"
	pkgr2Name := "repo2"
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, "pkg0.test.carvel.dev", "0.0.0")
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, "pkg0.test.carvel.dev", "0.0.0")

	assertPkgOwnedBy1 := func() {
		out := kubectl.Run([]string{"get", "package", pkgName, "-oyaml"})
		expectedOwnership := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s", env.Namespace, pkgr1Name)
		assert.Contains(t, out, expectedOwnership,
			"\n========\n Package Ownership Check Failed: expected ", expectedOwnership, "\n=======")
	}

	appName1 := "pkgr-1"
	appName2 := "pkgr-2"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName1})
		kapp.Run([]string{"delete", "-a", appName2})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out, _ := kubectl.RunWithOpts([]string{"get", "packages", "-A"}, e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, "pkg0.test.carvel.dev.0.0.0")

		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy1()
	})

	logger.Section("deploy pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName2},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr2), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out := kapp.Run([]string{"inspect", "-a", appName2, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy1()
	})
	pkgrUpdate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %[1]s.tankyu.carvel.dev
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

	pkgr1u := fmt.Sprintf(pkgrUpdate, pkgr1Name, "pkg0.test.carvel.dev", "0.0.0")
	pkgr2u := fmt.Sprintf(pkgrUpdate, pkgr2Name, "pkg0.test.carvel.dev", "0.0.0")

	logger.Section("updated pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1u), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out := kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy1()

		out = kubectl.Run([]string{"get", "pkg", pkgName, "-oyaml"})
		assert.Contains(t, out, "releaseNotes")
	})

	logger.Section("update pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName2},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr2u), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out := kapp.Run([]string{"inspect", "-a", appName2, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy1()
	})
}

func Test_PackageReposWithOverlappingPackages_identicalPackagesReconcile(t *testing.T) {
	pkgName := "pkg0.test.carvel.dev.0.0.0"
	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s.tankyu.carvel.dev
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
	pkgr1Name := "repo1"
	pkgr2Name := "repo2"
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, pkgName)
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, pkgName)

	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	appName1 := "pkgr-1"
	appName2 := "pkgr-2"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName1})
		kapp.Run([]string{"delete", "-a", appName2})
	}
	cleanUp()
	defer cleanUp()

	assertPkgOwnedBy := func(pkgrName string) {
		out := kubectl.Run([]string{"get", "package", pkgName, "-oyaml"})

		expectedOwnership := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s", env.Namespace, pkgrName)
		assert.Contains(t, out, expectedOwnership,
			"\n========\n Package Ownership Check Failed: expected ", expectedOwnership, "\n=======")
	}

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out, _ := kubectl.RunWithOpts([]string{"get", "packages", "-A"}, e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, "pkg0.test.carvel.dev.0.0.0")

		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name)
	})

	logger.Section("deploy pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName2},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr2), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out := kapp.Run([]string{"inspect", "-a", appName2, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name)
	})

	// test cases where the two packages aren't quite identical so there's still an error
	logger.Section("deploy pkgr3 but it fails because the annotations are different", func() {
		pkgr3Name := "repo3"
		pkgNameAndAnn := fmt.Sprintf("%s\n            annotations: {some.co.internal.ann: value}", pkgName)
		pkgr3 := fmt.Sprintf(pkgrTemplate, pkgr3Name, pkgNameAndAnn)
		appName3 := "pkgr-3"

		defer func() {
			kapp.Run([]string{"delete", "-a", appName3})
		}()

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName3},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr3), AllowError: true})
		assert.Error(t, err)
		if err != nil {
			out := string(kubectl.Run([]string{"get", "pkgr", fmt.Sprintf("%s.tankyu.carvel.dev", pkgr3Name), "-oyaml"}))
			assert.Contains(t, out, "is already present but not identical (mismatch in metadata.annotations)")
		}
	})
	logger.Section("deploy pkgr4 but it fails because the labels are different", func() {
		pkgr4Name := "repo4"
		pkgNameAndLabel := fmt.Sprintf("%s\n            labels: {some.co.internal.label: label-value}", pkgName)
		pkgr4 := fmt.Sprintf(pkgrTemplate, pkgr4Name, pkgNameAndLabel)
		appName4 := "pkgr-4"

		defer func() {
			kapp.Run([]string{"delete", "-a", appName4})
		}()

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName4},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr4), AllowError: true})
		assert.Error(t, err)
		if err != nil {
			out := string(kubectl.Run([]string{"get", "pkgr", fmt.Sprintf("%s.tankyu.carvel.dev", pkgr4Name), "-oyaml"}))
			assert.Contains(t, out, "is already present but not identical (mismatch in metadata.labels)")
		}
	})
	logger.Section("deploy pkgr5 but it fails because the specs are different", func() {
		pkgr5Name := "repo5"
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
		appName5 := "pkgr-5"

		defer func() {
			kapp.Run([]string{"delete", "-a", appName5})
		}()

		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName5},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr5), AllowError: true})
		assert.Error(t, err)
		if err != nil {
			out := string(kubectl.Run([]string{"get", "pkgr", fmt.Sprintf("%s.tankyu.carvel.dev", pkgr5Name), "-oyaml"}))
			assert.Contains(t, out, "is already present but not identical (mismatch in spec.template)")
		}
	})
}

func Test_PackageReposWithOverlappingPackages_localAndGlobalNS(t *testing.T) {
	pkgName := "pkg0.test.carvel.dev.0.0.0"
	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s.tankyu.carvel.dev
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
	pkgr1Name := "repo1"
	pkgr2Name := "repo2"
	env := e2e.BuildEnv(t)
	pkgr1NS := env.Namespace
	pkgr2NS := env.PackagingGlobalNS
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, pkgr1NS, pkgName)
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, pkgr2NS, pkgName)

	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	nsApp := `
apiVersion: v1
kind: Namespace
metadata:
  name: dest-ns`
	nsAppName := "dest-ns-app"
	pkgi := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: pkg0-install.test.carvel.dev.0.0.0
  namespace: dest-ns
spec:
  packageRef:
    refName: pkg0.test.carvel.dev
    versionSelection:
      constraints: 0.0.0`

	appName1 := "pkgr-1"
	appName2 := "pkgr-2"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName1})
		kapp.RunWithOpts([]string{"delete", "-a", appName2, "-n", pkgr2NS}, e2e.RunOpts{NoNamespace: true})
		kapp.Run([]string{"delete", "-a", nsAppName})
		kubectl.RunWithOpts([]string{"delete", "pkgi", "pkg0-install.test.carvel.dev.0.0.0"}, e2e.RunOpts{AllowError: true})
	}
	cleanUp()
	defer cleanUp()

	assertPkgOwnedBy := func(pkgrName string, ns string) {
		out, _ := kubectl.RunWithOpts(
			[]string{"get", "package", pkgName, "-oyaml", "-n", ns},
			e2e.RunOpts{NoNamespace: true})

		expectedOwnership := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s", ns, pkgrName)
		assert.Contains(t, out, expectedOwnership,
			"\n========\n Package Ownership Check Failed: expected ", expectedOwnership, "\n=======")
	}

	logger.Section("deploy pkgr1", func() {
		out, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1, "-n", pkgr1NS},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1), NoNamespace: true, OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out, err = kubectl.RunWithOpts([]string{"get", "packages", "-A"}, e2e.RunOpts{NoNamespace: true})
		require.NoError(t, err)
		require.Contains(t, out, "pkg0.test.carvel.dev.0.0.0")

		out, err = kapp.RunWithOpts([]string{"inspect", "-a", appName1, "--json", "-n", pkgr1NS}, e2e.RunOpts{NoNamespace: true})
		require.NoError(t, err)
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name, pkgr1NS)
	})

	logger.Section("deploy the pkgi app but it fails bc namespace restrictions", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", nsAppName}, e2e.RunOpts{StdinReader: strings.NewReader(nsApp)})
		require.NoError(t, err)
		kubectl.RunWithOpts([]string{"apply", "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgi), NoNamespace: true})
		out, _ := kubectl.RunWithOpts([]string{"get", "pkgi", "-n", "dest-ns"}, e2e.RunOpts{NoNamespace: true})
		// annoyingly we need to retry this a little...
		for i := 1; i < 5; i++ {
			if strings.Contains(out, "Reconcile failed") {
				break
			}
			time.Sleep(time.Duration(i) * time.Millisecond)
			out, _ = kubectl.RunWithOpts([]string{"get", "pkgi", "-n", "dest-ns"}, e2e.RunOpts{NoNamespace: true})
		}
		assert.Contains(t, out, "Reconcile failed: Expected to find at least one version")
		assert.NotContains(t, out, "Expected at least one fetch option") // it shouldn't contain this error message until the pkg is in an ns we can see.
	})

	logger.Section("deploy pkgr2 successfully, but pkg is still owned by pkgr1", func() {
		out, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName2, "-n", pkgr2NS},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr2), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}, NoNamespace: true})

		out, err = kapp.RunWithOpts([]string{"inspect", "-a", appName2, "--json", "-n", pkgr2NS}, e2e.RunOpts{NoNamespace: true})
		assert.NoError(t, err)
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name, pkgr1NS)
	})

	logger.Section("deploy the pkgi app but it fails after trying to install an empty package", func() {
		out, _ := kubectl.RunWithOpts([]string{"apply", "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgi), NoNamespace: true})
		out, _ = kubectl.RunWithOpts([]string{"get", "pkgi", "-n", "dest-ns", "-oyaml"}, e2e.RunOpts{NoNamespace: true})
		// in my observation this second one hasn't actually needed a retry, but i do hate flaky tests as I hate all montagues...
		for i := 1; i < 5; i++ {
			if strings.Contains(out, "Reconcile failed") {
				break
			}
			time.Sleep(time.Duration(i) * time.Millisecond)
			out, _ = kubectl.RunWithOpts([]string{"get", "pkgi", "-n", "dest-ns"}, e2e.RunOpts{NoNamespace: true})
		}

		assert.Contains(t, out, "see .status.usefulErrorMessage for details")
		assert.Contains(t, out, "Expected at least one fetch option")
	})
}

func Test_PackageReposWithOverlappingPackages_packagesHaveDifferentRevisions(t *testing.T) {
	pkgName := "pkg0.test.carvel.dev.0.0.0"
	pkgrTemplate := `
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: %s.tankyu.carvel.dev
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
	pkgr1Name := "repo1"
	pkgr2Name := "repo2"
	pkgr1 := fmt.Sprintf(pkgrTemplate, pkgr1Name, pkgName, "1")
	pkgr2 := fmt.Sprintf(pkgrTemplate, pkgr2Name, pkgName, "2")

	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	appName1 := "pkgr-1"
	appName2 := "pkgr-2"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName1})
		kapp.Run([]string{"delete", "-a", appName2})
	}
	cleanUp()
	defer cleanUp()

	assertPkgOwnedBy := func(pkgrName string) {
		out := kubectl.Run([]string{"get", "package", pkgName, "-oyaml"})
		expectedOwnership := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s", env.Namespace, pkgrName)
		assert.Contains(t, out, expectedOwnership,
			"\n========\n Package Ownership Check Failed: expected ", expectedOwnership, "\n=======")
	}

	logger.Section("deploy pkgr1", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out, _ := kubectl.RunWithOpts([]string{"get", "packages", "-A"}, e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, "pkg0.test.carvel.dev.0.0.0")

		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name)
	})

	logger.Section("deploy pkgr2 successfully, and it overrides bc it has higher rev", func() {
		out, _ := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName2},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr2), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out = kapp.Run([]string{"inspect", "-a", appName2, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr2Name)
	})

	logger.Section("uninstall and reinstall pkgr1 but it never takes ownership of the pkg bc it has lower rev", func() {
		kapp.Run([]string{"delete", "-a", appName1})
		out, _ := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr2Name)
	})

	logger.Section("install some pkgrs with some other versions", func() {
		pkgr3Name := "repo3"
		pkgr3 := fmt.Sprintf(pkgrTemplate, pkgr3Name, pkgName, "2.0") // 2.0 should take priorit over "2"
		appName3 := "pkgr-3"
		cleaner3 := func() {
			kapp.Run([]string{"delete", "-a", appName3})
		}
		cleaner3()
		defer cleaner3()
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName3},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr3), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})
		out := kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr3Name)

		/// now install one of lower rev
		pkgr4Name := "repo4"
		pkgr4 := fmt.Sprintf(pkgrTemplate, pkgr4Name, pkgName, "1.6.8") // 2.0 should still take priority
		appName4 := "pkgr-4"
		cleaner4 := func() {
			kapp.Run([]string{"delete", "-a", appName4})
		}
		cleaner4()
		defer cleaner4()
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName4},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr4), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})
		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr3Name)

		/// now install one of higher rev
		pkgr5Name := "repo5"
		pkgr5 := fmt.Sprintf(pkgrTemplate, pkgr5Name, pkgName, "2.1.0") // 2.1.0 should take priority
		appName5 := "pkgr-5"
		cleaner5 := func() {
			kapp.Run([]string{"delete", "-a", appName5})
		}
		cleaner5()
		defer cleaner5()
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName5},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr5), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})
		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr5Name)
	})

}

func Test_PackageReposWithOverlappingPackages_NonTrivialPackages(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

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

	assertPkgOwnedBy := func(pkgrName string, pkgName string) {
		out := kubectl.Run([]string{"get", "package", pkgName, "-oyaml"})
		expectedOwnership := fmt.Sprintf("packaging.carvel.dev/package-repository-ref: %s/%s", env.Namespace, pkgrName)
		assert.Contains(t, out, expectedOwnership,
			"\n========\n Package Ownership Check Failed: expected ", expectedOwnership, "\n=======")
	}

	appName1 := "pkgrapp1"
	pkgr1Name := "repo-1.tankyu.carvel.dev"

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName1})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy pkgr1", func() {
		out, _ := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName1},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr1), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out, _ = kubectl.RunWithOpts([]string{"get", "packages", "-A"}, e2e.RunOpts{NoNamespace: true})
		require.Contains(t, out, "shirt-mgr.co.uk.5.5.5")
		require.Contains(t, out, "shirt-mgr.co.uk.5.6.0")
		require.Contains(t, out, "coredino.co.uk.32.76.7")

		out = kapp.Run([]string{"inspect", "-a", appName1, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name, "shirt-mgr.co.uk.5.5.5")
		assertPkgOwnedBy(pkgr1Name, "shirt-mgr.co.uk.5.6.0")
		assertPkgOwnedBy(pkgr1Name, "coredino.co.uk.32.76.7")
	})

	pkgr2 := fmt.Sprintf(
		"%s%s%s",
		fmt.Sprintf(pkgrPreamble, 2),
		fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.5.5"),
		fmt.Sprintf(pkgTemplate, "contooor.co.uk", "0.22.0"))
	appName2 := "pkgrapp2"
	pkgr2Name := "repo-2.tankyu.carvel.dev"

	cleanUp2 := func() {
		kapp.Run([]string{"delete", "-a", appName2})
	}
	cleanUp2()
	defer cleanUp2()

	logger.Section("deploy and check pkgr2", func() {
		out, _ := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName2},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr2), OnErrKubectl: []string{"get", "pkgr", "-A", "-oyaml"}})

		out = kapp.Run([]string{"inspect", "-a", appName2, "--json"})
		require.Contains(t, out, `"reconcile_state": "ok"`)

		assertPkgOwnedBy(pkgr1Name, "shirt-mgr.co.uk.5.5.5")
		assertPkgOwnedBy(pkgr2Name, "contooor.co.uk.0.22.0")
	})

	// pkgr3 should fail to reconcile because the contooor package will conflict in the spec
	pkgr3 := fmt.Sprintf(
		"%s%s%s",
		fmt.Sprintf(pkgrPreamble, 3),
		fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.5.5"),
		strings.Replace(fmt.Sprintf(pkgTemplate, "contooor.co.uk", "0.22.0"), "k8slt/kctrl-example-pkg:v1.0.0", "k8slt/some-other-image:latest", -1))
	appName3 := "pkgrapp3"
	pkgr3Name := "repo-3.tankyu.carvel.dev"

	cleanUp3 := func() {
		kapp.Run([]string{"delete", "-a", appName3})
	}
	cleanUp3()
	defer cleanUp3()

	logger.Section("deploy and check pkgr3", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName3},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr3), AllowError: true})
		assert.Error(t, err)

		out := string(kubectl.Run([]string{"get", "pkgr", pkgr3Name, "-oyaml"}))
		assert.Contains(t, out, "Conflicting resources: Package/contooor.co.uk.0.22.0 is already present but not identical (mismatch in spec.template)")

		out = kapp.Run([]string{"inspect", "-a", appName3, "--json"})
		assert.Contains(t, out, `"reconcile_state": "fail"`)

		assertPkgOwnedBy(pkgr1Name, "shirt-mgr.co.uk.5.5.5")
		assertPkgOwnedBy(pkgr2Name, "contooor.co.uk.0.22.0")
	})

	// pkgr4 will fail because of PackageMetadatas conflict
	pkgr4 := fmt.Sprintf(
		"%s%s%s",
		fmt.Sprintf(pkgrPreamble, 4),
		fmt.Sprintf(pkgTemplate, "shirt-mgr.co.uk", "5.6.0"),
		fmt.Sprintf(pkgMetadataTemplate, "shirt-mgr.co.uk", "shirt manager", "now with dress shirts"),
	)
	appName4 := "pkgrapp4"
	pkgr4Name := "repo-4.tankyu.carvel.dev"

	cleanUp4 := func() {
		kapp.Run([]string{"delete", "-a", appName4})
	}
	cleanUp4()
	defer cleanUp4()

	logger.Section("deploy and check pkgr4", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName4},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgr4), AllowError: true})
		assert.Error(t, err)
		out := string(kubectl.Run([]string{"get", "pkgr", pkgr4Name, "-oyaml"}))
		assert.Contains(t, out, " Conflicting resources: PackageMetadata/shirt-mgr.co.uk is already present but not identical (mismatch in spec.shortDescription)")

		out = kapp.Run([]string{"inspect", "-a", appName4, "--json"})
		assert.Contains(t, out, `"reconcile_state": "fail"`)
	})
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
