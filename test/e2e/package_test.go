// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
)

func Test_PackageIsValidated(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	appName := "invalid-pkg-name-test"

	invalidPackageName := "I am invalid"

	invalidPkgYML := fmt.Sprintf(`---
apiVersion: package.carvel.dev/v1alpha1
kind: Package
metadata:
  name: %s
spec:
  shortDescription: I am invalid
`, invalidPackageName)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", appName})
	}
	defer cleanUp()

	logger.Section("deploy package", func() {
		_, err := kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", appName},
			RunOpts{StdinReader: strings.NewReader(invalidPkgYML), AllowError: true})

		if err == nil {
			t.Fatalf("Expected package creation to fail")
		}

		if !strings.Contains(err.Error(), "is invalid: metadata.name") {
			t.Fatalf("Expected package creation error to contain message about invalid name")
		}
	})
}
