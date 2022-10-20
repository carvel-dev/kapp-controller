// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	buildconfigs "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/buildconfigs"
	sources "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/sources"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	AppFileName = "app.yml"
	//TODO: Where should this constant reside, it is also used by package init
	LocalFetchAnnotationKey = "kctrl.carvel.dev/local-fetch-0"
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

	appBuild, err := buildconfigs.NewAppBuild()
	if err != nil {
		return err
	}

	err = o.getAppBuildName(appBuild)
	if err != nil {
		return err
	}

	err = o.writeAppFile(appBuild)
	if err != nil {
		return err
	}

	fetchMode, sourceConfiguration, err := sources.NewSource(o.ui, appBuild).Configure()
	if err != nil {
		return err
	}

	// Source does not need to be conifgured if manifests are in the local directory
	if sourceConfiguration != nil {
		err = sourceConfiguration.Configure()
		if err != nil {
			return err
		}
	}

	err = sources.NewVendirRunner(o.ui).Sync()
	if err != nil {
		return err
	}

	err = sources.NewTemplateConfiguration(o.ui, appBuild).Configure(fetchMode)
	if err != nil {
		return err
	}

	appSpec := appBuild.GetAppSpec()
	if appSpec.Deploy == nil {
		appSpec.Deploy = []kcv1alpha1.AppDeploy{{Kapp: &kcv1alpha1.AppDeployKapp{}}}
	}
	appBuild.SetAppSpec(appSpec)

	buildconfigs.ConfigureExportSection(appBuild, fetchMode == sources.LocalDirectory, sources.VendirSyncDirectory)
	err = appBuild.Save()
	if err != nil {
		return err
	}
	return nil
}

func (o *InitOptions) getAppBuildName(appBuild *buildconfigs.AppBuild) error {
	o.ui.PrintHeaderText("\nBasic Information")
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	defaultAppName := filepath.Base(wd)
	appBuildObjectMeta := appBuild.GetObjectMeta()
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
	appName, err := o.ui.AskForText(textOpts)
	if err != nil {
		return err
	}
	appBuildObjectMeta.Name = appName
	appBuild.SetObjectMeta(appBuildObjectMeta)

	err = appBuild.Save()
	if err != nil {
		return err
	}
	return nil
}

func (o *InitOptions) writeAppFile(appBuild *buildconfigs.AppBuild) error {
	appConfig, err := o.generateApp(appBuild)
	if err != nil {
		return err
	}

	appContent, err := yaml.Marshal(appConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(AppFileName, appContent, os.ModePerm)
	if err != nil {
		return err
	}
	o.ui.PrintHeaderText("\nOutput")
	o.ui.PrintInformationalText("Successfully updated app-build.yml\n")
	o.ui.PrintInformationalText("Successfully updated app.yml\n")
	o.ui.PrintHeaderText("\n**Next steps**")
	o.ui.PrintInformationalText("Created files can be consumed in following ways:\n1. Optionally, use 'kctrl dev deploy' to iterate on the app and deploy locally.\n2. Use 'kctrl app release' to release the app.\n")
	return nil
}

func (o *InitOptions) generateApp(appBuild *buildconfigs.AppBuild) (kcv1alpha1.App, error) {
	var app kcv1alpha1.App
	_, err := os.Stat(AppFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return o.createAppFromAppBuild(appBuild), nil
		}
		return kcv1alpha1.App{}, err
	}

	data, err := os.ReadFile(AppFileName)
	if err != nil {
		return kcv1alpha1.App{}, err
	}
	err = yaml.Unmarshal(data, &app)
	if err != nil {
		return kcv1alpha1.App{}, err
	}

	return app, nil

}

func (o *InitOptions) createAppFromAppBuild(appBuild *buildconfigs.AppBuild) kcv1alpha1.App {
	appName := "microservices-demo"
	serviceAccountName := fmt.Sprintf("%s-sa", appName)
	appAnnotation := map[string]string{
		LocalFetchAnnotationKey: ".",
	}
	appTemplateSection := appBuild.GetAppSpec().Template
	appDeploySection := appBuild.GetAppSpec().Deploy
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
