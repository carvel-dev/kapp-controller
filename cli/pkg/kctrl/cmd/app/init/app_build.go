// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"os"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	FileName = "app-build.yml"
)

type AppBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Spec `json:"spec,omitempty"`
}

type Spec struct {
	Export  []Export                  `json:"export,omitempty"`
	App     *v1alpha1.AppTemplateSpec `json:"app,omitempty"`
	Release []Release                 `json:"release,omitempty"`
}

type Export struct {
	ImgpkgBundle *ImgpkgBundle `json:"imgpkgBundle,omitempty"`
	IncludePaths []string      `json:"includePaths,omitempty"`
}

type Release struct {
	Resource *ReleaseResource `json:"resource,omitempty"`
	YAML     *ReleaseYAML     `json:"yaml,omitempty"`
}

type ReleaseYAML struct {
}

type ImgpkgBundle struct {
	Image             string `json:"image,omitempty"`
	UseKbldImagesLock bool   `json:"useKbldImagesLock,omitempty"`
}

type ReleaseResource struct {
}

// Save will persist the appBuild onto the fileSystem. Before saving, it will remove the Annotations from the AppBuild.
func (b AppBuild) Save() error {
	// We dont want to persist the annotations.
	b.ObjectMeta.Annotations = nil
	content, err := yaml.Marshal(b)
	if err != nil {
		return err
	}

	return WriteFile(FileName, content)
}

func NewAppBuild() (AppBuild, error) {
	var appBuild AppBuild
	exists, err := IsFileExists(FileName)
	if err != nil {
		return AppBuild{}, err
	}

	if exists {
		appBuild, err = NewAppBuildFromFile(FileName)
		if err != nil {
			return AppBuild{}, err
		}

		// In case user has manually removed the app section from the app-build
		if appBuild.Spec.App == nil {
			appBuild.Spec.App = NewDefaultAppTemplateSpec()
		}
	} else {
		appBuild = NewDefaultAppBuild()
	}

	return appBuild, nil
}

func NewAppBuildFromFile(filePath string) (AppBuild, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return AppBuild{}, err
	}
	var appBuild AppBuild
	err = yaml.Unmarshal(content, &appBuild)
	return appBuild, err
}

func NewDefaultAppTemplateSpec() *v1alpha1.AppTemplateSpec {
	return &v1alpha1.AppTemplateSpec{
		Spec: &kcv1alpha1.AppSpec{
			Fetch:    []kcv1alpha1.AppFetch{},
			Template: []kcv1alpha1.AppTemplate{},
			Deploy: []kcv1alpha1.AppDeploy{
				{Kapp: &kcv1alpha1.AppDeployKapp{}},
			},
		},
	}
}

func NewDefaultAppBuild() AppBuild {
	return AppBuild{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AppBuild",
			APIVersion: "kctrl.carvel.dev/v1alpha1",
		},
		Spec: Spec{
			App: NewDefaultAppTemplateSpec(),
		},
	}
}

func (b AppBuild) GetAppSpec() *kcv1alpha1.AppSpec {
	return b.Spec.App.Spec
}

func (b *AppBuild) SetAppSpec(appSpec *kcv1alpha1.AppSpec) {
	if b.Spec.App == nil {
		b.Spec.App = &v1alpha1.AppTemplateSpec{}
	}
	b.Spec.App.Spec = appSpec
}

func (b AppBuild) GetObjectMeta() *metav1.ObjectMeta {
	return &b.ObjectMeta
}

func (b *AppBuild) SetObjectMeta(metaObj *metav1.ObjectMeta) {
	b.ObjectMeta = *metaObj
	return
}

func (b AppBuild) GetExport() *[]Export {
	return &b.Spec.Export
}

func (b *AppBuild) SetExport(exportObj *[]Export) {
	b.Spec.Export = *exportObj
	return
}
