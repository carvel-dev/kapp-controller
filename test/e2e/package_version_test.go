// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
)

func Test_PackageVersionIsValidated_Name(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	appName := "invalid-pkg-version-name-test"

	invalidPackageVersionName := "notThePackage-notTheVersion"

	invalidPkgVersionYML := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: %s
spec:
  packageName: test-pkg.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt:
          paths:
          - "config.yml"
          - "values.yml"
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}`, invalidPackageVersionName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	defer cleanUp()

	logger.Section("deploy package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName},
			RunOpts{StdinReader: strings.NewReader(invalidPkgVersionYML), AllowError: true})

		if err == nil {
			t.Fatalf("Expected package version creation to fail")
		}

		if !strings.Contains(err.Error(), "is invalid: metadata.name") {
			t.Fatalf("Expected package version creation error to contain message about invalid name, but got: %v", err)
		}

		if !strings.Contains(err.Error(), "must begin with <spec.packageName> + '.'") {
			t.Fatalf("Expected error message to contain required form for package version name, got: %v", err)
		}
	})
}

func Test_PackageVersionWithValuesSchema_PreservesSchemaData(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t: t, namespace: env.Namespace, l: logger}
	appName := "test-package-version-schema"
	packageName := "pkg.test.carvel.dev.1.0.0"
	version := "1.0.0"

	pkgYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: %s.%s
spec:
  packageName: %s
  version: %s
  valuesSchema:
    openAPIv3:
      properties:
        svc_port:
          description: Port number for service. Defaults to 80.
          type: int
        hello_msg:
          description: The message simple-app will display
          type: string
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt:
          paths:
          - "config.yml"
          - "values.yml"
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}`, packageName, version, packageName, version)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	cleanUp()
	defer cleanUp()

	kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(pkgYaml)})

	out := kubectl.Run([]string{"get", "packageversions/" + fmt.Sprintf("%s.%s", packageName, version), "-o=jsonpath={.spec.valuesSchema.openAPIv3}"})
	if !strings.Contains(out, "properties") && !strings.Contains(out, "hello_msg") && !strings.Contains(out, "svc_port") {
		t.Fatalf("Could not find properties on values schema. Got:\n%s", out)
	}

	out = kapp.Run([]string{"inspect", "-a", appName, "--raw", "--tty=false", "--filter-kind=PackageVersion"})
	var cr v1alpha1.PackageVersion
	err := yaml.Unmarshal([]byte(out), &cr)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	var into interface{}
	err = yaml.Unmarshal(cr.Spec.ValuesSchema.OpenAPIv3.Raw, &into)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
}

func Test_PackageVersion_FieldSelectors(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	name := "test-package-version-field-selector"
	packageName := "test-package.carvel.dev"
	filteredPackageName := "you-shouldnt-see-me.carvel.dev"
	packcageYamls := fmt.Sprintf(`---
kind: Package
apiVersion: data.packaging.carvel.dev/v1alpha1
metadata:
  name: %[1]s
spec:
  shortDescription: "Package for testing"
---
kind: PackageVersion
apiVersion: data.packaging.carvel.dev/v1alpha1
metadata:
  name: %[1]s.1.0.0
spec:
  packageName: %[1]s
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
                name: configmap
              data:
                hello_msg: hi
      template:
      - ytt: {}
      deploy:
      - kapp: {}
---
kind: Package
apiVersion: data.packaging.carvel.dev/v1alpha1
metadata:
  name: %[2]s
spec:
  shortDescription: "Package for testing"
---
kind: PackageVersion
apiVersion: data.packaging.carvel.dev/v1alpha1
metadata:
  name: %[2]s.1.0.0
spec:
  packageName: %[2]s
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
                name: configmap
              data:
                hello_msg: hi
      template:
      - ytt: {}
      deploy:
      - kapp: {}
`, packageName, filteredPackageName)

	cleanup := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	defer cleanup()

	logger.Section("deploy package and package version", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{StdinReader: strings.NewReader(packcageYamls)})
	})

	logger.Section("check field selector", func() {
		out, err := kubectl.RunWithOpts([]string{"get", "packageversions",
			"--field-selector", fmt.Sprintf("spec.packageName=%s", packageName)}, RunOpts{AllowError: true})

		if err != nil {
			t.Fatalf("Expected field selector to successfully return a package version but got error: %v", err)
		}

		if strings.Contains(out, filteredPackageName) {
			t.Fatalf("Expected not to see filtered package in output:\n %s", out)
		}
	})
}

func TestOverridePackageVersionDelete(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	k := Kubectl{t, env.Namespace, logger}

	localNS := env.Namespace
	globalNS := env.PackagingGlobalNS

	packagesYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: pkg1.test.carvel.dev.1.0.0-global
  namespace: %s
spec:
  packageName: pkg1.test.carvel.dev
  version: 1.0.0
  template:
    spec: {}
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: pkg1.test.carvel.dev.1.0.0
  namespace: %s
spec:
  packageName: pkg1.test.carvel.dev
  version: 1.0.0
  template:
    spec: {}`, globalNS, localNS)

	cleanup := func() {
		k.RunWithOpts([]string{"delete", "packageversions/pkg1.test.carvel.dev.1.0.0-global", "-n", globalNS}, RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", "packageversions/pkg1.test.carvel.dev.1.0.0", "-n", localNS}, RunOpts{NoNamespace: true, AllowError: true})
	}
	defer cleanup()

	logger.Section("cleanup", cleanup)

	logger.Section("deploy packages", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(packagesYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected package version application to succeed, but got: %v", err)
		}
	})

	logger.Section("attempt to delete the local package", func() {
		timeout := 30 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		_, err := k.RunWithOpts([]string{"delete", "packageversions/pkg1.test.carvel.dev.1.0.0", "-n", localNS}, RunOpts{Ctx: ctx, NoNamespace: true, AllowError: true})
		if err != nil {
			t.Fatalf("Expected delete of local package version to succeed in %v, but got: %v", timeout, err)
		}
	})
}

func TestOverridePackageVersionNamespaceDelete(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	k := Kubectl{t, env.Namespace, logger}

	localNS := "test-ns"
	globalNS := env.PackagingGlobalNS

	packagesYaml := fmt.Sprintf(`---
apiVersion: v1
kind: Namespace
metadata:
  name: %[1]s
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: pkg1.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  packageName: pkg1.test.carvel.dev
  version: 1.0.0
  template:
    spec: {}
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: pkg1.test.carvel.dev.1.0.0-global
  namespace: %[2]s
spec:
  packageName: pkg1.test.carvel.dev
  version: 1.0.0
  template:
    spec: {}`, localNS, globalNS)

	cleanup := func() {
		k.RunWithOpts([]string{"delete", "packageversions/pkg1.test.carvel.dev.1.0.0-global", "-n", globalNS}, RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", "packageversions/pkg1.test.carvel.dev.1.0.0", "-n", localNS}, RunOpts{NoNamespace: true, AllowError: true})
		k.RunWithOpts([]string{"delete", fmt.Sprintf("namespaces/%s", localNS)}, RunOpts{NoNamespace: true, AllowError: true})
	}
	defer logger.Section("post test cleanup", cleanup)

	logger.Section("pre test cleanup", cleanup)

	logger.Section("deploy packages and namespace", func() {
		_, err := k.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(packagesYaml), NoNamespace: true})
		if err != nil {
			t.Fatalf("Expected package version application to succeed, but got: %v", err)
		}
	})

	logger.Section("attempt to delete the local namespace", func() {
		timeout := 30 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		_, err := k.RunWithOpts([]string{"delete", fmt.Sprintf("namespaces/%s", localNS)}, RunOpts{Ctx: ctx, NoNamespace: true, AllowError: true})
		if err != nil {
			if strings.Contains(err.Error(), "signal: interrupt") {
				t.Fatalf("Timed out waiting for delete of namespace '%s'", localNS)
			}
			t.Fatalf("Expected delete of local namespace '%s' to succeed, but got: %v", localNS, err)
		}

		_, err = k.RunWithOpts([]string{"get", fmt.Sprintf("namespaces/%s", localNS)}, RunOpts{NoNamespace: true, AllowError: true})
		if err == nil {
			t.Fatalf("Expected not to find local namespace '%s', but did", localNS)
		}
	})
}
