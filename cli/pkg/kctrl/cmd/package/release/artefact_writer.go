// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type ArtefactWriter struct {
	Package string
	Version string
}

const (
	artefactDir = "carvel-artefacts"
	packageDir  = "packages"
)

func NewArtefactWriter(pkg string, version string) *ArtefactWriter {
	return &ArtefactWriter{Package: pkg, Version: version}
}

func (w *ArtefactWriter) CreatePackageDir() error {
	path := filepath.Join(artefactDir, packageDir, w.Package)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (w *ArtefactWriter) TouchPackageMetadata() error {
	metadata := kcdatav1alpha1.PackageMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Name: w.Package,
		},
		Spec: kcdatav1alpha1.PackageMetadataSpec{
			DisplayName: w.Package[:strings.IndexByte(w.Package, '.')],
		},
	}
	template := `# longDescription: Detailed description of package
# shortDescription: Concise description of package
# providerName: Organisation/entity providing this package
# maintainers:
#   - name: Maintainer 1
#   - name: Maintainer 2
`
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}
	path := filepath.Join(artefactDir, packageDir, w.Package, "metadata.yml")
	return w.createFileIfNotExists(path, append(metadataBytes, []byte(template)...))
}

func (w *ArtefactWriter) WritePackageFile(imgpkgBundleLocation string, buildAppSpec *kcv1alpha1.AppSpec, useKbldLockOutput bool) error {
	packageMeta := kcdatav1alpha1.Package{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", w.Package, w.Version),
		},
		Spec: kcdatav1alpha1.PackageSpec{
			Version: w.Version,
			RefName: w.Package,
			Template: kcdatav1alpha1.AppTemplateSpec{
				Spec: &kcv1alpha1.AppSpec{
					Fetch: []kcv1alpha1.AppFetch{
						{
							ImgpkgBundle: &kcv1alpha1.AppFetchImgpkgBundle{
								Image: imgpkgBundleLocation,
							},
						},
					},
					Template: buildAppSpec.Template,
					Deploy:   buildAppSpec.Deploy,
				},
			},
		},
	}
	if useKbldLockOutput {
		packageMeta.Spec.Template.Spec.Template = append(packageMeta.Spec.Template.Spec.Template,
			kcv1alpha1.AppTemplate{Kbld: &kcv1alpha1.AppTemplateKbld{
				Paths: []string{"-", ".imgpkg/images.yml"},
			}})
	}

	packageBytes, err := yaml.Marshal(packageMeta)
	if err != nil {
		return err
	}
	path := filepath.Join(artefactDir, packageDir, w.Package, fmt.Sprintf("%s.yml", w.Version))
	return w.createOrOverwriteFile(path, packageBytes)
}

func (w *ArtefactWriter) createFileIfNotExists(path string, data []byte) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := ioutil.WriteFile(path, data, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *ArtefactWriter) createOrOverwriteFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
