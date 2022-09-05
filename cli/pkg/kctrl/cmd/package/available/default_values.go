// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package available

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	structuralschema "k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
	"k8s.io/apimachinery/pkg/runtime"
)

type PackageSchema struct {
	Raw []byte
}

// DefaultValues returns a yaml byte array with values populated according to schema
func (s PackageSchema) DefaultValues() ([]byte, error) {
	jsonSchemaProps := &apiextensions.JSONSchemaProps{}

	if err := yaml.Unmarshal(s.Raw, jsonSchemaProps); err != nil {
		return nil, err
	}

	ss, err := structuralschema.NewStructural(jsonSchemaProps)
	if err != nil {
		return nil, err
	}

	unstructured := make(map[string]interface{})

	s.schemaDefault(unstructured, ss)

	var b bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)

	if err := yamlEncoder.Encode(unstructured); err != nil {
		return nil, err
	}

	if err := yamlEncoder.Close(); err != nil {
		return nil, err
	}

	return s.commentDefaultValues(b), nil
}

// schemaDefault does defaulting of x depending on default values in s.
// This is adopted from k8s.io/apiextensions-apiserver/pkg/apiserver/schema/defaulting with 2 changes
// 1. Prevent deep copy of int as it panics
// 2. For type object depth first search to see if there is any property with default
//
//gocyclo:ignore
func (s PackageSchema) schemaDefault(x interface{}, ss *structuralschema.Structural) {
	if ss == nil {
		return
	}

	switch x := x.(type) {
	case map[string]interface{}:
		for k, prop := range ss.Properties {
			if prop.Default.Object == nil {
				shouldCreateDefault := false
				var b []bool

				b = s.createDefault(&prop, b)

				for _, x := range b {
					if x {
						shouldCreateDefault = x
						break
					}
				}

				if shouldCreateDefault {
					prop.Default.Object = make(map[string]interface{})
				} else {
					continue
				}
			}

			if _, found := x[k]; !found || s.isNonNullableNull(x[k], &prop) {
				if s.isKindInt(prop.Default.Object) {
					x[k] = prop.Default.Object
				} else {
					x[k] = runtime.DeepCopyJSONValue(prop.Default.Object)
				}
			}
		}
		for k := range x {
			if prop, found := ss.Properties[k]; found {
				s.schemaDefault(x[k], &prop)
			} else if ss.AdditionalProperties != nil {
				if s.isNonNullableNull(x[k], ss.AdditionalProperties.Structural) {
					if s.isKindInt(ss.AdditionalProperties.Structural.Default.Object) {
						x[k] = ss.AdditionalProperties.Structural.Default.Object
					} else {
						x[k] = runtime.DeepCopyJSONValue(ss.AdditionalProperties.Structural.Default.Object)
					}
				}

				s.schemaDefault(x[k], ss.AdditionalProperties.Structural)
			}
		}
	case []interface{}:
		for i := range x {
			if s.isNonNullableNull(x[i], ss.Items) {
				if s.isKindInt(ss.Items.Default.Object) {
					x[i] = ss.Items.Default.Object
				} else {
					x[i] = runtime.DeepCopyJSONValue(ss.Items.Default.Object)
				}
			}

			s.schemaDefault(x[i], ss.Items)
		}
	default:
		// scalars, do nothing
	}
}

func (s PackageSchema) isNonNullableNull(x interface{}, ss *structuralschema.Structural) bool {
	return x == nil && ss != nil && !ss.Generic.Nullable
}

func (s PackageSchema) isKindInt(src interface{}) bool {
	if src != nil && reflect.TypeOf(src).Kind() == reflect.Int {
		return true
	}
	return false
}

func (s PackageSchema) createDefault(structural *structuralschema.Structural, b []bool) []bool {
	for _, v := range structural.Properties {
		// return true if there is a non-nested(not object) with a default value
		if v.Type != "object" && v.Default.Object != nil {
			b = append(b, true)
			return b
		}

		if v.Type == "object" && v.Default.Object == nil && v.Properties != nil {
			b = append(b, s.createDefault(&v, b)...)
		}
	}

	b = append(b, false)

	return b
}

func (s PackageSchema) commentDefaultValues(defaultValues bytes.Buffer) []byte {
	var commentedDefaultValues string
	for _, line := range strings.Split(strings.TrimSuffix(defaultValues.String(), "\n"), "\n") {
		line = fmt.Sprintf("# %s\n", line)
		commentedDefaultValues += line
	}
	return []byte(commentedDefaultValues)
}
