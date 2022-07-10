// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"os"

	appBuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PkgBuildAPIVersion = "kctrl.carvel.dev/v1alpha1"
	PkgBuildKind       = "PackageBuild"
)

type PackageBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PackageBuildSpec `json:"spec,omitempty"`
}

type PackageBuildSpec struct {
	Template Template           `json:"template,omitempty"`
	Release  []appBuild.Release `json:"release,omitempty"`
}

type Template struct {
	Spec appBuild.Spec `json:"spec"`
}

func (packageBuild PackageBuild) Save() error {
	content, err := yaml.Marshal(packageBuild)
	if err != nil {
		return err
	}

	return common.WriteFile(PkgBuildFileName, content)
}

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
		packageBuild = &PackageBuild{
			TypeMeta: metav1.TypeMeta{
				Kind:       PkgBuildKind,
				APIVersion: PkgBuildAPIVersion,
			},
		}
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

func (pkgBuild PackageBuild) GetAppSpec() *v1alpha12.AppSpec {
	if pkgBuild.Spec.Template.Spec.App == nil || pkgBuild.Spec.Template.Spec.App.Spec == nil {
		return nil
	}
	return pkgBuild.Spec.Template.Spec.App.Spec
}

func (pkgBuild *PackageBuild) SetAppSpec(appSpec *v1alpha12.AppSpec) {
	if pkgBuild.Spec.Template.Spec.App == nil {
		pkgBuild.Spec.Template.Spec.App = &v1alpha1.AppTemplateSpec{}
	}
	pkgBuild.Spec.Template.Spec.App.Spec = appSpec
}

func (pkgBuild PackageBuild) GetObjectMeta() *metav1.ObjectMeta {
	return &pkgBuild.ObjectMeta
}

func (pkgBuild *PackageBuild) SetObjectMeta(metaObj *metav1.ObjectMeta) {
	pkgBuild.ObjectMeta = *metaObj
	return
}

func (pkgBuild PackageBuild) GetExport() *[]appBuild.Export {
	return &pkgBuild.Spec.Template.Spec.Export
}

func (pkgBuild *PackageBuild) SetExport(exportObj *[]appBuild.Export) {
	pkgBuild.Spec.Template.Spec.Export = *exportObj
	return
}
