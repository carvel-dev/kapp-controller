// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/util/jsonpath"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Values struct {
	ValuesFrom []v1alpha1.AppTemplateValuesSource

	genericOpts GenericOpts
	coreClient  kubernetes.Interface
}

func (t Values) AsPaths(dirPath string) ([]string, func(), error) {
	var valuesDirs []*memdir.TmpDir

	cleanUpFunc := func() {
		for _, valDir := range valuesDirs {
			_ = valDir.Remove()
		}
	}

	var allPaths []string

	for _, source := range t.ValuesFrom {
		var paths []string
		var err error

		valuesDir := memdir.NewTmpDir("template-values")

		err = valuesDir.Create()
		if err != nil {
			cleanUpFunc()
			return nil, nil, err
		}

		valuesDirs = append(valuesDirs, valuesDir)

		switch {
		case source.SecretRef != nil:
			paths, err = t.writeFromSecret(valuesDir.Path(), *source.SecretRef)

		case source.ConfigMapRef != nil:
			paths, err = t.writeFromConfigMap(valuesDir.Path(), *source.ConfigMapRef)

		case len(source.Path) > 0:
			if source.Path == stdinPath {
				paths = append(paths, stdinPath)
			} else {
				checkedPath, err := memdir.ScopedPath(dirPath, source.Path)
				if err == nil {
					paths = append(paths, checkedPath)
				}
			}

		case source.DownwardAPI != nil:
			paths, err = t.writeFromDownwardAPI(valuesDir.Path(), *source.DownwardAPI)

		default:
			err = fmt.Errorf("Expected either secretRef, configMapRef or path as a source")
		}
		if err != nil {
			cleanUpFunc()
			return nil, nil, fmt.Errorf("Writing paths: %s", err)
		}

		allPaths = append(allPaths, paths...)
	}

	return allPaths, cleanUpFunc, nil
}

func (t Values) writeFromSecret(dstPath string,
	secretRef v1alpha1.AppTemplateValuesSourceRef) ([]string, error) {

	secret, err := t.coreClient.CoreV1().Secrets(t.genericOpts.Namespace).Get(
		context.Background(), secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	for name, val := range secret.Data {
		path, err := t.writeFile(dstPath, name, val)
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}

	sort.Strings(result)

	return result, nil
}

func (t Values) writeFromDownwardAPI(dstPath string, downwardAPIRef v1alpha1.AppTemplateValuesDownwardAPI) ([]string, error) {
	var result []string

	for idx, item := range downwardAPIRef.Items {
		fieldPathExpression, err := relaxedJSONPathExpression(item.FieldPath)
		if err != nil {
			return nil, err
		}

		content, err := t.extractFieldPathAsKeyValue(item.Name, fieldPathExpression)
		if err != nil {
			return nil, err
		}

		path, err := t.writeFile(dstPath, fmt.Sprintf("downwardapi_%d.yaml", idx), content)
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}

	return result, nil
}

func (t Values) extractFieldPathAsKeyValue(name string, fieldPath string) ([]byte, error) {
	path := jsonpath.New(name).AllowMissingKeys(false)
	err := path.Parse(fieldPath)
	if err != nil {
		return nil, err
	}

	results, err := path.FindResults(t.genericOpts.Metadata)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 || len(results[0]) != 1 {
		return nil, errors.New("invalid field spec provided to DownwardAPI. Only single supported fields are allowed")
	}

	result := results[0][0]
	kind := result.Kind()
	if kind == reflect.Interface {
		kind = result.Elem().Kind()
	}
	if kind == reflect.Map {
		return t.keyValues(name, result.Interface().(map[string]string))
	}

	return t.keyValue(name, fmt.Sprintf("%v", result.Interface())), nil
}

func (t Values) keyValue(key string, val interface{}) []byte {
	return []byte(fmt.Sprintf("%s: %q", key, val))
}

func (t Values) keyValues(key string, m map[string]string) ([]byte, error) {
	// output with keys in sorted order to provide stable output
	fmtStr := fmt.Sprintf("%s: \n", key)
	keys := sets.NewString()
	for k := range m {
		keys.Insert(k)
	}
	tabIndent := "  "
	for _, k := range keys.List() {
		fmtStr += fmt.Sprintf("%s%v: %q\n", tabIndent, k, m[k])
	}
	fmtStr = strings.TrimSuffix(fmtStr, "\n")

	return []byte(fmtStr), nil
}

func (t Values) writeFromConfigMap(dstPath string,
	configMapRef v1alpha1.AppTemplateValuesSourceRef) ([]string, error) {

	configMap, err := t.coreClient.CoreV1().ConfigMaps(t.genericOpts.Namespace).Get(
		context.Background(), configMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	for name, val := range configMap.Data {
		path, err := t.writeFile(dstPath, name, []byte(val))
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}

	sort.Strings(result)

	return result, nil
}

func (t Values) writeFile(dstPath, subPath string, content []byte) (string, error) {
	newPath, err := memdir.ScopedPath(dstPath, subPath)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(newPath, content, 0600)
	if err != nil {
		return "", fmt.Errorf("Writing file '%s': %s", newPath, err)
	}

	return newPath, nil
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
