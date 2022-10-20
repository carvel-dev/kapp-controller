// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package buildconfigs

import (
	"os"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	PkgBuildFileName = "package-build.yml"
)

type PackageBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PackageBuildSpec `json:"spec,omitempty"`
}

type PackageBuildSpec struct {
	Template Template  `json:"template,omitempty"`
	Release  []Release `json:"release,omitempty"`
}

type Template struct {
	Spec Spec `json:"spec"`
}

var _ Build = &PackageBuild{}

func (b *PackageBuild) Save() error {
	content, err := yaml.Marshal(b)
	if err != nil {
		return err
	}

	err = os.WriteFile(PkgBuildFileName, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
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

func (b *PackageBuild) GetAppSpec() *kcv1alpha1.AppSpec {
	return b.getAppSpec()
}

func (b *PackageBuild) getAppSpec() *kcv1alpha1.AppSpec {
	if b.Spec.Template.Spec.App == nil || b.Spec.Template.Spec.App.Spec == nil {
		return nil
	}
	return b.Spec.Template.Spec.App.Spec
}

func (b *PackageBuild) setAppSpec(appSpec *kcv1alpha1.AppSpec) {
	if b.Spec.Template.Spec.App == nil {
		b.Spec.Template.Spec.App = &v1alpha1.AppTemplateSpec{}
	}
	b.Spec.Template.Spec.App.Spec = appSpec
}

func (b *PackageBuild) SetAppSpec(appSpec *kcv1alpha1.AppSpec) {
	b.setAppSpec(appSpec)
}

func (b *PackageBuild) GetObjectMeta() *metav1.ObjectMeta {
	return &b.ObjectMeta
}

func (b *PackageBuild) SetObjectMeta(metaObj *metav1.ObjectMeta) {
	b.ObjectMeta = *metaObj
	return
}

func (b *PackageBuild) GetExport() []Export {
	return b.Spec.Template.Spec.Export
}

func (b *PackageBuild) SetExport(exportObj *[]Export) {
	b.Spec.Template.Spec.Export = *exportObj
	return
}

func (b *PackageBuild) HasHelmTemplate() bool {
	appSpec := b.getAppSpec()
	if appSpec == nil || appSpec.Template == nil {
		return false
	}
	appTemplates := appSpec.Template
	for _, appTemplate := range appTemplates {
		if appTemplate.HelmTemplate != nil {
			return true
		}
	}
	return false
}

func (b *PackageBuild) InitializeOrKeepDeploySection() {
	appSpec := b.getAppSpec()
	if appSpec.Deploy == nil {
		appSpec.Deploy = []kcv1alpha1.AppDeploy{{Kapp: &kcv1alpha1.AppDeployKapp{}}}
	}
	b.setAppSpec(appSpec)
}
