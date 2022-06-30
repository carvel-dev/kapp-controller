// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"os"

	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	kapp "github.com/k14s/kapp/pkg/kapp/resources"
	appBuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type PackageBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline"`
	Spec              PackageBuildSpec `json:"spec, omitempty"`
}

type PackageBuildSpec struct {
	Template appBuild.AppBuild  `json:"template, omitempty"`
	Release  []appBuild.Release `json:"release, omitempty"`
}

type Template struct {
	Spec TemplateSpec `json:"spec, omitempty"`
}

type TemplateSpec struct {
	Export []appBuild.Export         `json:"export, omitempty"`
	App    *v1alpha1.AppTemplateSpec `json:"app, omitempty"`
}

type ImgpkgBundle struct {
	Image             string   `json:"image, omitempty"`
	UseKbldImagesLock bool     `json:"useKbldImagesLock, omitempty"`
	IncludePaths      []string `json:"includePaths, omitempty"`
}

// func (packageBuild PackageBuild) WriteToFile() error {
// 	content, err := yaml.Marshal(packageBuild)
// 	if err != nil {
// 		return err
// 	}

// 	return common.WriteFile(common.PackageBuildFileName, content)
// }

func GetPackageBuild(pkgBuildFilePath string) (*PackageBuild, error) {
	var packageBuild *PackageBuild
	exists, err := common.IsFileExists(pkgBuildFilePath)
	if err != nil {
		return nil, err
	}

	if exists {
		packageBuild, err = NewPackageBuildFromFile(pkgBuildFilePath)
		if err != nil {
			return nil, err
		}

		//TODO In case user has manually removed the app section from the app-build
		/*if packageBuild.Spec.App == nil {
			defaultApp, err := NewDefaultAppTemplateSpec()
			if err != nil {
				return PackageBuild{}, err
			}
			appBuild.Spec.App = defaultApp
		}*/
	} else {
		// packageBuild, err = NewDefaultPackageBuild()
		// if err != nil {
		// 	return nil, err
		// }
	}

	return packageBuild, nil
}

func NewPackageBuildFromFile(filePath string) (*PackageBuild, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var packageBuild PackageBuild
	err = yaml.Unmarshal(content, &packageBuild)
	if err != nil {
		return nil, err
	}
	return &packageBuild, nil
}

func GetPackage(filePath string) (*v1alpha1.Package, error) {
	resources, err := getResources(filePath)
	if err != nil {
		return nil, err
	}
	var pkg v1alpha1.Package
	for _, resource := range resources {
		if resource.Kind() == "Package" {
			data, err := resource.AsYAMLBytes()
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(data, &pkg)
			if err != nil {
				return nil, err
			}
			return &pkg, nil
		}
	}

	// pkg, err = NewDefaultPackage()
	// if err != nil {
	// 	return nil, err
	// }
	// return &pkg, nil
	return nil, nil
}

func GetPackageMetadata(filePath string) (*v1alpha1.PackageMetadata, error) {
	resources, err := getResources(filePath)
	if err != nil {
		return nil, err
	}
	var pkgMetadata v1alpha1.PackageMetadata
	for _, resource := range resources {
		if resource.Kind() == "PackageMetadata" {
			data, err := resource.AsYAMLBytes()
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(data, &pkgMetadata)
			if err != nil {
				return nil, err
			}
			return &pkgMetadata, nil
		}
	}

	// pkgMetadata, err = NewDefaultPackageMetadata()
	// if err != nil {
	// 	return nil, err
	// }
	// return &pkgMetadata, nil
	return nil, nil
}

func GetPackageInstall(filePath string) (*v1alpha12.PackageInstall, error) {
	resources, err := getResources(filePath)
	if err != nil {
		return nil, err
	}
	var packageInstall v1alpha12.PackageInstall
	for _, resource := range resources {
		if resource.Kind() == "PackageInstall" {
			data, err := resource.AsYAMLBytes()
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(data, &packageInstall)
			if err != nil {
				return nil, err
			}
			return &packageInstall, nil
		}
	}

	// packageInstall, err = NewDefaultPackageInstall()
	// if err != nil {
	// 	return nil, err
	// }
	// return &packageInstall, nil
	return nil, nil
}

func getResources(filePath string) ([]kapp.Resource, error) {
	fileExists, err := common.IsFileExists(filePath)
	if err != nil {
		return nil, err
	}
	if !fileExists {
		return nil, nil
	}
	//TODO should we write our own the struct from kapp? Will it lead to circular dependency?
	var allResources []kapp.Resource
	fileRs, err := kapp.NewFileResources(filePath)
	if err != nil {
		return nil, err
	}

	for _, fileRes := range fileRs {
		resources, err := fileRes.Resources()
		if err != nil {
			return nil, err
		}

		allResources = append(allResources, resources...)
	}

	return allResources, nil
}
