// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/installpackage/v1alpha1"
)

func Test_PackageVersionWithValuesSchema_PreservesSchemaData(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t: t, namespace: env.Namespace, l: logger}
	name := "pkg-with-schema.1.0.0"

	pkgYaml := fmt.Sprintf(`---
apiVersion: package.carvel.dev/v1alpha1
kind: PackageVersion
metadata:
  name: %s
spec:
  version: 1.0.0
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
      - kapp: {}`, name)

	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(pkgYaml)})

	out := kubectl.Run([]string{"get", "packageversions/" + name, "-o=jsonpath={.spec.valuesSchema.openAPIv3}"})
	if !strings.Contains(out, "properties") && !strings.Contains(out, "hello_msg") && !strings.Contains(out, "svc_port") {
		t.Fatalf("Could not find properties on values schema. Got:\n%s", out)
	}

	out = kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind=PackageVersion"})
	var cr v1alpha1.InternalPackageVersion
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
