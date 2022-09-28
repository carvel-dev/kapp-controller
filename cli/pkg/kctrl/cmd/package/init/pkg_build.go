// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"os"

	appbuild "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type PackageBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PackageBuildSpec `json:"spec,omitempty"`
}

type PackageBuildSpec struct {
	Template Template           `json:"template,omitempty"`
	Release  []appbuild.Release `json:"release,omitempty"`
}

type Template struct {
	Spec appbuild.Spec `json:"spec"`
}

func (b *PackageBuild) Save() error {
	content, err := yaml.Marshal(b)
	if err != nil {
		return err
	}

	return appbuild.WriteFile(pkgBuildFileName, content)
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

func (b *PackageBuild) GetAppSpec() *v1alpha12.AppSpec {
	return b.getAppSpec()
}

func (b *PackageBuild) getAppSpec() *v1alpha12.AppSpec {
	if b.Spec.Template.Spec.App == nil || b.Spec.Template.Spec.App.Spec == nil {
		return nil
	}
	return b.Spec.Template.Spec.App.Spec
}

func (b *PackageBuild) SetAppSpec(appSpec *v1alpha12.AppSpec) {
	if b.Spec.Template.Spec.App == nil {
		b.Spec.Template.Spec.App = &v1alpha1.AppTemplateSpec{}
	}
	b.Spec.Template.Spec.App.Spec = appSpec
}

func (b *PackageBuild) GetObjectMeta() *metav1.ObjectMeta {
	return &b.ObjectMeta
}

func (b *PackageBuild) SetObjectMeta(metaObj *metav1.ObjectMeta) {
	b.ObjectMeta = *metaObj
	return
}

func (b *PackageBuild) GetExport() *[]appbuild.Export {
	return &b.Spec.Template.Spec.Export
}

func (b *PackageBuild) SetExport(exportObj *[]appbuild.Export) {
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

func (b *PackageBuild) ConfigureExportSection() {
	fetchSource := b.GetObjectMeta().Annotations[appbuild.FetchContentAnnotationKey]
	exportSection := *b.GetExport()
	// In case of pkg init rerun with FetchFromLocalDirectory, today we overwrite the includePaths
	// with what we get from template section.
	// Alternatively, we can merge the includePaths with template section.
	// It becomes complex to merge already existing includePaths with template section especially scenario 2
	// Scenario 1: During rerun, something is added in the app template section
	// Scenario 2: During rerun, something is removed from the app template section
	if exportSection == nil || len(exportSection) == 0 || fetchSource == appbuild.FetchFromLocalDirectory {
		appTemplates := b.GetAppSpec().Template
		includePaths := []string{}
		for _, appTemplate := range appTemplates {
			if appTemplate.HelmTemplate != nil {
				includePaths = append(includePaths, appbuild.UpstreamFolderName)
			}

			if appTemplate.Ytt != nil {
				for _, path := range appTemplate.Ytt.Paths {
					if path == appbuild.StdIn {
						continue
					}
					includePaths = append(includePaths, path)
				}
			}
		}

		if len(exportSection) == 0 {
			exportSection = []appbuild.Export{{}}
		}
		exportSection[0].IncludePaths = includePaths

		b.SetExport(&exportSection)
	}
}
