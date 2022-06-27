// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	AppBuildFileName = "app-build.yml"
)

type AppBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Spec `json:"spec, omitempty"`
}

type Spec struct {
	Export  []Export                  `json:"export, omitempty"`
	App     *v1alpha1.AppTemplateSpec `json:"app, omitempty"`
	Release []Release                 `json:"release, omitempty"`
}

type Export struct {
	ImgpkgBundle *ImgpkgBundle `json:"imgpkg_bundle, omitempty"`
}

type Release struct {
	Resource interface{} `json:"resource, omitempty"`
	Yaml     interface{} `json:"yaml, omitempty"`
}

type ImgpkgBundle struct {
	Image             string   `json:"image, omitempty"`
	UseKbldImagesLock bool     `json:"useKbldImagesLock, omitempty"`
	IncludePaths      []string `json:"includePaths, omitempty"`
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
			defaultApp, err := NewDefaultAppTemplateSpec()
			if err != nil {
				return AppBuild{}, err
			}
			appBuild.Spec.App = defaultApp
		}
	} else {
		appBuild, err = NewDefaultAppBuild()
		if err != nil {
			return AppBuild{}, err
		}
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
