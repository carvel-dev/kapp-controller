// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		// Pick up stdin input, without trying to traverse `-`
		if file == "-" {
			bytes, err := cmdcore.NewInputFile(file).Bytes()
			if err != nil {
				return nil, fmt.Errorf("Reading file: %s", err.Error())
			}
			filePathsMap[fmt.Sprintf("%04d-stdin.yaml", i)] = bytes
			continue
		}

		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("Checking file '%s'", file)
		}

		if fileInfo.IsDir() {
			err = filepath.Walk(file, func(path string, fi os.FileInfo, err error) error {
				if err != nil || fi.IsDir() {
					return err
				}
				ext := filepath.Ext(path)
				for _, allowedExt := range fileResourcesAllowedExts {
					if allowedExt == ext {
						relPath, err := filepath.Rel(file, path)
						if err != nil {
							return err
						}

						// Ensures that nested directories are not allowed
						// Accounts for windows environments
						if strings.Count(relPath, "/") > 0 {
							return fmt.Errorf("Nested directories are not supported by `--ytt-overlay-file`")
						}

						// Ignore hidden files like `.git`
						// TODO: Should we explicitly exclude hidden files in Windows?
						if filepath.Base(relPath)[0:1] == "." {
							continue
						}

						bytes, err := cmdcore.NewInputFile(path).Bytes()
						if err != nil {
							return fmt.Errorf("Reading file: %s", err.Error())
						}

						key := fmt.Sprintf("%04d-%s", i, relPath)

						filePathsMap[key] = bytes
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("Recursing through directory: %s: %s", file, err)
			}
		} else {
			for _, allowedExt := range fileResourcesAllowedExts {
				ext := filepath.Ext(file)
				if allowedExt == ext {
					bytes, err := cmdcore.NewInputFile(file).Bytes()
					if err != nil {
						return nil, fmt.Errorf("Reading file: %s", err.Error())
					}
					filePathsMap[fmt.Sprintf("%04d-%s", i, filepath.Base(file))] = bytes
				}
			}
		}
	}
	secret.Data = filePathsMap

	return secret, nil
}
