package template

import (
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"strings"
)

type TransformAppTemplates struct {
	appTemplates *[]kcv1alpha1.AppTemplate
}

const (
	StdIn              = "-"
	UpstreamFolderName = "upstream"
)

func NewTransformAppTemplates(appTemplates *[]kcv1alpha1.AppTemplate) TransformAppTemplates {
	return TransformAppTemplates{appTemplates: appTemplates}
}

func (t TransformAppTemplates) AddUpstreamAsPathToYttIfNotExist() {
	for _, appTemplate := range *t.appTemplates {
		if appTemplate.Ytt != nil {
			for _, path := range appTemplate.Ytt.Paths {
				if strings.HasPrefix(path, UpstreamFolderName) {
					return
				}
			}
		}
	}
	appTemplateWithYtt := kcv1alpha1.AppTemplate{
		Ytt: &kcv1alpha1.AppTemplateYtt{
			Paths: []string{UpstreamFolderName},
		},
	}
	*t.appTemplates = append([]kcv1alpha1.AppTemplate{appTemplateWithYtt}, *t.appTemplates...)
}

func (t TransformAppTemplates) AddStdInAsPathToYttIfNotExist() {
	for _, appTemplate := range *t.appTemplates {
		if appTemplate.Ytt != nil {
			for _, path := range appTemplate.Ytt.Paths {
				if path == StdIn {
					return
				}
			}
		}
	}
	appTemplateWithYtt := kcv1alpha1.AppTemplate{
		Ytt: &kcv1alpha1.AppTemplateYtt{
			Paths: []string{StdIn},
		},
	}

	//YttTemplate with Stdin should be the immediate next template to the helmTemplate.
	index := 0
	var appTemplate kcv1alpha1.AppTemplate
	if len(*t.appTemplates) == 1 {
		index++
	} else {
		for index, appTemplate = range *t.appTemplates {
			if appTemplate.HelmTemplate == nil {
				break
			}
		}
	}

	//var oldAppTemplates []kcv1alpha1.AppTemplate
	oldAppTemplates := make([]kcv1alpha1.AppTemplate, len(*t.appTemplates))
	copy(oldAppTemplates, *t.appTemplates)
	*t.appTemplates = append(oldAppTemplates[:index], appTemplateWithYtt)
	*t.appTemplates = append(*t.appTemplates, oldAppTemplates[index:]...)

	return
}

func (t TransformAppTemplates) AddUpstreamAsPathToHelmIfNotExist() {
	for _, appTemplate := range *t.appTemplates {
		if appTemplate.HelmTemplate != nil {
			path := appTemplate.HelmTemplate.Path
			//If a helmTemplate exist, it will always be the first one in the template section. Theoretically, it can exist anywhere but every real use case needs it to be first.
			//TODO confirm above understanding. Not handled the scenario if helmTemplate exists as not the first element of slice.
			if strings.HasPrefix(path, UpstreamFolderName) {
				return
			}
		}
	}

	appTemplateWithHelm := kcv1alpha1.AppTemplate{
		HelmTemplate: &kcv1alpha1.AppTemplateHelmTemplate{
			Path: UpstreamFolderName,
		},
	}
	*t.appTemplates = append([]kcv1alpha1.AppTemplate{appTemplateWithHelm}, *t.appTemplates...)
}

func (t TransformAppTemplates) GetAppTemplates() []kcv1alpha1.AppTemplate {
	return *t.appTemplates
}
