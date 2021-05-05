// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
)

func Test_PackageWithoutTwoPeriodsInName_ResultsInError(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	name := "bad-package-name"

	pkgYaml := fmt.Sprintf(`---
apiVersion: package.carvel.dev/v1alpha1
kind: Package
metadata:
  name: %s
spec:
  publicName: pkg.fail.carvel.dev
  version: 1.0.0
  displayName: "Test Package in repo"
  description: "Package used for testing"
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
      - kapp: {}`, name)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	out, err := kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(pkgYaml), AllowError: true})
	if err == nil {
		t.Fatalf("expected error from Package name not passing validation but got:\n%s", out)
	}

	if !strings.Contains(err.Error(), "package name requires at least two periods. Recommended package naming convention is publicName.packageRepo.version") {
		t.Fatalf("unexpected error:\n%s", err.Error())
	}
}
