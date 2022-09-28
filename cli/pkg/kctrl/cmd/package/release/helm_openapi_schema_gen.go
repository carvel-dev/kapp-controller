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
	typeProp        = "type"
	formatProp      = "format"
	descriptionProp = "description"
	itemsProp       = "items"
	propertiesProp  = "properties"
	defaultProp     = "default"
)

var propOrder = map[string]int{
	typeProp:        1,
	formatProp:      2,
	descriptionProp: 3,
	itemsProp:       4,
	propertiesProp:  5,
	defaultProp:     6,
}

const (
	objectVal = "object"
	arrayVal  = "array"
	floatVal  = "float"
)

// Yaml tags that differentiate the type of scalar object in the node
const (
	nullTag  = "!!null"
	boolTag  = "!!bool"
	intTag   = "!!int"
	floatTag = "!!float"
)

type MapItem struct {
	Key   interface{}
	Value interface{}
}

type Map struct {
	Items []*MapItem
}

func (m *Map) Iterate(iterFunc func(k, v interface{})) {
	for _, item := range m.Items {
		iterFunc(item.Key, item.Value)
	}
}

type openAPIKeys []*MapItem

func (o openAPIKeys) Len() int {
	return len(o)
}

func (o openAPIKeys) Less(i, j int) bool {
	return propOrder[o[i].Key.(string)] < propOrder[o[j].Key.(string)]
}

func (o openAPIKeys) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
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
	if document.Kind == 0 {
		return &kcdatav1alpha1.ValuesSchema{
			OpenAPIv3: runtime.RawExtension{Raw: nil},
		}, nil
	}
	if document.Kind != yaml3.DocumentNode {
		// return proper error message
		return nil, fmt.Errorf("invalid node kind supplied: %d", document.Kind)
	}
	if document.Content[0].Kind != yaml3.MappingNode {
		// return proper error message
		return nil, fmt.Errorf("values file must resolve to a map (was %d)", document.Content[0].Kind)
	}
	openAPIProperties, err := h.calculateProperties(nil, document.Content[0])
	if err != nil {
		return nil, err
	}

	bytes, err := yaml2.Marshal(h.convertToYAML(openAPIProperties))
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
	helmChartDir := h.dirPath
	fileInfo, err := os.Stat(helmChartDir)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %s", h.dirPath, err.Error())
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("expected %s to be directory", h.dirPath)
	}

	_, err = os.Stat(filepath.Join(helmChartDir, "values.yaml"))
	if err != nil {
		// It is possible that helm chart doesn't have values.yml file.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return os.ReadFile(filepath.Join(helmChartDir, "values.yaml"))
}

func (h HelmValuesSchemaGen) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func (h HelmValuesSchemaGen) convertToYAML(val interface{}) interface{} {
	switch typedVal := val.(type) {
	case *Map:
		result := yaml2.MapSlice{}
		typedVal.Iterate(func(k, v interface{}) {
			result = append(result, yaml2.MapItem{
				Key:   k,
				Value: h.convertToYAML(v),
			})
		})
		return result
	default:
		return val
	}
}

func (h HelmValuesSchemaGen) calculateProperties(key *yaml3.Node, value *yaml3.Node) (*Map, error) {
	switch value.Kind {
	case yaml3.MappingNode:
		var apiKeys openAPIKeys
		description, isPresent := h.getDescriptionFromNode(key)
		if isPresent {
			apiKeys = append(apiKeys, &MapItem{Key: descriptionProp, Value: description})
		}
		apiKeys = append(apiKeys, &MapItem{Key: typeProp, Value: objectVal})

		if len(value.Content) == 0 {
			apiKeys = append(apiKeys, &MapItem{Key: defaultProp, Value: "{}"})
		} else {
			var properties openAPIKeys
			for i := 0; i < len(value.Content); i += 2 {
				k := value.Content[i]
				v := value.Content[i+1]
				calculatedProperties, err := h.calculateProperties(k, v)
				if err != nil {
					return nil, err
				}
				valueItems := calculatedProperties.Items[0]
				properties = append(properties, &MapItem{Key: valueItems.Key, Value: valueItems.Value})
			}
			apiKeys = append(apiKeys, &MapItem{Key: propertiesProp, Value: &Map{Items: properties}})
		}

		sort.Sort(apiKeys)
		if key == nil {
			return &Map{Items: apiKeys}, nil
		}
		mi := &MapItem{Key: key.Value, Value: &Map{Items: apiKeys}}
		return &Map{Items: []*MapItem{mi}}, nil
	case yaml3.SequenceNode:
		var defaultVals []interface{}
		var apiKeys openAPIKeys
		apiKeys = append(apiKeys, &MapItem{Key: typeProp, Value: arrayVal})
		description, isPresent := h.getDescriptionFromNode(key)
		if isPresent {
			apiKeys = append(apiKeys, &MapItem{Key: descriptionProp, Value: description})
		}

		if len(value.Content) != 0 {
			properties := &Map{[]*MapItem{}}
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
						if item.Key == propertiesProp {
							properties.Items = append(properties.Items, item.Value.(*Map).Items...)
						}
					}
				case yaml3.SequenceNode:
					calculatedProperties, err := h.calculateProperties(nil, v)
					if err != nil {
						return nil, err
					}
					for _, item := range calculatedProperties.Items {
						if item.Key == itemsProp {
							properties.Items = append(properties.Items, item.Value.(*Map).Items...)
						}
					}
				case yaml3.ScalarNode:
					val, err := h.getDefaultValue(h.openAPITypeFor(value.Content[0].Tag), v.Value)
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
						if item.Key == itemsProp {
							properties.Items = append(properties.Items, item.Value.(*Map).Items...)
						}
					}
				default:
					return nil, fmt.Errorf("Unrecognized type %T", v.Kind)
				}
			}

			var itemsProperties *Map
			switch value.Content[0].Kind {
			case yaml3.MappingNode, yaml3.AliasNode:
				itemsProperties = &Map{[]*MapItem{
					&MapItem{Key: typeProp, Value: "object"},
					&MapItem{Key: propertiesProp, Value: properties}}}
			case yaml3.SequenceNode:
				itemsProperties = &Map{[]*MapItem{
					&MapItem{Key: typeProp, Value: "array"},
					&MapItem{Key: defaultProp, Value: "[]"},
					&MapItem{Key: itemsProp, Value: properties}}}
			case yaml3.ScalarNode:
				itemsProperties = &Map{[]*MapItem{
					&MapItem{Key: typeProp, Value: h.openAPITypeFor(value.Content[0].Tag)}}}
			}
			apiKeys = append(apiKeys, &MapItem{Key: itemsProp, Value: itemsProperties})
		}
		apiKeys = append(apiKeys, &MapItem{Key: defaultProp, Value: defaultVals})
		sort.Sort(apiKeys)
		if key == nil {
			return &Map{Items: apiKeys}, nil
		}
		mi := &MapItem{Key: key.Value, Value: &Map{Items: apiKeys}}
		return &Map{Items: []*MapItem{mi}}, nil
	case yaml3.ScalarNode:
		var apiKeys openAPIKeys
		description, isPresent := h.getDescriptionFromNode(key)
		if isPresent {
			apiKeys = append(apiKeys, &MapItem{Key: descriptionProp, Value: description})
		}

		valType := h.openAPITypeFor(value.Tag)
		defaultVal, err := h.getDefaultValue(valType, value.Value)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, &MapItem{Key: typeProp, Value: valType})
		apiKeys = append(apiKeys, &MapItem{Key: defaultProp, Value: defaultVal})
		if value.Tag == floatTag {
			apiKeys = append(apiKeys, &MapItem{Key: formatProp, Value: floatVal})
		}

		sort.Sort(apiKeys)
		mi := &MapItem{Key: key.Value, Value: &Map{Items: apiKeys}}
		return &Map{Items: []*MapItem{mi}}, nil
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
	if strings.HasPrefix(comment, fmt.Sprintf("# %s", node.Value)) {
		return strings.TrimSpace(strings.TrimPrefix(comment, fmt.Sprintf("# %s", node.Value))), true
	} else if strings.HasPrefix(comment, "# --") {
		return strings.TrimSpace(strings.TrimPrefix(comment, "# --")), true
	} else if strings.HasPrefix(comment, "##") {
		return strings.TrimSpace(strings.TrimPrefix(comment, "##")), true
	} else if strings.HasPrefix(comment, "#") {
		return strings.TrimSpace(strings.TrimPrefix(comment, "#")), true
	}
	return "", false
}

func (h HelmValuesSchemaGen) openAPITypeFor(tag string) string {
	switch tag {
	case boolTag:
		return "boolean"
	case floatTag:
		return "number"
	case intTag:
		return "integer"
	case nullTag:
		return "null"
	default:
		return "string"
	}
}

func (h HelmValuesSchemaGen) getDefaultValue(tag, value string) (interface{}, error) {
	switch tag {
	case "boolean":
		return strconv.ParseBool(value)
	case "integer":
		return strconv.Atoi(value)
	case "number":
		return strconv.ParseFloat(value, 64)
	default:
		return value, nil
	}
}
