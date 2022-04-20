// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlUnique(t *testing.T) {
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

	observedOutput, err := FilterResources([]byte(input))
	assert.NoError(t, err)
	assert.NotContains(t, observedOutput, "unexpected")
}
