// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	"github.com/stretchr/testify/require"
)

func TestAPIServerProtobuf(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}

	config := `
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: protobuf-test.carvel.dev.2.0.0
spec:
  refName: protobuf-test.carvel.dev
  version: 2.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg
          tagSelection:
            semver:
              constraints: "2.0.0 <98765.0.0"
      template:
      - ytt:
          paths:
          - config-step-2-template
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}
`

	customNs := "test-kappctrl-protobuf"
	cleanUp := func() {
		kubectl.Run([]string{"delete", "ns", customNs, "--ignore-not-found"})
	}

	cleanUp()
	defer cleanUp()

	// From @liggitt (in #sig-api-machinery Slack channel):
	//   protobuf tags are not used on CRDs (kube-apiserver only supports
	//   json/yaml requests/responses and json storage for CRD-based APIs).
	//   a custom aggregated server can support protobuf requests/responses/storage;
	//   not supporting protobuf serialization on those types can cause:
	//   https://github.com/kubernetes/kubernetes/issues/86666
	logger.Section("create package in a new namespace", func() {
		kubectl.Run([]string{"create", "ns", customNs})

		kubectl.RunWithOpts([]string{"apply", "-f", "-", "-n", customNs},
			e2e.RunOpts{StdinReader: strings.NewReader(config), NoNamespace: true})

		kubectl.RunWithOpts([]string{"get", "packages", "-n", customNs}, e2e.RunOpts{NoNamespace: true})

		out, _ := kubectl.RunWithOpts([]string{"get", "packages", "protobuf-test.carvel.dev.2.0.0", "-o", "yaml", "-n", customNs},
			e2e.RunOpts{NoNamespace: true})

		// Check that semver constraints are returned
		// (These fields are represented by vendir's versions package)
		// (cheap check that content is returned via string contains vs unmarshaling)
		require.Equal(t, 1, strings.Count(out, `constraints: 2.0.0 <98765.0.0`))
	})

	logger.Section("delete succeeds", func() {
		kubectl.Run([]string{"delete", "ns", customNs})
	})
}
