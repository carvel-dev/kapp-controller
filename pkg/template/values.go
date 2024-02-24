// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
package template

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Values struct {
	ValuesFrom       []v1alpha1.AppTemplateValuesSource
	AdditionalValues AdditionalDownwardAPIValues
	appContext       AppContext
	coreClient       kubernetes.Interface
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
			downwardAPIValues := DownwardAPIValues{
				items:                       source.DownwardAPI.Items,
				metadata:                    t.appContext.Metadata,
				additionalDownwardAPIValues: t.AdditionalValues,
			}
			paths, err = t.writeFromDownwardAPI(valuesDir.Path(), downwardAPIValues)
		default:
			err = fmt.Errorf("expected one of secretRef, configMapRef, downwardAPI, or path as a source")
		}
		if err != nil {
			cleanUpFunc()
			return nil, nil, fmt.Errorf("preparing template values: %s", err)
		}
		allPaths = append(allPaths, paths...)
	}
	return allPaths, cleanUpFunc, nil
}
func (t Values) writeFromSecret(dstPath string,
	secretRef v1alpha1.AppTemplateValuesSourceRef) ([]string, error) {
	secret, err := t.coreClient.CoreV1().Secrets(t.appContext.Namespace).Get(
		context.Background(), secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	if secretRef.Key != "" {
		// If a key is provided, only write the value for that key
		val, ok := secret.Data[secretRef.Key]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found in secret '%s'", secretRef.Key, secretRef.Name)
		}
		path, err := t.writeFile(dstPath, secretRef.Key, val)
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	} else {
		// If no key is provided, write all values in the secret
		for name, val := range secret.Data {
			path, err := t.writeFile(dstPath, name, val)
			if err != nil {
				return nil, err
			}
			result = append(result, path)
		}
	}

	sort.Strings(result)
	return result, nil
}

func (t Values) writeFromDownwardAPI(dstPath string, valuesExtractor DownwardAPIValues) ([]string, error) {
	var result []string
	dataValues, err := valuesExtractor.AsYAMLs()
	if err != nil {
		return nil, err
	}
	for idx, content := range dataValues {
		path, err := t.writeFile(dstPath, fmt.Sprintf("downwardapi_%d.yaml", idx), content)
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	}
	return result, nil
}
func (t Values) writeFile(dstPath, subPath string, content []byte) (string, error) {
	newPath, err := memdir.ScopedPath(dstPath, subPath)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(newPath, content, 0600)
	if err != nil {
		return "", fmt.Errorf("writing file '%s': %s", newPath, err)
	}
	return newPath, nil
}
func (t Values) writeFromConfigMap(dstPath string,
	configMapRef v1alpha1.AppTemplateValuesSourceRef) ([]string, error) {
	configMap, err := t.coreClient.CoreV1().ConfigMaps(t.appContext.Namespace).Get(
		context.Background(), configMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	if configMapRef.Key != "" {
		// If a key is provided, only write the value for that key
		val, ok := configMap.Data[configMapRef.Key]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found in config map '%s'", configMapRef.Key, configMapRef.Name)
		}
		path, err := t.writeFile(dstPath, configMapRef.Key, []byte(val))
		if err != nil {
			return nil, err
		}
		result = append(result, path)
	} else {
		// If no key is provided, write all values in the config map
		for name, val := range configMap.Data {
			path, err := t.writeFile(dstPath, name, []byte(val))
			if err != nil {
				return nil, err
			}
			result = append(result, path)
		}
	}

	sort.Strings(result)
	return result, nil
}
