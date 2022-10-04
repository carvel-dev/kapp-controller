// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/release"
	"sigs.k8s.io/yaml"
)

func cleanup() {
	os.RemoveAll("tmp")
}

func TestHelmValuesSchemaGen_Schema_Success(t *testing.T) {
	type test struct {
		name  string
		input string
		want  string
	}

	tests := []test{
		{
			name: "array with different values",
			input: `
# arrKeyWithStringValues default value is stringVal1, stringVal2
arrKeyWithStringValues:
- stringVal1
- stringVal2
# arrKeyWithIntValues default value is 1, 2
arrKeyWithIntValues:
- 1
- 2
# arrKeyWithFloatValues default value is 1.1, 1.2
arrKeyWithFloatValues:
- 1.1
- 1.2
`,
			want: `properties:
  arrKeyWithFloatValues:
    default:
    - 1.1
    - 1.2
    description: default value is 1.1, 1.2
    items:
      type: number
    type: array
  arrKeyWithIntValues:
    default:
    - 1
    - 2
    description: default value is 1, 2
    items:
      type: integer
    type: array
  arrKeyWithStringValues:
    default:
    - stringVal1
    - stringVal2
    description: default value is stringVal1, stringVal2
    items:
      type: string
    type: array
type: object
`},
		{
			name: "object with different values",
			input: `
# Image details
image: my-docker-image:1.0.0
# Container name
name: test-container
# Boolean example
boolExample: true
# Array example
arrExample:
- arr1
# Float example
floatExample: 2.3
# Integer example
intExample: 3
`,
			want: `properties:
  arrExample:
    default:
    - arr1
    description: Array example
    items:
      type: string
    type: array
  boolExample:
    default: true
    description: Boolean example
    type: boolean
  floatExample:
    default: 2.3
    description: Float example
    format: float
    type: number
  image:
    default: my-docker-image:1.0.0
    description: Image details
    type: string
  intExample:
    default: 3
    description: Integer example
    type: integer
  name:
    default: test-container
    description: Container name
    type: string
type: object
`},
		{
			name: "nested complex object",
			input: `
containers:
  # Image details
  image: my-docker-image:1.0.0
  # Array example
  env:
  # key1 example
  - key1: val1
  -
  # key2 example
    key2: val2
  # Float example
  floatExample: 2.3
`,
			want: `properties:
  containers:
    properties:
      env:
        default: []
        description: Array example
        items:
          properties:
            key1:
              default: val1
              description: example
              type: string
            key2:
              default: val2
              description: example
              type: string
          type: object
        type: array
      floatExample:
        default: 2.3
        description: Float example
        format: float
        type: number
      image:
        default: my-docker-image:1.0.0
        description: Image details
        type: string
    type: object
type: object
`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer cleanup()
			dirName := "tmp"
			err := createDir(dirName)
			require.NoError(t, err)
			fileName := "values.yaml"
			err = createFile(filepath.Join(dirName, fileName), []byte(test.input))
			require.NoError(t, err)
			valuesSchema, err := release.NewHelmValuesSchemaGen("tmp").Schema()
			output, err := yaml.JSONToYAML(valuesSchema.OpenAPIv3.Raw)
			require.Equal(t, test.want, string(output), "Expected valuesSchema to match")
		})
	}
}

func TestHelmValuesSchemaGen_Schema_EmptyFile(t *testing.T) {
	defer cleanup()
	dirName := "tmp"
	err := createDir(dirName)
	require.NoError(t, err)
	fileName := "values.yaml"
	err = createFile(filepath.Join(dirName, fileName), []byte(""))
	require.NoError(t, err)
	valuesSchema, err := release.NewHelmValuesSchemaGen("tmp").Schema()
	require.Equal(t, 0, len(valuesSchema.OpenAPIv3.Raw), "Expected valuesSchema.OpenAPIv3.Raw to be empty")
}

func TestHelmValuesSchemaGen_Schema_File_Not_Present(t *testing.T) {
	defer cleanup()
	cleanup()
	dirName := "tmp"
	err := createDir(dirName)
	require.NoError(t, err)
	valuesSchema, err := release.NewHelmValuesSchemaGen("tmp").Schema()
	require.Equal(t, 0, len(valuesSchema.OpenAPIv3.Raw), "Expected valuesSchema.OpenAPIv3.Raw to be empty")
}

func createDir(dirName string) error {
	err := os.Mkdir(dirName, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create %s directory: %s", dirName, err.Error())
	}
	return nil
}

func createFile(filePath string, fileData []byte) error {
	err := os.WriteFile(filePath, fileData, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create %s file: %s", fileData, err.Error())
	}
	return nil
}
