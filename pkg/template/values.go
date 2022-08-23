// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"context"
	"fmt"
	"os"
	"sort"

	semver "github.com/k14s/semver/v4"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/clusterclient"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Values contains the appContext and the ability to fetch the values
type Values struct {
	ValuesFrom []v1alpha1.AppTemplateValuesSource

	appContext     AppContext
	versionFetcher *fetch.VersionFetch
	coreClient     kubernetes.Interface
}

// ValuesFactory abstracts the factories for fetching the values away from the template factory
type ValuesFactory struct {
	appContext   AppContext
	fetchFactory fetch.Factory
	coreClient   kubernetes.Interface
}

// NewValues returns a Values struct based on the app template source and the app context
func (vf ValuesFactory) NewValues(valuesFrom []v1alpha1.AppTemplateValuesSource, appContext AppContext) Values {
	return Values{ValuesFrom: valuesFrom, appContext: appContext, versionFetcher: vf.fetchFactory.NewVersionFetcher(), coreClient: vf.coreClient}
}

// AsPaths returns a list of directories containing values files which are passed to the various templating tools
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
				items:    source.DownwardAPI.Items,
				metadata: t.appContext.Metadata,
				kubernetesVersion: func() (semver.Version, error) {
					return t.versionFetcher.GetKubernetesVersion(t.appContext.ServiceAccountName, t.appContext.AppSpec.Cluster, clusterclient.GenericOpts{Name: t.appContext.Name, Namespace: t.appContext.Namespace})
				},
				kappControllerVersion: t.versionFetcher.GetKappControllerVersion(),
			}
			paths, err = t.writeFromDownwardAPI(valuesDir.Path(), downwardAPIValues)

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

	dataValues, err := valuesExtractor.AsYAML()
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
		return "", fmt.Errorf("Writing file '%s': %s", newPath, err)
	}

	return newPath, nil
}

func (t Values) writeFromConfigMap(dstPath string, configMapRef v1alpha1.AppTemplateValuesSourceRef) ([]string, error) {
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
