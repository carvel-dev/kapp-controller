// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sources

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	buildconfigs "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/buildconfigs"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

type TemplateConfiguration struct {
	ui    cmdcore.AuthoringUI
	build buildconfigs.Build
}

func NewTemplateConfiguration(ui cmdcore.AuthoringUI, build buildconfigs.Build) *TemplateConfiguration {
	return &TemplateConfiguration{ui: ui, build: build}
}

func (t *TemplateConfiguration) Configure(fetchMode string) error {
	appSpec := t.build.GetAppSpec()
	if appSpec == nil {
		appSpec = &v1alpha1.AppSpec{}
	}
	existingTemplates := appSpec.Template

	/* In case of pkg init rerun, if user has selected anything else except FetchFromLocalDirectory,
	we will return from Template section without touching it.
	We dont want to reset the modification user have done. */
	if len(existingTemplates) > 0 {
		if fetchMode == LocalDirectory {
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
		if fetchMode == ChartFromGit || fetchMode == HelmRepo {
			appTemplate, err := t.getHelmAppTemplate(fetchMode)
			if err != nil {
				return err
			}
			appTemplates = append(appTemplates, appTemplate)
		}

		//  Define YttPaths
		var defaultYttPaths []string
		if fetchMode == HelmRepo || fetchMode == ChartFromGit {
			defaultYttPaths = []string{buildconfigs.StdIn}
		} else if fetchMode == LocalDirectory {
			var err error
			defaultYttPaths, err = t.getYttPathsForLocalDirectory("")
			if err != nil {
				return err
			}
		} else {
			defaultYttPaths = []string{VendirSyncDirectory}
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

func (t *TemplateConfiguration) getHelmAppTemplate(fetchMode string) (v1alpha1.AppTemplate, error) {
	var pathFromVendir string
	if fetchMode == ChartFromGit {
		vendirConfig := NewVendirConfig(vendirFileName)
		err := vendirConfig.Load()
		if err != nil {
			return v1alpha1.AppTemplate{}, err
		}
		fmt.Println(vendirConfig.Config)
		pathFromVendir = vendirConfig.Contents()[0].IncludePaths[0]
		// Remove all the trailing `/` from the string
		pathFromVendir = strings.TrimRight(pathFromVendir, "/")
		pathFromVendir = strings.TrimSuffix(pathFromVendir, "/**/*")
		pathFromVendir = strings.TrimSuffix(pathFromVendir, "/**")
		pathFromVendir = strings.TrimSuffix(pathFromVendir, "/*")
	}
	appTemplateWithHelm := v1alpha1.AppTemplate{
		HelmTemplate: &v1alpha1.AppTemplateHelmTemplate{
			Path: filepath.Join(VendirSyncDirectory, pathFromVendir),
		}}
	return appTemplateWithHelm, nil
}

func (t *TemplateConfiguration) getYttPathsForLocalDirectory(defaultIncludedPath string) ([]string, error) {
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
