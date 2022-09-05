// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package available

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
)

// DataValueProperty holds the details of each property under Carvel package.spec.valuesSchema.openAPIv3.properties.
// The example of the schema could be found at: https://carvel.dev/kapp-controller/docs/latest/packaging/#package-1
// From above example, we would have the following:
//
//	DataValueProperty.Key = "namespace"
//	DataValueProperty.Type = "string"
//	DataValueProperty.Description = "Namespace where fluent-bit will be installed."
//	DataValueProperty.Default = "fluent-bit"
type DataValueProperty struct {
	Key         string
	Type        string
	Description string
	Default     interface{}
}

// PackageValuesSchemaParser loads Carvel package values schema and extract property details
type PackageValuesSchemaParser struct {
	Doc                 *openapi3.T
	DataValueProperties []DataValueProperty
	fieldPath           []string // This is the internal variable to temporarily hold the path of a field when traversing on Doc.
}

func NewValuesSchemaParser(valuesSchema v1alpha1.ValuesSchema) (*PackageValuesSchemaParser, error) {
	loader := openapi3.NewLoader()
	doc, loadErr := loader.LoadFromData(valuesSchema.OpenAPIv3.Raw)
	if loadErr != nil {
		return nil, loadErr
	}
	return &PackageValuesSchemaParser{Doc: doc, DataValueProperties: []DataValueProperty{}}, nil
}

// ParseProperties parses the loaded doc and feed the details into []DataValueProperty
func (parser *PackageValuesSchemaParser) ParseProperties() ([]DataValueProperty, error) {
	walkErr := parser.walkOnValueSchemaProperties(parser.Doc.Extensions)
	if walkErr != nil {
		// returning []DataValueProperty{} is to be safe when caller uses len() or for loop, nil in above case will cause
		// panic unexpectedly.
		return []DataValueProperty{}, walkErr
	}
	return parser.DataValueProperties, nil
}

// walkOnValueSchemaProperties is a recursive function to walk on the given docMap and store the details of each property
// within res *[]DataValueProperty. Caller could use the returned *[]DataValueProperty for any purposes.
// Parameter docMap map[string]interface{} is a map representation of Carvel package.spec.valuesSchema.openAPIv3. From the example
// of https://carvel.dev/kapp-controller/docs/latest/packaging/#package-1, docMap should have `title`, `examples` and
// `properties` as top-level keys which might hold nested maps.
func (parser *PackageValuesSchemaParser) walkOnValueSchemaProperties(docMap map[string]interface{}) error {
	var properties interface{}
	var exist bool

	// Base case one: if current level does not have properties field, we do not need to go deeper to look for
	// the nested properties.
	if properties, exist = docMap["properties"]; !exist && len(parser.fieldPath) > 0 {
		propertyType, _ := docMap["type"].(string)
		description, _ := docMap["description"].(string)
		parser.DataValueProperties = append(parser.DataValueProperties, DataValueProperty{
			Key:         strings.Join(parser.fieldPath, "."),
			Type:        propertyType,
			Description: description,
			Default:     docMap["default"],
		})
		return nil
	}

	var propertiesMap map[string]interface{}
	var err error

	// properties is an interface, it needs to be casted to proper type
	switch t := properties.(type) {
	case json.RawMessage:
		err = json.Unmarshal(properties.(json.RawMessage), &propertiesMap)
		if err != nil {
			return err
		}
	case map[string]interface{}:
		propertiesMap = properties.(map[string]interface{})
	default:
		errMsg := fmt.Sprintf("unable to parse the value schema, the value of key named 'properties' has unsupported"+
			" type %v. Expected types are: [map[string]interface{}, json.RawMessage]", t)
		return fmt.Errorf(errMsg)
	}

	// Base case two: if current level does have properties field but that interface is an empty map, we do not need to
	// look for nested properties.
	if len(propertiesMap) == 0 {
		propertyType, _ := docMap["type"].(string)
		description, _ := docMap["description"].(string)
		parser.DataValueProperties = append(parser.DataValueProperties, DataValueProperty{
			Key:         strings.Join(parser.fieldPath, "."),
			Type:        propertyType,
			Description: description,
			Default:     docMap["default"],
		})
		return nil
	}

	// Recursively look into nested map entries in order to load details of properties
	for k, v := range propertiesMap {
		if _, ok := v.(map[string]interface{}); ok {
			parser.fieldPath = append(parser.fieldPath, k)
			_ = parser.walkOnValueSchemaProperties(v.(map[string]interface{}))
			parser.fieldPath = parser.fieldPath[:len(parser.fieldPath)-1]
		}
	}
	return nil
}
