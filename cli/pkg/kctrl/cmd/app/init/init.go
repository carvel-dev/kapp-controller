// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	AppFileName = "app.yml"
)

type InitOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger
	chdir       string
}

func NewInitOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *InitOptions {
	return &InitOptions{ui: cmdcore.NewAuthoringUIImpl(ui), depsFactory: depsFactory, logger: logger}
}

func NewInitCmd(o *InitOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize App",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Working directory with app-build and other config")
	return cmd
}

func (o *InitOptions) Run() error {
	if o.chdir != "" {
		err := os.Chdir(o.chdir)
		if err != nil {
			return err
		}
	}

	o.ui.PrintHeaderText("\nPre-requisite")
	o.ui.PrintInformationalText("Welcome! Before we start on the app creation journey, please ensure the following pre-requites are met:\n* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n* You have access to an OCI registry, and authenticated locally so that images can be pushed. e.g. docker login <REGISTRY URL>\n")

	appBuild, err := NewAppBuild()
	if err != nil {
		return err
	}
	createStep := NewCreateStep(o.ui, &appBuild, o.logger, o.depsFactory, true)
	err = Run(createStep)
	if err != nil {
		return err
	}

	return nil
}

type CreateStep struct {
	ui                        cmdcore.AuthoringUI
	build                     Build
	logger                    logger.Logger
	depsFactory               cmdcore.DepsFactory
	isAppCommandRunExplicitly bool
}

func NewCreateStep(ui cmdcore.AuthoringUI, build Build, logger logger.Logger, depsFactory cmdcore.DepsFactory, isAppCommandRunExplicitly bool) *CreateStep {
	return &CreateStep{
		ui:                        ui,
		build:                     build,
		logger:                    logger,
		depsFactory:               depsFactory,
		isAppCommandRunExplicitly: isAppCommandRunExplicitly,
	}
}

func (c CreateStep) GetAppBuild() Build {
	return c.build
}

func (c *CreateStep) PreInteract() error { return nil }

func (c *CreateStep) Interact() error {
	if c.isAppCommandRunExplicitly {
		c.ui.PrintHeaderText("\nBasic Information")
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		defaultAppName := filepath.Base(wd)
		appBuildObjectMeta := c.build.GetObjectMeta()
		if appBuildObjectMeta == nil {
			appBuildObjectMeta = &metav1.ObjectMeta{}
		}
		if len(appBuildObjectMeta.Name) != 0 {
			defaultAppName = appBuildObjectMeta.Name
		}

		textOpts := ui.TextOpts{
			Label:        "Enter the app name",
			Default:      defaultAppName,
			ValidateFunc: nil,
		}
		appName, err := c.ui.AskForText(textOpts)
		if err != nil {
			return err
		}
		appBuildObjectMeta.Name = appName
		c.build.SetObjectMeta(appBuildObjectMeta)
		err = c.build.Save()
		if err != nil {
			return err
		}
	}

	fetchConfiguration := NewFetchStep(c.ui, c.build, c.isAppCommandRunExplicitly)
	err := Run(fetchConfiguration)
	if err != nil {
		return err
	}

	templateConfiguration := NewTemplateStep(c.ui, c.build)
	err = Run(templateConfiguration)
	if err != nil {
		return err
	}

	return nil
}

func (c *CreateStep) PostInteract() error {
	appSpec := c.build.GetAppSpec()
	if appSpec.Deploy == nil {
		appSpec.Deploy = []kcv1alpha1.AppDeploy{kcv1alpha1.AppDeploy{Kapp: &kcv1alpha1.AppDeployKapp{}}}
	}
	c.build.SetAppSpec(appSpec)
	c.configureExportSection()

	if c.isAppCommandRunExplicitly {
		appConfig, err := c.generateApp()
		if err != nil {
			return err
		}

		appContent, err := yaml.Marshal(appConfig)
		if err != nil {
			return err
		}
		err = WriteFile(AppFileName, appContent)
		if err != nil {
			return err
		}
		c.ui.PrintHeaderText("\nOutput")
		c.ui.PrintInformationalText("Successfully updated app-build.yml\n")
		c.ui.PrintInformationalText("Successfully updated app.yml\n")
		c.ui.PrintHeaderText("\n**Next steps**")
		c.ui.PrintInformationalText("Created files can be consumed in following ways:\n1. Optionally, use 'kctrl dev deploy' to iterate on the app and deploy locally.\n2. Use 'kctrl app release' to release the app.\n")
	}

	return nil
}

func (c CreateStep) generateApp() (kcv1alpha1.App, error) {
	var app kcv1alpha1.App
	exists, err := IsFileExists(AppFileName)
	if err != nil {
		return kcv1alpha1.App{}, err
	}

	if exists {
		data, err := os.ReadFile(AppFileName)
		if err != nil {
			return kcv1alpha1.App{}, err
		}
		err = yaml.Unmarshal(data, &app)
		if err != nil {
			return kcv1alpha1.App{}, err
		}
	} else {
		app = c.createAppFromAppBuild()
	}
	return app, nil

}

func (c CreateStep) createAppFromAppBuild() kcv1alpha1.App {
	appName := "microservices-demo"
	serviceAccountName := fmt.Sprintf("%s-sa", appName)
	appAnnotation := map[string]string{
		LocalFetchAnnotationKey: ".",
	}
	appTemplateSection := c.build.GetAppSpec().Template
	appDeploySection := c.build.GetAppSpec().Deploy
	return kcv1alpha1.App{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kappctrl.k14s.io/v1alpha1",
			Kind:       "App",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Annotations: appAnnotation,
		},
		Spec: kcv1alpha1.AppSpec{
			ServiceAccountName: serviceAccountName,
			Template:           appTemplateSection,
			Deploy:             appDeploySection,
		},
	}
}

func (c CreateStep) configureExportSection() {
	fetchSource := c.build.GetObjectMeta().Annotations[FetchContentAnnotationKey]
	exportSection := *c.build.GetExport()
	// In case of pkg init rerun with FetchFromLocalDirectory, today we overwrite the includePaths
	// with what we get from template section.
	// Alternatively, we can merge the includePaths with template section.
	// It becomes complex to merge already existing includePaths with template section especially scenario 2
	// Scenario 1: During rerun, something is added in the app template section
	// Scenario 2: During rerun, something is removed from the app template section
	if exportSection == nil || len(exportSection) == 0 || fetchSource == FetchFromLocalDirectory {
		appTemplates := c.build.GetAppSpec().Template
		includePaths := []string{}
		for _, appTemplate := range appTemplates {
			if appTemplate.HelmTemplate != nil {
				includePaths = append(includePaths, UpstreamFolderName)
			}

			if appTemplate.Ytt != nil {
				for _, path := range appTemplate.Ytt.Paths {
					if path == StdIn {
						continue
					}
					includePaths = append(includePaths, path)
				}
			}
		}

		if len(exportSection) == 0 {
			exportSection = []Export{Export{}}
		}
		exportSection[0].IncludePaths = includePaths

		c.build.SetExport(&exportSection)
	}
	return
}
