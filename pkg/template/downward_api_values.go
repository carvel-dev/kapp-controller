// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/util/jsonpath"
)

// AdditionalDownwardAPIValues holds values that are not computed discoverable on the resource itself
type AdditionalDownwardAPIValues struct {
	KappControllerVersion func() (string, error)
	KubernetesVersion     func() (string, error)
	KubernetesAPIs        func() ([]string, error)
}

// DownwardAPIValues produces multiple key-values based on the DownwardAPI config
// queried against the object metadata
type DownwardAPIValues struct {
	items                       []v1alpha1.AppTemplateValuesDownwardAPIItem
	metadata                    PartialObjectMetadata
	additionalDownwardAPIValues AdditionalDownwardAPIValues
}

// AsYAMLs returns many key-values queried (using jsonpath) against an object metadata provided, and use additionalValues.
func (a DownwardAPIValues) AsYAMLs() ([][]byte, error) {
	dataValues := [][]byte{}
	keyValueContent := []byte{}
	var err error

	for _, item := range a.items {
		err = a.validateName(item.Name)
		if err != nil {
			return nil, err
		}

		switch {
		case item.FieldPath != "":
			var fieldPathExpression string
			fieldPathExpression, err = relaxedJSONPathExpression(item.FieldPath)
			if err != nil {
				return nil, err
			}
			keyValueContent, err = a.extractFieldPathAsKeyValue(item.Name, fieldPathExpression)
		case item.KubernetesVersion != nil:
			if item.KubernetesVersion.Version == "" { // I wish there was a ternary operator in Go
				v, err := a.additionalDownwardAPIValues.KubernetesVersion()
				if err != nil {
					return nil, err
				}
				keyValueContent, err = yaml.Marshal(map[string]string{item.Name: v})
			} else {
				keyValueContent, err = yaml.Marshal(map[string]string{item.Name: item.KubernetesVersion.Version})
			}
		case item.KappControllerVersion != nil:
			if item.KappControllerVersion.Version == "" {
				v, err := a.additionalDownwardAPIValues.KappControllerVersion()
				if err != nil {
					return nil, err
				}
				keyValueContent, err = yaml.Marshal(map[string]string{item.Name: v})
			} else {
				keyValueContent, err = yaml.Marshal(map[string]string{item.Name: item.KappControllerVersion.Version})
			}
		case item.KubernetesAPIs != nil:
			if item.KubernetesAPIs.GroupVersions == nil {
				v, err := a.additionalDownwardAPIValues.KubernetesAPIs()
				if err != nil {
					return nil, err
				}
				keyValueContent, err = yaml.Marshal(map[string]interface{}{item.Name: v})
			} else {
				keyValueContent, err = yaml.Marshal(map[string]interface{}{item.Name: item.KubernetesAPIs.GroupVersions})
			}
		default:
			return nil, fmt.Errorf("Invalid downward API item given")
		}

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

	return yaml.Marshal(root)
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
	return append(tokens, string(token))
}

// Copied from https://github.com/kubernetes/kubectl/blob/ac26f503e81287d9903761a1a8ded25fdebec6a7/pkg/cmd/get/customcolumn.go#L38
var jsonRegexp = regexp.MustCompile(`^\{\.?([^{}]+)\}$|^\.?([^{}]+)$`)

// RelaxedJSONPathExpression attempts to be flexible with JSONPath expressions, it accepts:
//   - metadata.name (no leading '.' or curly braces '{...}'
//   - {metadata.name} (no leading '.')
//   - .metadata.name (no curly braces '{...}')
//   - {.metadata.name} (complete expression)
//
// And transforms them all into a valid jsonpath expression:
//
//	{.metadata.name}
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
