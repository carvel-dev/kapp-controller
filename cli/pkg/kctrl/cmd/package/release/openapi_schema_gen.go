// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"bytes"
	"fmt"
	goexec "os/exec"

	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type ValuesSchemaGen struct {
	paths []string
}

func NewValuesSchemaGen(paths []string) ValuesSchemaGen {
	return ValuesSchemaGen{paths}
}

func (g ValuesSchemaGen) Schema() (*kcdatav1alpha1.ValuesSchema, error) {
	cmdArgs := []string{"--data-values-schema-inspect", "--output=openapi-v3"}
	for _, path := range g.paths {
		cmdArgs = append(cmdArgs, "-f", path)
	}

	var stderrBuf, stdoutBuf bytes.Buffer
	cmd := goexec.Command("ytt", cmdArgs...)
	cmd.Stderr = &stderrBuf
	cmd.Stdout = &stdoutBuf

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, stderrBuf.String())
	}

	jsonEncodedBytes, err := g.stdoutToJSONBytes(stdoutBuf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Parsing ytt open api output: %s", err)
	}

	return &kcdatav1alpha1.ValuesSchema{
		OpenAPIv3: runtime.RawExtension{Raw: jsonEncodedBytes},
	}, nil
}

type partialYttSchemaOutput struct {
	Components struct {
		Schemas struct {
			DataValues interface{} `yaml:"dataValues"`
		} `yaml:"schemas"`
	} `yaml:"components"`
}

func (g ValuesSchemaGen) stdoutToJSONBytes(stdout []byte) ([]byte, error) {
	// Extract the part of YAML needed by the package
	var partialYttSchemaOutput partialYttSchemaOutput
	err := yaml.Unmarshal(stdout, &partialYttSchemaOutput)
	if err != nil {
		return nil, err
	}
	partialYAMLBytes, err := yaml.Marshal(partialYttSchemaOutput.Components.Schemas.DataValues)
	if err != nil {
		return nil, err
	}

	jsonEncodedBytes, err := yaml.YAMLToJSON(partialYAMLBytes)
	if err != nil {
		return nil, err
	}

	return jsonEncodedBytes, nil
}
