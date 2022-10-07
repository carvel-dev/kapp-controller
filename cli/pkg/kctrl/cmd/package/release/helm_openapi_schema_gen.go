// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	yaml2 "gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

// keys used when generating an OpenAPI Document
const (
	typeKey        = "type"
	formatKey      = "format"
	descriptionKey = "description"
	itemsKey       = "items"
	propertiesKey  = "properties"
	defaultKey     = "default"
)

var keyOrder = map[string]int{
	typeKey:        1,
	formatKey:      2,
	descriptionKey: 3,
	itemsKey:       4,
	propertiesKey:  5,
	defaultKey:     6,
}

const (
	objectVal = "object"
	arrayVal  = "array"
	floatVal  = "float"
)

// YAML tags that differentiate the type of scalar object in the node
const (
	nullTag  = "!!null"
	boolTag  = "!!bool"
	intTag   = "!!int"
	floatTag = "!!float"
)

type MapItem struct {
	Key   string
	Value interface{}
}

type Map struct {
	Items []*MapItem
}

func (m *Map) Iterate(iterFunc func(k string, v interface{})) {
	for _, item := range m.Items {
		iterFunc(item.Key, item.Value)
	}
}

type byKey []*MapItem

func (k byKey) Len() int {
	return len(k)
}

func (k byKey) Less(i, j int) bool {
	return keyOrder[k[i].Key] < keyOrder[k[j].Key]
}

func (k byKey) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

type HelmValuesSchemaGen struct {
	// dirPath to the helm chart directory
	dirPath string
}

func NewHelmValuesSchemaGen(path string) HelmValuesSchemaGen {
	return HelmValuesSchemaGen{path}
}

func (h HelmValuesSchemaGen) Schema() (*kcdatav1alpha1.ValuesSchema, error) {
	fileData, err := h.readValuesFile()
	if err != nil {
		return nil, err
	}
	if len(fileData) == 0 {
		return &kcdatav1alpha1.ValuesSchema{
			OpenAPIv3: runtime.RawExtension{Raw: nil},
		}, nil
	}

	var document yaml3.Node
	err = yaml3.Unmarshal(fileData, &document)
	if err != nil {
		return nil, err
	}
	if document.IsZero() {
		return &kcdatav1alpha1.ValuesSchema{
			OpenAPIv3: runtime.RawExtension{Raw: nil},
		}, nil
	}

	if document.Kind != yaml3.DocumentNode {
		return nil, fmt.Errorf("invalid node kind supplied: %d", document.Kind)
	}
	if document.Content[0].Kind != yaml3.MappingNode {
		return nil, fmt.Errorf("values file must resolve to a map (was %d)", document.Content[0].Kind)
	}

	openAPIProperties, err := h.calculateProperties(nil, document.Content[0])
	if err != nil {
		return nil, err
	}

	bytes, err := yaml2.Marshal(h.toYAML(openAPIProperties))
	if err != nil {
		return nil, err
	}
	jsonEncodedBytes, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		return nil, err
	}
	return &kcdatav1alpha1.ValuesSchema{
		OpenAPIv3: runtime.RawExtension{Raw: jsonEncodedBytes},
	}, nil
}

func (h HelmValuesSchemaGen) readValuesFile() ([]byte, error) {
	fileInfo, err := os.Stat(h.dirPath)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("expected %s to be directory", h.dirPath)
	}

	fileData, err := os.ReadFile(filepath.Join(h.dirPath, "values.yaml"))
	if err != nil {
		// It is possible that helm chart doesn't have values.yml file.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return fileData, nil
}

func (h HelmValuesSchemaGen) toYAML(val interface{}) interface{} {
	switch typedVal := val.(type) {
	case *Map:
		result := yaml2.MapSlice{}
		typedVal.Iterate(func(k string, v interface{}) {
			result = append(result, yaml2.MapItem{
				Key:   k,
				Value: h.toYAML(v),
			})
		})
		return result
	default:
		return val
	}
}

func (h HelmValuesSchemaGen) calculateProperties(key *yaml3.Node, value *yaml3.Node) (*Map, error) {
	var apiKeys []*MapItem
	description, isPresent := h.getDescriptionFromNode(key)
	if isPresent {
		apiKeys = append(apiKeys, &MapItem{Key: descriptionKey, Value: description})
	}

	switch value.Kind {
	case yaml3.MappingNode:
		var properties []*MapItem
		apiKeys = append(apiKeys, &MapItem{Key: typeKey, Value: objectVal})

		for i := 0; i < len(value.Content); i += 2 {
			calculatedProperties, err := h.calculateProperties(value.Content[i], value.Content[i+1])
			if err != nil {
				return nil, err
			}
			properties = append(properties, &MapItem{Key: calculatedProperties.Items[0].Key, Value: calculatedProperties.Items[0].Value})
		}
		if len(properties) > 0 {
			apiKeys = append(apiKeys, &MapItem{Key: propertiesKey, Value: &Map{Items: properties}})
		} else {
			apiKeys = append(apiKeys, &MapItem{Key: defaultKey, Value: "{}"})
		}

		sort.Sort(byKey(apiKeys))
		if key == nil {
			return &Map{Items: apiKeys}, nil
		}
		return &Map{Items: []*MapItem{&MapItem{Key: key.Value, Value: &Map{Items: apiKeys}}}}, nil
	case yaml3.SequenceNode:
		var defaultVals []interface{}
		properties := &Map{[]*MapItem{}}
		apiKeys = append(apiKeys, &MapItem{Key: typeKey, Value: arrayVal})

		for _, v := range value.Content {
			if len(v.Content) > 0 && len(v.Content[0].HeadComment) == 0 {
				v.Content[0].HeadComment = v.HeadComment
			}
			switch v.Kind {
			case yaml3.MappingNode:
				calculatedProperties, err := h.calculateProperties(nil, v)
				if err != nil {
					return nil, err
				}
				for _, item := range calculatedProperties.Items {
					if item.Key == propertiesKey {
						properties.Items = append(properties.Items, item.Value.(*Map).Items...)
					}
				}
			case yaml3.SequenceNode:
				calculatedProperties, err := h.calculateProperties(nil, v)
				if err != nil {
					return nil, err
				}
				for _, item := range calculatedProperties.Items {
					if item.Key == itemsKey {
						properties.Items = append(properties.Items, item.Value.(*Map).Items...)
					}
				}
			case yaml3.ScalarNode:
				val, err := h.getDefaultValue(v.Tag, v.Value)
				if err != nil {
					return nil, err
				}
				defaultVals = append(defaultVals, val)
			case yaml3.AliasNode:
				calculatedProperties, err := h.calculateProperties(nil, v.Alias)
				if err != nil {
					return nil, err
				}
				for _, item := range calculatedProperties.Items {
					if item.Key == itemsKey {
						properties.Items = append(properties.Items, item.Value.(*Map).Items...)
					}
				}
			default:
				return nil, fmt.Errorf("Unrecognized type %T", v.Kind)
			}
		}

		if len(value.Content) != 0 {
			var itemsProperties *Map
			switch value.Content[0].Kind {
			case yaml3.MappingNode, yaml3.AliasNode:
				itemsProperties = &Map{[]*MapItem{
					&MapItem{Key: typeKey, Value: "object"},
					&MapItem{Key: propertiesKey, Value: properties}}}
			case yaml3.SequenceNode:
				itemsProperties = &Map{[]*MapItem{
					&MapItem{Key: typeKey, Value: "array"},
					&MapItem{Key: defaultKey, Value: "[]"},
					&MapItem{Key: itemsKey, Value: properties}}}
			case yaml3.ScalarNode:
				itemsProperties = &Map{[]*MapItem{
					&MapItem{Key: typeKey, Value: h.openAPIType(value.Content[0].Tag, value.Content[0].Value)}}}
			}
			apiKeys = append(apiKeys, &MapItem{Key: itemsKey, Value: itemsProperties})
		}
		apiKeys = append(apiKeys, &MapItem{Key: defaultKey, Value: defaultVals})
		sort.Sort(byKey(apiKeys))
		if key == nil {
			return &Map{Items: apiKeys}, nil
		}
		return &Map{Items: []*MapItem{&MapItem{Key: key.Value, Value: &Map{Items: apiKeys}}}}, nil
	case yaml3.ScalarNode:
		defaultVal, err := h.getDefaultValue(value.Tag, value.Value)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, &MapItem{Key: typeKey, Value: h.openAPIType(value.Tag, value.Value)})
		apiKeys = append(apiKeys, &MapItem{Key: defaultKey, Value: defaultVal})
		if value.Tag == floatTag {
			apiKeys = append(apiKeys, &MapItem{Key: formatKey, Value: floatVal})
		}

		sort.Sort(byKey(apiKeys))
		return &Map{Items: []*MapItem{&MapItem{Key: key.Value, Value: &Map{Items: apiKeys}}}}, nil
	case yaml3.AliasNode:
		return h.calculateProperties(key, value.Alias)
	default:
		return nil, fmt.Errorf("Unrecognized type %T", value.Kind)
	}
	return nil, nil
}

func (h HelmValuesSchemaGen) getDescriptionFromNode(node *yaml3.Node) (string, bool) {
	if node == nil || node.HeadComment == "" {
		return "", false
	}

	comment := node.HeadComment
	comment = strings.ReplaceAll(comment, "\n##", "")
	comment = strings.ReplaceAll(comment, "\n#", "")
	switch {
	case strings.HasPrefix(comment, fmt.Sprintf("# %s", node.Value)):
		return strings.TrimSpace(strings.TrimPrefix(comment, fmt.Sprintf("# %s", node.Value))), true
	case strings.HasPrefix(comment, "# --"):
		return strings.TrimSpace(strings.TrimPrefix(comment, "# --")), true
	case strings.HasPrefix(comment, "##"):
		return strings.TrimSpace(strings.TrimPrefix(comment, "##")), true
	case strings.HasPrefix(comment, "#"):
		return strings.TrimSpace(strings.TrimPrefix(comment, "#")), true
	}
	return "", false
}

func (h HelmValuesSchemaGen) openAPIType(tag, value string) string {
	switch tag {
	case boolTag:
		return "boolean"
	case floatTag:
		return "number"
	case intTag:
		return "integer"
	case nullTag:
		if value == "null" {
			return "null"
		}
	}
	return "string"

}

func (h HelmValuesSchemaGen) getDefaultValue(tag, value string) (interface{}, error) {
	switch tag {
	case boolTag:
		return strconv.ParseBool(value)
	case intTag:
		return strconv.Atoi(value)
	case floatTag:
		return strconv.ParseFloat(value, 64)
	default:
		return value, nil
	}
}
