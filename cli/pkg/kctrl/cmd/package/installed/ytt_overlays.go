// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"fmt"
	"os"
	"path/filepath"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	fileResourcesAllowedExts = []string{".yaml", ".yml"}
)

type YttOverlays struct {
	packageInstall string
	namespace      string
	files          []string
}

func NewYttOverlays(files []string, packageInstall string, namespace string) *YttOverlays {
	return &YttOverlays{files: files, packageInstall: packageInstall, namespace: namespace}
}

func (o *YttOverlays) OverlaysSecret() (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s-overlays", o.packageInstall, o.namespace),
			Annotations: map[string]string{
				KctrlPkgAnnotation: NewCreatedResourceAnnotations(o.packageInstall, o.namespace).PackageAnnValue(),
			},
		},
	}

	filePathsMap := map[string][]byte{}
	for i, file := range o.files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("Checking file '%s'", file)
		}

		if fileInfo.IsDir() {
			err := filepath.Walk(file, func(path string, fi os.FileInfo, err error) error {
				if err != nil || fi.IsDir() {
					return err
				}
				ext := filepath.Ext(path)
				for _, allowedExt := range fileResourcesAllowedExts {
					if allowedExt == ext {
						bytes, err := cmdcore.NewLocalFileSource(path).Bytes()
						if err != nil {
							return fmt.Errorf("Reading file: %s", err.Error())
						}
						filePathsMap[fmt.Sprintf("%04d-%s", i, path)] = bytes
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("Listing files '%s'", file)
			}
		} else {
			for _, allowedExt := range fileResourcesAllowedExts {
				ext := filepath.Ext(file)
				if allowedExt == ext {
					bytes, err := cmdcore.NewLocalFileSource(file).Bytes()
					if err != nil {
						return nil, fmt.Errorf("Reading file: %s", err.Error())
					}
					filePathsMap[fmt.Sprintf("%04d-%s", i, file)] = bytes
				}
			}
		}
	}
	secret.Data = filePathsMap

	return secret, nil
}
