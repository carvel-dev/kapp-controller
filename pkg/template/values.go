// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Values struct {
	ValuesFrom []v1alpha1.AppTemplateValuesSource

	appContext AppContext
	coreClient kubernetes.Interface
}

// AsPaths returns set of paths that store values, and a single
// temp directory that acts as a shared parent (used for cleanup of all values).
func (t Values) AsPaths(dirPath string) ([]string, *memdir.TmpDir, error) {
	rootDir := memdir.NewTmpDir("template-values")

	err := rootDir.Create()
	if err != nil {
		return nil, nil, err
	}

	var allPaths []string

	for idx, source := range t.ValuesFrom {
		var paths []string
		var err error

		dstPath := filepath.Join(rootDir.Path(), fmt.Sprintf("%d", idx))

		err = os.Mkdir(dstPath, 0700)
		if err != nil {
			rootDir.Remove()
			return nil, nil, fmt.Errorf("Creating subdir: %s", err)
		}

		switch {
		case source.SecretRef != nil:
			paths, err = t.writeFromSecret(dstPath, *source.SecretRef)

		case source.ConfigMapRef != nil:
			paths, err = t.writeFromConfigMap(dstPath, *source.ConfigMapRef)

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
				items:    source.DownwardAPI.Items,
				metadata: t.appContext.Metadata,
			}
			paths, err = t.writeFromDownwardAPI(dstPath, downwardAPIValues)

		default:
			err = fmt.Errorf("Expected either secretRef, configMapRef or path as a source")
		}
		if err != nil {
			rootDir.Remove()
			return nil, nil, fmt.Errorf("Writing paths: %s", err)
		}

		allPaths = append(allPaths, paths...)
	}

	return allPaths, rootDir, nil
}

func (t Values) writeFromSecret(dstPath string,
	secretRef v1alpha1.AppTemplateValuesSourceRef) ([]string, error) {

	secret, err := t.coreClient.CoreV1().Secrets(t.appContext.Namespace).Get(
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

	err = ioutil.WriteFile(newPath, content, 0600)
	if err != nil {
		return "", fmt.Errorf("Writing file '%s': %s", newPath, err)
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
