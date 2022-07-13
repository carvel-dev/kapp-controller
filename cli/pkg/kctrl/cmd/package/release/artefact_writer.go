// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type ArtifactWriter struct {
	Package     string
	Version     string
	ArtefactDir string

	ui cmdcore.AuthoringUI
}

const (
	artefactDir = "carvel-artefacts"
	packageDir  = "packages"
)

func NewArtifactWriter(pkg string, version string, artefactDir string, ui cmdcore.AuthoringUI) *ArtifactWriter {
	return &ArtifactWriter{Package: pkg, Version: version, ArtefactDir: artefactDir, ui: ui}
}

func (w *ArtifactWriter) Write(appSpec *kcv1alpha1.AppSpec) error {
	path := filepath.Join(w.ArtefactDir, packageDir, w.Package)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = w.writePackageMetadata(filepath.Join(path, "metadata.yml"))
	if err != nil {
		return err
	}
	err = w.writePackage(filepath.Join(path, "package.yml"), appSpec)
	if err != nil {
		return err
	}

	return nil
}

func (w *ArtifactWriter) WriteRepoOutput(appSpec *kcv1alpha1.AppSpec, path string) error {
	path = filepath.Join(path, packageDir, w.Package)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = w.writePackageMetadata(filepath.Join(path, "metadata.yml"))
	if err != nil {
		return err
	}
	err = w.writePackage(filepath.Join(path, fmt.Sprintf("%s.yml", w.Version)), appSpec)
	if err != nil {
		return err
	}

	return nil
}

func (w *ArtifactWriter) writePackageMetadata(path string) error {
	metadata := kcdatav1alpha1.PackageMetadata{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "data.packaging.carvel.dev/v1alpha1",
			Kind:       "PackageMetadata",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: w.Package,
		},
		Spec: kcdatav1alpha1.PackageMetadataSpec{
			DisplayName: w.Package[:strings.IndexByte(w.Package, '.')],
		},
	}
	template := `# longDescription: Detailed description of package
# shortDescription: Concise description of package
# providerName: Organization/entity providing this package
# maintainers:
#   - name: Maintainer 1
#   - name: Maintainer 2
`
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}
	return w.createFileIfNotExists(path, append(metadataBytes, []byte(template)...))
}

func (w *ArtifactWriter) writePackage(path string, appSpec *kcv1alpha1.AppSpec) error {
	packageObj := kcdatav1alpha1.Package{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "data.packaging.carvel.dev/v1alpha1",
			Kind:       "Package",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", w.Package, w.Version),
		},
		Spec: kcdatav1alpha1.PackageSpec{
			ReleasedAt: metav1.Now(),
			Version:    w.Version,
			RefName:    w.Package,
			Template: kcdatav1alpha1.AppTemplateSpec{
				Spec: appSpec,
			},
		},
	}

	packageBytes, err := yaml.Marshal(packageObj)
	if err != nil {
		return err
	}
	return w.createOrOverwriteFile(path, packageBytes)
}

func (w *ArtifactWriter) createFileIfNotExists(path string, data []byte) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := ioutil.WriteFile(path, data, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	w.ui.PrintHeaderWithContextText("Artefact created", path)
	return nil
}

func (w *ArtifactWriter) createOrOverwriteFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	w.ui.PrintHeaderWithContextText("Artefact created", path)
	return nil
}
