package e2e

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYttTemplateWithBuildValues(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	configYAML := `
#@ load("@ytt:data", "data")
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: pkg-cm
  namespace: #@ data.values.namespace
data:
  foo: bar
`
	schemaYAML := `
#@data/values-schema
---
#@schema/validation min_len=1
namespace: ""
`
	packageBuildYAML := fmt.Sprintf(`
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
metadata:
  creationTimestamp: null
  name: samplepackage.corp.com
spec:
  release:
  - resource: {}
  template:
    spec:
      app:
        spec:
          deploy:
          - kapp: {}
          template:
          - ytt:
              paths:
              - config
          - kbld: {}
      export:
      - imgpkgBundle:
          image: %s
          useKbldImagesLock: true
        includePaths:
        - config
`, env.Image)
	packageResourcesYAML := `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  creationTimestamp: null
  name: samplepackage.corp.com.0.0.0
spec:
  refName: samplepackage.corp.com
  releasedAt: null
  template:
    spec:
      deploy:
      - kapp: {}
      fetch:
      - git: {}
      template:
      - ytt:
        paths:
        - config
      - kbld: {}
  valuesSchema:
    openAPIv3: null
  version: 0.0.0
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  creationTimestamp: null
  name: samplepackage.corp.com
spec:
  displayName: samplepackage
  longDescription: samplepackage.corp.com
  shortDescription: samplepackage.corp.com
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  annotations:
    kctrl.carvel.dev/local-fetch-0: .
  creationTimestamp: null
  name: samplepackage
spec:
  packageRef:
  refName: samplepackage.corp.com
  versionSelection:
    constraints: 0.0.0
  serviceAccountName: samplepackage-sa
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`

	values := `
namespace: test
`

	cleanUp := func() {
		os.RemoveAll(workingDir)
	}
	cleanUp()
	defer cleanUp()

	configDir := "config"
	err := os.Mkdir(workingDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(path.Join(workingDir, configDir), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, configDir, "config.yaml"), []byte(configYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, configDir, "schema.yaml"), []byte(schemaYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, "package-build.yml"), []byte(packageBuildYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, "package-resources.yml"), []byte(packageResourcesYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	logger.Section("run package release without build values", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "release", "--chdir", workingDir}, RunOpts{NoNamespace: true, AllowError: true})
		require.Error(t, err)
	})

	logger.Section("run package release with build values", func() {
		err = os.WriteFile(path.Join(workingDir, "build-values.yml"), []byte(values), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		// Verify that validation checks are not performed while running ytt to build packages
		kappCtrl.RunWithOpts([]string{"package", "release", "--chdir", workingDir, "--build-values", "build-values.yml"}, RunOpts{NoNamespace: true})
	})
}

func TestHelmTemplateWithBuildValues(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	configYAML := `
{{- $fooVal := .Values.fooVal | required ".Values.fooVal is required." -}}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-cm-1
data:
  foo {{ $fooVal }}
`
	chartYAML := `
apiVersion: v1
name: test-chart
version: 1.0.0
`
	packageBuildYAML := fmt.Sprintf(`
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
metadata:
  creationTimestamp: null
  name: samplepackage.corp.com
spec:
  release:
  - resource: {}
  template:
    spec:
      app:
        spec:
          deploy:
          - kapp: {}
          template:
          - helmTemplate:
              path: config
          - kbld: {}
      export:
      - imgpkgBundle:
          image: %s
          useKbldImagesLock: true
        includePaths:
        - config
`, env.Image)
	packageResourcesYAML := `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  creationTimestamp: null
  name: samplepackage.corp.com.0.0.0
spec:
  refName: samplepackage.corp.com
  releasedAt: null
  template:
    spec:
      deploy:
      - kapp: {}
      fetch:
      - git: {}
      template:
      - helmTemplate:
        path: config
      - kbld: {}
  valuesSchema:
    openAPIv3: null
  version: 0.0.0
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  creationTimestamp: null
  name: samplepackage.corp.com
spec:
  displayName: samplepackage
  longDescription: samplepackage.corp.com
  shortDescription: samplepackage.corp.com
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  annotations:
    kctrl.carvel.dev/local-fetch-0: .
  creationTimestamp: null
  name: samplepackage
spec:
  packageRef:
  refName: samplepackage.corp.com
  versionSelection:
    constraints: 0.0.0
  serviceAccountName: samplepackage-sa
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`

	values := `
fooVal: bar
`

	cleanUp := func() {
		os.RemoveAll(workingDir)
	}
	cleanUp()
	defer cleanUp()

	configDir := "config"
	err := os.Mkdir(workingDir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	templateDir := "templates"
	err = os.Mkdir(path.Join(workingDir, configDir), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(path.Join(workingDir, configDir, templateDir), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, configDir, templateDir, "config.yaml"), []byte(configYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, configDir, "Chart.yaml"), []byte(chartYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, "package-build.yml"), []byte(packageBuildYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(workingDir, "package-resources.yml"), []byte(packageResourcesYAML), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	logger.Section("run helm templated package release without build values", func() {
		_, err := kappCtrl.RunWithOpts([]string{"package", "release", "--chdir", workingDir}, RunOpts{NoNamespace: true, AllowError: true})
		require.Error(t, err)
	})

	logger.Section("run helm templated package release with build values", func() {
		err = os.WriteFile(path.Join(workingDir, "build-values.yml"), []byte(values), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		kappCtrl.RunWithOpts([]string{"package", "release", "--chdir", workingDir, "--build-values", "build-values.yml"}, RunOpts{NoNamespace: true})
	})
}
