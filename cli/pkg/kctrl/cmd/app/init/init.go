// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"os"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/configure/template"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	AppFileName = "app.yml"
)

type InitOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger
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

	return cmd
}

func (o *InitOptions) Run() error {
	o.ui.PrintHeaderText("\nPre-requisite")
	o.ui.PrintInformationalText("Welcome! Before we start on the app creation journey, please ensure the following pre-requites are met:\n* The Carvel suite of tools are installed. Do get familiar with the following Carvel tools: ytt, imgpkg, vendir, and kbld.\n* You have access to an OCI registry, and authenticated locally so that images can be pushed. e.g. docker login <REGISTRY URL>\n")

	appBuild, err := appbuild.NewAppBuild()
	if err != nil {
		return err
	}
	createStep := NewCreateStep(o.ui, &appBuild, o.logger, o.depsFactory, true)
	err = common.Run(createStep)
	if err != nil {
		return err
	}
	printInformation(o.ui)
	printNextStep(o.ui)

	return nil
}

type CreateStep struct {
	ui                        cmdcore.AuthoringUI
	appBuild                  *appbuild.AppBuild
	logger                    logger.Logger
	depsFactory               cmdcore.DepsFactory
	isAppCommandRunExplicitly bool
}

func NewCreateStep(ui cmdcore.AuthoringUI, appBuild *appbuild.AppBuild, logger logger.Logger, depsFactory cmdcore.DepsFactory, isAppCommandRunExplicitly bool) *CreateStep {
	return &CreateStep{
		ui:                        ui,
		appBuild:                  appBuild,
		logger:                    logger,
		depsFactory:               depsFactory,
		isAppCommandRunExplicitly: isAppCommandRunExplicitly,
	}
}

func (createStep CreateStep) GetAppBuild() *appbuild.AppBuild {
	return createStep.appBuild
}

func (createStep CreateStep) printStartBlock() {
}

func (createStep *CreateStep) PreInteract() error {
	createStep.printStartBlock()
	return nil
}

func (createStep *CreateStep) Interact() error {
	err := createStep.configureFetchSection()
	if err != nil {
		return err
	}

	err = createStep.configureTemplateSection()
	if err != nil {
		return err
	}

	createStep.configureExportSection()

	return nil
}

func (createStep *CreateStep) configureFetchSection() error {
	fetchConfiguration := fetch.NewFetchStep(createStep.ui, createStep.appBuild, createStep.isAppCommandRunExplicitly)
	err := common.Run(fetchConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (createStep *CreateStep) configureTemplateSection() error {
	templateConfiguration := template.NewTemplateStep(createStep.ui, createStep.appBuild)
	err := common.Run(templateConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (createStep *CreateStep) PostInteract() error {
	if createStep.isAppCommandRunExplicitly {
		appConfig, err := createStep.generateApp()
		if err != nil {
			return err
		}

		appContent, err := yaml.Marshal(appConfig)
		if err != nil {
			return err
		}
		err = common.WriteFile(AppFileName, appContent)
		if err != nil {
			return err
		}

		createStep.ui.PrintInformationalText(fmt.Sprintf("Both the files can be accessed from the following location:"))
	}

	return nil
}

func (createStep CreateStep) generateApp() (kcv1alpha1.App, error) {
	var app kcv1alpha1.App
	exists, err := common.IsFileExists(AppFileName)
	if err != nil {
		return kcv1alpha1.App{}, err
	}

	if exists {
		app, err = createStep.updateExistingApp()
		if err != nil {
			return kcv1alpha1.App{}, err
		}
	} else {
		app = createStep.createAppFromAppBuild()
	}
	timestamp := v1.NewTime(time.Now().UTC()).Rfc3339Copy()
	app.ObjectMeta.CreationTimestamp = timestamp
	return app, nil

}

func printNextStep(ui cmdcore.AuthoringUI) {
	ui.PrintInformationalText("\n**Next steps**")
	ui.PrintInformationalText("\nCreated app can be consumed in following ways:\n")
}

func printInformation(ui cmdcore.AuthoringUI) {
	ui.PrintInformationalText("\n**Information**\napp-build.yml is generated as part of this flow. This file can be used for further updating and adding complex scenarios while using the `kctrl dev deploy` command. Please read the link'ed documentation for more explanation.")
}

func (createStep CreateStep) createAppFromAppBuild() kcv1alpha1.App {
	//TODO Should we ask for the app Name ?
	appName := "microservices-demo"
	serviceAccountName := fmt.Sprintf("sa-%s", appName)
	appAnnotation := map[string]string{
		fetch.LocalFetchAnnotationKey: ".",
	}
	appTemplateSection := createStep.appBuild.Spec.App.Spec.Template
	//TODO should we remove the fetch section as it is not beind used and add it dynamically during dev deploy.
	appFetchSection := []kcv1alpha1.AppFetch{kcv1alpha1.AppFetch{
		HTTP: &kcv1alpha1.AppFetchHTTP{},
	}}
	appDeploySection := createStep.appBuild.Spec.App.Spec.Deploy
	return kcv1alpha1.App{
		TypeMeta: v1.TypeMeta{
			APIVersion: "kappctrl.k14s.io/v1alpha1",
			Kind:       "App",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:        appName,
			Annotations: appAnnotation,
		},
		Spec: kcv1alpha1.AppSpec{
			ServiceAccountName: serviceAccountName,
			Fetch:              appFetchSection,
			Template:           appTemplateSection,
			Deploy:             appDeploySection,
		},
	}
}

func (createStep CreateStep) updateExistingApp() (kcv1alpha1.App, error) {
	var existingApp kcv1alpha1.App
	data, err := os.ReadFile(AppFileName)
	if err != nil {
		return kcv1alpha1.App{}, err
	}
	err = yaml.Unmarshal(data, &existingApp)
	if err != nil {
		return kcv1alpha1.App{}, err
	}
	fetchSource := createStep.appBuild.ObjectMeta.Annotations[fetch.FetchContentAnnotationKey]
	if fetchSource == fetch.FetchFromLocalDirectory {
		return existingApp, nil
	}

	templateSectionFromExistingApp := existingApp.Spec.Template
	transformAppTemplates := template.NewTransformAppTemplates(&templateSectionFromExistingApp)
	//As fetchSource is not from Local directory, we should add upstream folder in ytt path and helmTemplate(if required).
	if fetchSource == fetch.FetchChartFromHelmRepo || fetchSource == fetch.FetchChartFromGithub {
		//TODO Figure out the edge scenarios
		transformAppTemplates.AddUpstreamAsPathToHelmIfNotExist()
		transformAppTemplates.AddStdInAsPathToYttIfNotExist()

	}
	transformAppTemplates.AddUpstreamAsPathToYttIfNotExist()

	return existingApp, nil
}

func (createStep CreateStep) configureExportSection() {
	fetchSource := createStep.appBuild.ObjectMeta.Annotations[fetch.FetchContentAnnotationKey]
	if fetchSource == fetch.FetchFromLocalDirectory {
		return
	}

	// TODO current implementation is if export section is already defined, we will not touch it. Confirm the same.
	if createStep.appBuild.Spec.Export == nil || len(createStep.appBuild.Spec.Export) == 0 {
		createStep.appBuild.Spec.Export = append(createStep.appBuild.Spec.Export, appbuild.Export{
			ImgpkgBundle: nil,
			IncludePaths: []string{template.UpstreamFolderName},
		})
	}
	return
}
