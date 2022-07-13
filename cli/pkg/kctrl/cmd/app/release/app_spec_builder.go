// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/appbuild"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdlocal "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppSpecBuilder struct {
	workingDirectory string
	depsFactory      cmdcore.DepsFactory
	logger           logger.Logger
	ui               cmdcore.AuthoringUI
	opts             AppSpecBuilderOpts
}

type AppSpecBuilderOpts struct {
	BuildTemplate []kcv1alpha1.AppTemplate
	BuildDeploy   []kcv1alpha1.AppDeploy
	BuildExport   []appbuild.Export
	BundleImage   string
	Debug         bool
}

func NewAppSpecBuilder(workingDirectory string, depsFactory cmdcore.DepsFactory, logger logger.Logger, ui cmdcore.AuthoringUI, opts AppSpecBuilderOpts) *AppSpecBuilder {
	return &AppSpecBuilder{workingDirectory, depsFactory, logger, ui, opts}
}

const (
	LockOutputFolder = ".imgpkg"
	LockOutputFile   = "images.yml"
)

func (b *AppSpecBuilder) Build() (kcv1alpha1.AppSpec, error) {
	// In-memory app for building and pushing images
	builderApp := kcv1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kctrl-builder",
			Namespace: "in-memory",
			Annotations: map[string]string{
				"kctrl.carvel.dev/local-fetch-0": ".",
			},
		},
		Spec: kcv1alpha1.AppSpec{
			Fetch: []kcv1alpha1.AppFetch{
				{
					// To be replaced by local fetch
					Git: &kcv1alpha1.AppFetchGit{},
				},
			},
			Template: b.opts.BuildTemplate,
			Deploy:   b.opts.BuildDeploy,
		},
	}
	buildConfigs := cmdlocal.Configs{
		Apps: []kcv1alpha1.App{builderApp},
	}

	// Make lock output directory if it does not exist
	tmpImgpkgFolder := filepath.Join(b.workingDirectory, LockOutputFolder)
	_, err := os.Stat(tmpImgpkgFolder)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(tmpImgpkgFolder, os.ModePerm)
		if err != nil {
			return kcv1alpha1.AppSpec{}, err
		}
	}
	defer os.RemoveAll(filepath.Join(b.workingDirectory, LockOutputFolder))

	// Build images and resolved references using reconciler
	tempImgpkgLockPath := filepath.Join(b.workingDirectory, LockOutputFolder, LockOutputFile)
	cmdRunner := NewReleaseCmdRunner(os.Stdout, b.opts.Debug, tempImgpkgLockPath, b.ui)
	reconciler := cmdlocal.NewReconciler(b.depsFactory, cmdRunner, b.logger)

	err = reconciler.Reconcile(buildConfigs, cmdlocal.ReconcileOpts{
		Local:     true,
		KbldBuild: true,
	})

	bundleURL := ""
	useKbldImagesLock := false
	for _, exportStep := range b.opts.BuildExport {
		switch {
		case exportStep.ImgpkgBundle != nil:
			useKbldImagesLock = exportStep.ImgpkgBundle.UseKbldImagesLock
			imgpkgRunner := ImgpkgRunner{
				BundlePath:        fmt.Sprintf("%s:build-%d", exportStep.ImgpkgBundle.Image, time.Now().Unix()),
				Paths:             exportStep.IncludePaths,
				UseKbldImagesLock: exportStep.ImgpkgBundle.UseKbldImagesLock,
				ImgLockFilepath:   tempImgpkgLockPath,
				UI:                b.ui,
			}
			bundleURL, err = imgpkgRunner.Run()
			if err != nil {
				return kcv1alpha1.AppSpec{}, err
			}

		default:
			continue
		}
	}

	appSpec := kcv1alpha1.AppSpec{
		Fetch: []kcv1alpha1.AppFetch{
			{
				ImgpkgBundle: &kcv1alpha1.AppFetchImgpkgBundle{
					Image: bundleURL,
				},
			},
		},
		Template: b.opts.BuildTemplate,
		Deploy:   b.opts.BuildDeploy,
	}
	if useKbldImagesLock {
		for _, templateStage := range appSpec.Template {
			if templateStage.Kbld != nil {
				templateStage.Kbld.Paths = append(templateStage.Kbld.Paths, "-")
				templateStage.Kbld.Paths = append(templateStage.Kbld.Paths, ".imgpkg/images.yml")
			}
		}
	}

	return appSpec, nil
}