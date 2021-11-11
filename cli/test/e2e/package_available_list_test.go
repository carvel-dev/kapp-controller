// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageAvailableList(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, env.KappBinaryPath, logger}
	kappCtrl := Kapp{t, env.Namespace, env.KappCtrlBinaryPath, logger}

	appName := "test-package-name"

	packageMetadataName := "test-pkg.carvel.dev"

	packageMetadata := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: %s
spec:
  displayName: "Carvel Test Package"
  shortDescription: "Carvel package for testing installation"`, packageMetadataName)

	package1Name := "test-pkg.carvel.dev.1.0.0"
	package1Version := "1.0.0"

	package2Name := "test-pkg.carvel.dev.1.1.0"
	package2Version := "1.1.0"

	packageCR := `---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: %s
spec:
  refName: test-pkg.carvel.dev
  version: %s
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

	packageCR1 := fmt.Sprintf(packageCR, package1Name, package1Version)
	packageCR2 := fmt.Sprintf(packageCR, package2Name, package2Version)
	yaml := fmt.Sprintf("%s\n%s\n%s", packageMetadata, packageCR1, packageCR2)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	defer cleanUp()

	logger.Section("package available list with no package present", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list"}, RunOpts{})
		require.NoError(t, err)

		out = strings.TrimSpace(replaceSpaces(replaceTarget(out)))

		expectedOutput := strings.TrimSpace(replaceSpaces(`

Available packages in namespace 'kapp-test'

Name  Display-Name  Short-Description  

0 Packages Available

Succeeded`))

		require.Equal(t, expectedOutput, out)

	})

	logger.Section("Adding test package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-a", appName, "-f", "-"}, RunOpts{
			StdinReader: strings.NewReader(yaml), AllowError: true,
		})
		require.NoError(t, err)
	})

	logger.Section("package available list with one package available", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list"}, RunOpts{})
		require.NoError(t, err)

		out = strings.TrimSpace(replaceSpaces(replaceTarget(out)))

		expectedOutput := strings.TrimSpace(replaceSpaces(`

Available packages in namespace 'kapp-test'

Name                 Display-Name         Short-Description  
test-pkg.carvel.dev  Carvel Test Package  Carvel package for testing installation  

1 Packages Available

Succeeded`))

		require.Equal(t, expectedOutput, out)
	})

	logger.Section("package available list versions of a package", func() {
		out, err := kappCtrl.RunWithOpts([]string{"package", "available", "list", "-p", "test-pkg.carvel.dev"}, RunOpts{})
		require.NoError(t, err)

		out = strings.TrimSpace(replaceSpaces(replaceTarget(out)))

		expectedOutput := strings.TrimSpace(replaceSpaces(`

Available package versions for 'test-pkg.carvel.dev' in namespace 'kapp-test'

Name                 Version  Released-At  
test-pkg.carvel.dev  1.0.0    0001-01-01 00:00:00 +0000 UTC  
test-pkg.carvel.dev  1.1.0    0001-01-01 00:00:00 +0000 UTC  

2 Package Versions Available

Succeeded`))

		require.Equal(t, expectedOutput, out)
	})
}

func replaceSpaces(result string) string {
	// result = strings.Replace(result, " ", "_", -1) // useful for debugging
	result = strings.Replace(result, " \n", " $\n", -1) // explicit endline
	return result
}

func replaceTarget(result string) string {
	return regexp.MustCompile("Target cluster .+\n").ReplaceAllString(result, "")
}
