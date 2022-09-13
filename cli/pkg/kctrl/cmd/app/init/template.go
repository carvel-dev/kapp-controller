// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	UpstreamFolderName = "upstream"
	StdIn              = "-"
)

type TemplateStep struct {
	ui    cmdcore.AuthoringUI
	build Build
}

func NewTemplateStep(ui cmdcore.AuthoringUI, build Build) *TemplateStep {
	templateStep := TemplateStep{
		ui:    ui,
		build: build,
	}
	return &templateStep
}

func (t *TemplateStep) PreInteract() error { return nil }

func (t *TemplateStep) Interact() error {
	appSpec := t.build.GetAppSpec()
	if appSpec == nil {
		appSpec = &v1alpha1.AppSpec{}
	}
	existingTemplates := appSpec.Template
	fetchSource := t.build.GetObjectMeta().Annotations[FetchContentAnnotationKey]

	/* In case of pkg init rerun, if user has selected anything else except FetchFromLocalDirectory,
	we will return from Template section without touching it.
	We dont want to reset the modification user have done. */
	if len(existingTemplates) > 0 {
		if fetchSource == FetchFromLocalDirectory {
			for _, template := range existingTemplates {
				if template.Ytt != nil {
					var defaultIncludedPath string
					defaultIncludedPath = strings.Join(template.Ytt.Paths, ",")
					defaultYttPaths, err := t.getYttPathsForLocalDirectory(defaultIncludedPath)
					if err != nil {
						return err
					}
					template.Ytt.Paths = defaultYttPaths
					t.build.SetAppSpec(appSpec)
					return t.build.Save()
				}
			}
		} else {
			return nil
		}
	} else {
		appTemplates := []v1alpha1.AppTemplate{}

		// Add helmTemplate
		if fetchSource == FetchChartFromGit || fetchSource == FetchFromHelmRepo {
			appTemplate, err := t.getHelmAppTemplate(fetchSource)
			if err != nil {
				return err
			}
			appTemplates = append(appTemplates, appTemplate)
		}

		//  Define YttPaths
		var defaultYttPaths []string
		if fetchSource == FetchFromHelmRepo || fetchSource == FetchChartFromGit {
			defaultYttPaths = []string{StdIn}
		} else if fetchSource == FetchFromLocalDirectory {
			var err error
			defaultYttPaths, err = t.getYttPathsForLocalDirectory("")
			if err != nil {
				return err
			}
		} else {
			defaultYttPaths = []string{UpstreamFolderName}
		}
		// Add yttTemplate
		appTemplateWithYtt := v1alpha1.AppTemplate{
			Ytt: &v1alpha1.AppTemplateYtt{
				Paths: defaultYttPaths,
			},
		}
		appTemplates = append(appTemplates, appTemplateWithYtt)

		// Add kbldTemplate
		appTemplates = append(appTemplates, v1alpha1.AppTemplate{Kbld: &v1alpha1.AppTemplateKbld{}})

		appSpec.Template = appTemplates
		t.build.SetAppSpec(appSpec)
		return t.build.Save()
	}
	return nil
}

func (t *TemplateStep) getHelmAppTemplate(fetchSource string) (v1alpha1.AppTemplate, error) {
	var pathFromVendir string
	if fetchSource == FetchChartFromGit {
		vendirConf, err := ReadVendirConfig()
		if err != nil {
			return v1alpha1.AppTemplate{}, err
		}
		pathFromVendir = vendirConf.Directories[0].Contents[0].IncludePaths[0]
		// Remove all the trailing `/` from the string
		pathFromVendir = strings.TrimRight(pathFromVendir, "/")
		pathFromVendir = strings.TrimSuffix(pathFromVendir, "/**/*")
		pathFromVendir = strings.TrimSuffix(pathFromVendir, "/**")
		pathFromVendir = strings.TrimSuffix(pathFromVendir, "/*")
	}
	appTemplateWithHelm := v1alpha1.AppTemplate{
		HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{
			Path: filepath.Join(UpstreamFolderName, pathFromVendir),
		}}
	return appTemplateWithHelm, nil
}

func (t *TemplateStep) getYttPathsForLocalDirectory(defaultIncludedPath string) ([]string, error) {
	t.ui.PrintInformationalText("We need to include files/ directories which contain Kubernetes manifests. " +
		"Multiple values can be included using a comma separator.")
	textOpts := ui.TextOpts{
		Label:        "Enter the paths which contain Kubernetes manifests",
		Default:      defaultIncludedPath,
		ValidateFunc: nil,
	}
	includePaths, err := t.ui.AskForText(textOpts)
	if err != nil {
		return nil, err
	}
	defaultYttPaths := strings.Split(includePaths, ",")
	for i := range defaultYttPaths {
		defaultYttPaths[i] = strings.TrimSpace(defaultYttPaths[i])
	}
	return defaultYttPaths, nil
}

func (t *TemplateStep) PostInteract() error { return nil }
