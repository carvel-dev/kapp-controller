// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package appbuild

import (
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
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
	Resource interface{} `json:"resource,omitempty"`
	Yaml     interface{} `json:"yaml,omitempty"`
}

type ImgpkgBundle struct {
	Image             string `json:"image,omitempty"`
	UseKbldImagesLock bool   `json:"useKbldImagesLock,omitempty"`
}

func (appBuild AppBuild) WriteToFile() error {
	content, err := yaml.Marshal(appBuild)
	if err != nil {
		return err
	}

	return common.WriteFile(AppBuildFileName, content)
}

func GetAppBuild() (AppBuild, error) {
	var appBuild AppBuild
	appBuildFilePath := filepath.Join(AppBuildFileName)
	exists, err := common.IsFileExists(appBuildFilePath)
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
