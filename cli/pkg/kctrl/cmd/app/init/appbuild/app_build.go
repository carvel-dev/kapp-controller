// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package appbuild

import (
	"fmt"
	"os"
	"path/filepath"

	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	AppBuildFileName = "app-build.yml"

	AppBuildAPIVersion = "kctrl.carvel.dev/v1alpha1"
	AppBuildKind       = "AppBuild"
)

type AppBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Spec `json:"spec,omitempty"`
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
	Yaml     interface{}      `json:"yaml,omitempty"`
}

type ImgpkgBundle struct {
	Image             string `json:"image,omitempty"`
	UseKbldImagesLock bool   `json:"useKbldImagesLock,omitempty"`
}

type ReleaseResource struct {
}

// Save will persist the appBuild onto the fileSystem. Before saving, it will remove the Annotations from the AppBuild.
func (appBuild AppBuild) Save() error {
	// We dont want to persist the annotations.
	appBuild.ObjectMeta.Annotations = nil
	content, err := yaml.Marshal(appBuild)
	if err != nil {
		return err
	}

	return WriteFile(AppBuildFileName, content)
}

func NewAppBuild() (AppBuild, error) {
	var appBuild AppBuild
	appBuildFilePath := filepath.Join(AppBuildFileName)
	exists, err := IsFileExists(appBuildFilePath)
	if err != nil {
		return AppBuild{}, err
	}

	if exists {
		appBuild, err = NewAppBuildFromFile(appBuildFilePath)
		if err != nil {
			return AppBuild{}, err
		}

		//In case user has manually removed the app section from the app-build
		if appBuild.Spec.App == nil {
			defaultApp := NewDefaultAppTemplateSpec()
			appBuild.Spec.App = defaultApp
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
	if err != nil {
		return AppBuild{}, err
	}
	return appBuild, nil
}

func NewDefaultAppTemplateSpec() *v1alpha1.AppTemplateSpec {
	appSpec := v1alpha12.AppSpec{
		Fetch:    []v1alpha12.AppFetch{},
		Template: []v1alpha12.AppTemplate{},
		Deploy: []v1alpha12.AppDeploy{
			{Kapp: &v1alpha12.AppDeployKapp{}},
		},
	}
	return &v1alpha1.AppTemplateSpec{&appSpec}
}

func NewDefaultAppBuild() AppBuild {
	appBuild := AppBuild{
		TypeMeta: metav1.TypeMeta{
			Kind:       AppBuildKind,
			APIVersion: AppBuildAPIVersion,
		},
		Spec: Spec{
			App: NewDefaultAppTemplateSpec(),
		},
	}
	return appBuild
}

func (appBuild AppBuild) GetAppSpec() *v1alpha12.AppSpec {
	return appBuild.Spec.App.Spec
}

func (appBuild *AppBuild) SetAppSpec(appSpec *v1alpha12.AppSpec) {
	if appBuild.Spec.App == nil {
		appBuild.Spec.App = &v1alpha1.AppTemplateSpec{}
	}
	appBuild.Spec.App.Spec = appSpec
}

func (appBuild AppBuild) GetObjectMeta() *metav1.ObjectMeta {
	return nil
}

func (appBuild *AppBuild) SetObjectMeta(metaObj *metav1.ObjectMeta) {
	appBuild.ObjectMeta = *metaObj
	return
}

func (appBuild AppBuild) GetExport() *[]Export {
	return nil
}

func (appBuild *AppBuild) SetExport(exportObj *[]Export) {
	appBuild.Spec.Export = *exportObj
	return
}

// Check if file exists
func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check for the existence of file. Error is: %s", err.Error())
	}
}

// Write binary content to file
func WriteFile(filePath string, data []byte) error {
	// Create creates or truncates the named file. If the file already exists, it is truncated.
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
