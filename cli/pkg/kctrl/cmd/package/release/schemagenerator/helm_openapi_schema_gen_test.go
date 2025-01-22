// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package schemagenerator_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/release/schemagenerator"
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
# arrKeyWithStringValues default value is stringVal1
arrKeyWithStringValues:
- stringVal1
# arrKeyWithIntValues default value is 1
arrKeyWithIntValues:
- 1
# arrKeyWithFloatValues default value is 1.1. 1.2 is ignored
arrKeyWithFloatValues:
- 1.1
- 1.2
arrKeyEmpty: []
`,
			want: `properties:
  arrKeyEmpty:
    default: []
    items: {}
    type: array
  arrKeyWithFloatValues:
    default: []
    description: default value is 1.1. 1.2 is ignored
    items:
      default: 1.1
      format: float
      type: number
    type: array
  arrKeyWithIntValues:
    default: []
    description: default value is 1
    items:
      default: 1
      type: integer
    type: array
  arrKeyWithStringValues:
    default: []
    description: default value is stringVal1
    items:
      default: stringVal1
      type: string
    type: array
type: object
`,
		},
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
# Object example
objExample: {}
`,
			want: `properties:
  arrExample:
    default: []
    description: Array example
    items:
      default: arr1
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
  objExample:
    default: {}
    description: Object example
    type: object
type: object
`,
		},
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
  env2:
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
          type: object
        type: array
      env2:
        default: []
        items:
          properties:
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
`,
		},
		{
			name: "Alias Node",
			input: `
# alias Example
aliasExample: &aliasEx
  # foo bar example
  foo: bar

# global Example
global:
  - *aliasEx

# alaisex as key, not an aliasNode here
aliasEx:
  # foo1bar1 Example
  foo1: bar1
`,
			want: `properties:
  aliasEx:
    description: alaisex as key, not an aliasNode here
    properties:
      foo1:
        default: bar1
        description: bar1 Example
        type: string
    type: object
  aliasExample:
    description: alias Example
    properties:
      foo:
        default: bar
        description: bar example
        type: string
    type: object
  global:
    default: []
    description: Example
    items:
      properties:
        foo:
          default: bar
          description: bar example
          type: string
      type: object
    type: array
type: object
`,
		},
		{
			name: "unknown type",
			input: `
# a field without a type
anything: null
`,
			want: `properties:
  anything:
    description: a field without a type
    oneOf:
    - default: null
      nullable: true
      type: integer
    - default: null
      nullable: true
      type: number
    - default: null
      nullable: true
      type: boolean
    - default: null
      nullable: true
      type: string
    - default: null
      nullable: true
      type: object
    - default: null
      items: {}
      nullable: true
      type: array
type: object
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup()
			defer cleanup()
			dirName := "tmp"
			err := os.Mkdir(dirName, fs.ModePerm)
			require.NoError(t, err)
			fileName := "values.yaml"
			err = os.WriteFile(filepath.Join(dirName, fileName), []byte(test.input), fs.ModePerm)
			require.NoError(t, err)
			valuesSchema, err := schemagenerator.NewHelmValuesSchemaGen("tmp").Schema()
			output, err := yaml.JSONToYAML(valuesSchema.OpenAPIv3.Raw)
			require.NoError(t, err)
			require.Equal(t, test.want, string(output), "Expected valuesSchema to match")
		})
	}
}

func TestHelmValuesSchemaGen_Schema_EmptyFile(t *testing.T) {
	cleanup()
	defer cleanup()
	dirName := "tmp"
	err := os.Mkdir(dirName, fs.ModePerm)
	require.NoError(t, err)
	fileName := "values.yaml"
	err = os.WriteFile(filepath.Join(dirName, fileName), []byte(""), fs.ModePerm)
	require.NoError(t, err)
	valuesSchema, err := schemagenerator.NewHelmValuesSchemaGen("tmp").Schema()
	require.NoError(t, err)
	require.Equal(t, 0, len(valuesSchema.OpenAPIv3.Raw), "Expected valuesSchema.OpenAPIv3.Raw to be empty")
}

func TestHelmValuesSchemaGen_Schema_File_Not_Present(t *testing.T) {
	cleanup()
	defer cleanup()
	cleanup()
	dirName := "tmp"
	err := os.Mkdir(dirName, fs.ModePerm)
	require.NoError(t, err)
	valuesSchema, err := schemagenerator.NewHelmValuesSchemaGen("tmp").Schema()
	require.NoError(t, err)
	require.Equal(t, 0, len(valuesSchema.OpenAPIv3.Raw), "Expected valuesSchema.OpenAPIv3.Raw to be empty")
}
