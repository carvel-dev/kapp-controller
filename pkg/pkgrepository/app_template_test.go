// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository_test

import (
	"testing"

	"carvel.dev/kapp-controller/pkg/pkgrepository"
	"github.com/stretchr/testify/assert"
)

func TestFilterResourcesYAMLUnique(t *testing.T) {
	input := `
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg-standalone.test.carvel.dev
  namespace: pkg-standalone
spec:
  displayName: "Test Package standalone"
  shortDescription: "Package used for testing"
  unexpectedKeyNotInSpec: "this is just a test"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg-standalone.test.carvel.dev.2.0.0
  namespace: pkg-standalone
spec:
  refName: pkg-standalone.test.carvel.dev
  version: 2.0.0
  releasedAt: 2021-05-05T18:57:06Z
  template: # type of App CR
    spec:
      fetch:
      - git:
          url: https://github.com/k14s/k8s-simple-app-example
          ref: origin/develop
          unexpectedSubKeyNotInSpec: hahaha
      template:
      - ytt:
          paths:
          - config-step-2-template
          - config-step-2a-overlays
      unexpectedKeyNotInSpec: foo
      deploy:
      - kapp: {}`

	observedOutput, err := pkgrepository.FilterResources(input)
	assert.NoError(t, err)
	assert.Contains(t, observedOutput, "kind: PackageMetadata\n")
	assert.Contains(t, observedOutput, "kind: Package\n")
	assert.NotContains(t, observedOutput, "unexpected") // does not include unexpected* keys
}
