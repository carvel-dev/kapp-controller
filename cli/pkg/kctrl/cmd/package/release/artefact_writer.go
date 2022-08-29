// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcdatav1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type ArtifactWriter struct {
	packageName string
	version     string
	artifactDir string

	metadataTemplate *kcdatav1alpha1.PackageMetadata
	packageTemplate  *kcdatav1alpha1.Package

	ui cmdcore.AuthoringUI
}

const (
	artifactDir = "carvel-artifacts"
	packageDir  = "packages"
)

func NewArtifactWriter(pkg string, version string, packageTemplate *kcdatav1alpha1.Package,
	metadataTemplate *kcdatav1alpha1.PackageMetadata, artifactDir string, ui cmdcore.AuthoringUI) *ArtifactWriter {
	return &ArtifactWriter{packageName: pkg, version: version, artifactDir: artifactDir, metadataTemplate: metadataTemplate,
		packageTemplate: packageTemplate, ui: ui}
}

func (w *ArtifactWriter) Write(appSpec *kcv1alpha1.AppSpec, valuesSchema kcdatav1alpha1.ValuesSchema) error {
	path := filepath.Join(w.artifactDir, packageDir, w.packageName)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = w.writePackageMetadata(filepath.Join(path, "metadata.yml"))
	if err != nil {
		return err
	}
	err = w.writePackage(filepath.Join(path, "package.yml"), appSpec, valuesSchema)
	if err != nil {
		return err
	}

	return nil
}

func (w *ArtifactWriter) WriteRepoOutput(appSpec *kcv1alpha1.AppSpec, valuesSchema kcdatav1alpha1.ValuesSchema, path string) error {
	path = filepath.Join(path, packageDir, w.packageName)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = w.writePackageMetadata(filepath.Join(path, "metadata.yml"))
	if err != nil {
		return err
	}
	err = w.writePackage(filepath.Join(path, fmt.Sprintf("%s.yml", w.version)), appSpec, valuesSchema)
	if err != nil {
		return err
	}

	return nil
}

func (w *ArtifactWriter) writePackageMetadata(path string) error {
	metadataBytes, err := yaml.Marshal(w.metadataTemplate)
	if err != nil {
		return err
	}
	return w.createOrOverwriteFile(path, metadataBytes)
}

func (w *ArtifactWriter) writePackage(path string, appSpec *kcv1alpha1.AppSpec, valuesSchema kcdatav1alpha1.ValuesSchema) error {
	w.packageTemplate.SetName(fmt.Sprintf("%s.%s", w.packageName, w.version))
	w.packageTemplate.Spec.ReleasedAt = metav1.Now()
	w.packageTemplate.Spec.Version = w.version
	w.packageTemplate.Spec.RefName = w.packageName
	w.packageTemplate.Spec.Template = kcdatav1alpha1.AppTemplateSpec{Spec: appSpec}
	w.packageTemplate.Spec.ValuesSchema = valuesSchema

	packageBytes, err := yaml.Marshal(w.packageTemplate)
	if err != nil {
		return err
	}
	return w.createOrOverwriteFile(path, packageBytes)
}

func (w *ArtifactWriter) createOrOverwriteFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	w.ui.PrintHeaderWithContextText("Artifact created", path)
	return nil
}
