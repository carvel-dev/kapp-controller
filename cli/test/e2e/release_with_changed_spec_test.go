// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	packageBuildFile     = "package-build.yml"
	packageResourcesFile = "package-resources.yml"
)

func TestPackageReleaseWithChangedSpec(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kappCtrl := Kctrl{t, env.Namespace, env.KctrlBinaryPath, logger}

	cleanUp := func() {
		os.RemoveAll(workingDir)
	}
	cleanUp()
	defer cleanUp()

	refName := "test-kctrl.carvel.dev"
	releaseNotes := "Made things better"
	shortDescription := "shorter description"
	longDescription := "longer description"

	packageBuild := fmt.Sprintf(`
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
metadata:
  creationTimestamp: null
  name: %s
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
`, refName, env.Image)

	packageResources := fmt.Sprintf(`
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  creationTimestamp: null
  name: %s.0.0.0
spec:
  licenses:
  - MIT
  refName: %s
  releasedAt: null
  releaseNotes: %s
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
  name: %s
spec:
  displayName: test-kctrl
  longDescription: %s
  shortDescription: %s

---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  annotations:
    kctrl.carvel.dev/local-fetch-0: .
  creationTimestamp: null
  name: test-kctrl
spec:
  packageRef:
  refName: test-kctrl.carvel.dev
  versionSelection:
    constraints: 0.0.0
  serviceAccountName: test-kctrl-sa
status:
  conditions: null
  friendlyDescription: ""
  observedGeneration: 0
`, refName, refName, releaseNotes, refName, longDescription, shortDescription)

	configManifests := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: kctrl-test
data:
  foo: bar	
`

	// Create manifest and files created by `init` command (with changes)
	err := os.MkdirAll(filepath.Join(workingDir, "config"), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workingDir, packageResourcesFile), []byte(packageResources), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workingDir, packageBuildFile), []byte(packageBuild), os.ModePerm)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(workingDir, "config", "config.yml"), []byte(configManifests), os.ModePerm)
	require.NoError(t, err)

	logger.Section("release package", func() {
		kappCtrl.RunWithOpts([]string{"package", "release", "--chdir", workingDir}, RunOpts{NoNamespace: true})

		packagePath := "carvel-artifacts/packages"
		packageOutput, err := os.ReadFile(filepath.Join(workingDir, packagePath, refName, "package.yml"))
		require.NoError(t, err)
		require.Contains(t, string(packageOutput), fmt.Sprintf("releaseNotes: %s", releaseNotes))
		require.Contains(t, string(packageOutput), "licenses:\n  - MIT")

		packageMetadataOutput, err := os.ReadFile(filepath.Join(workingDir, packagePath, refName, "metadata.yml"))
		require.NoError(t, err)
		require.Contains(t, string(packageMetadataOutput), fmt.Sprintf("shortDescription: %s", shortDescription))
		require.Contains(t, string(packageMetadataOutput), fmt.Sprintf("longDescription: %s", longDescription))
	})
}
