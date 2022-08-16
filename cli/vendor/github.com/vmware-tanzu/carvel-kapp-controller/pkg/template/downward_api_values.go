// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/util/jsonpath"
)

// DownwardAPIValues produces multiple key-values based on the DownwardAPI config
// queried against the object metadata
type DownwardAPIValues struct {
	items    []v1alpha1.AppTemplateValuesDownwardAPIItem
	metadata PartialObjectMetadata
}

// AsYAMLs returns many key-values queried (using jsonpath) against an object metadata provided.
func (a DownwardAPIValues) AsYAMLs() ([][]byte, error) {
	dataValues := [][]byte{}
	for _, item := range a.items {
		err := a.validateName(item.Name)
		if err != nil {
			return nil, err
		}

		fieldPathExpression, err := relaxedJSONPathExpression(item.FieldPath)
		if err != nil {
			return nil, err
		}

		keyValueContent, err := a.extractFieldPathAsKeyValue(item.Name, fieldPathExpression)
		if err != nil {
			return nil, err
		}
		dataValues = append(dataValues, keyValueContent)
	}

	return dataValues, nil
}

func (a DownwardAPIValues) validateName(name string) error {
	if strings.HasSuffix(name, ".") || strings.HasPrefix(name, ".") || strings.Contains(name, "..") {
		return fmt.Errorf("Invalid name was provided '%s' (hint: separate paths should only use a single '.' character)", name)
	}

	return nil
}

func (a DownwardAPIValues) extractFieldPathAsKeyValue(name string, fieldPath string) ([]byte, error) {
	path := jsonpath.New(name).AllowMissingKeys(false)
	err := path.Parse(fieldPath)
	if err != nil {
		return nil, err
	}

	results, err := path.FindResults(a.metadata)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 || len(results[0]) != 1 {
		return nil, errors.New("Invalid field spec provided to DownwardAPI. Only single supported fields are allowed")
	}

	return a.nestedKeyValue(name, results[0][0].Interface())
}

// operator may wish to assign a downward API value into a nested key structure to use within their template
func (DownwardAPIValues) nestedKeyValue(name string, val interface{}) ([]byte, error) {
	root := map[string]interface{}{}
	nestedMap := root

	nestedKeys := splitWithEscaping(name, '.', '\\')
	for i := 0; i < len(nestedKeys); i++ {
		if i == len(nestedKeys)-1 {
			nestedMap[nestedKeys[i]] = val
			break
		}

		nextLevel := map[string]interface{}{}
		nestedMap[nestedKeys[i]] = nextLevel
		nestedMap = nextLevel
	}

	yamlContent, err := yaml.Marshal(root)
	return yamlContent, err
}

func splitWithEscaping(s string, separator, escape byte) []string {
	var token []byte
	var tokens []string
	for i := 0; i < len(s); i++ {
		if s[i] == separator {
			tokens = append(tokens, string(token))
			token = token[:0]
		} else if s[i] == escape && i+1 < len(s) {
			i++
			token = append(token, s[i])
		} else {
			token = append(token, s[i])
		}
	}
	tokens = append(tokens, string(token))
	return tokens
}

// Copied from https://github.com/kubernetes/kubectl/blob/ac26f503e81287d9903761a1a8ded25fdebec6a7/pkg/cmd/get/customcolumn.go#L38
var jsonRegexp = regexp.MustCompile(`^\{\.?([^{}]+)\}$|^\.?([^{}]+)$`)

// RelaxedJSONPathExpression attempts to be flexible with JSONPath expressions, it accepts:
//   * metadata.name (no leading '.' or curly braces '{...}'
//   * {metadata.name} (no leading '.')
//   * .metadata.name (no curly braces '{...}')
//   * {.metadata.name} (complete expression)
// And transforms them all into a valid jsonpath expression:
//   {.metadata.name}
func relaxedJSONPathExpression(pathExpression string) (string, error) {
	if len(pathExpression) == 0 {
		return pathExpression, nil
	}
	submatches := jsonRegexp.FindStringSubmatch(pathExpression)
	if submatches == nil {
		return "", fmt.Errorf("unexpected path string, expected a 'name1.name2' or '.name1.name2' or '{name1.name2}' or '{.name1.name2}'")
	}
	if len(submatches) != 3 {
		return "", fmt.Errorf("unexpected submatch list: %v", submatches)
	}
	var fieldSpec string
	if len(submatches[1]) != 0 {
		fieldSpec = submatches[1]
	} else {
		fieldSpec = submatches[2]
	}
	return fmt.Sprintf("{.%s}", fieldSpec), nil
}
