// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	cmdlocal "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	buildconfigs "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local/buildconfigs"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	fakekc "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppSpecBuilder struct {
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger
	ui          cmdcore.AuthoringUI
	opts        AppSpecBuilderOpts
}

type AppSpecBuilderOpts struct {
	BuildTemplate []kcv1alpha1.AppTemplate
	BuildDeploy   []kcv1alpha1.AppDeploy
	BuildExport   []buildconfigs.Export
	BundleImage   string
	Debug         bool
	BundleTag     string
}

func NewAppSpecBuilder(depsFactory cmdcore.DepsFactory, logger logger.Logger, ui cmdcore.AuthoringUI, opts AppSpecBuilderOpts) *AppSpecBuilder {
	return &AppSpecBuilder{depsFactory, logger, ui, opts}
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
	tmpImgpkgFolder := LockOutputFolder
	_, err := os.Stat(tmpImgpkgFolder)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(tmpImgpkgFolder, os.ModePerm)
		if err != nil {
			return kcv1alpha1.AppSpec{}, err
		}
	}
	defer os.RemoveAll(LockOutputFolder)

	// Build images and resolve references using reconciler
	tempImgpkgLockPath := filepath.Join(LockOutputFolder, LockOutputFile)
	cmdRunner := NewReleaseCmdRunner(os.Stdout, b.opts.Debug, tempImgpkgLockPath, b.ui)
	reconciler := cmdlocal.NewReconciler(b.depsFactory, cmdRunner, b.logger)

	err = reconciler.Reconcile(buildConfigs, cmdlocal.ReconcileOpts{
		Local:             true,
		KbldBuild:         true,
		AfterAppReconcile: b.checkForErrorsAfterReconciliation,
	})
	if err != nil {
		return kcv1alpha1.AppSpec{}, err
	}

	bundleURL := ""
	useKbldImagesLock := false
	tag := fmt.Sprintf("build-%d", time.Now().Unix())
	if b.opts.BundleTag != "" {
		tag = b.opts.BundleTag
	}
	for _, exportStep := range b.opts.BuildExport {
		switch {
		case exportStep.ImgpkgBundle != nil:
			useKbldImagesLock = exportStep.ImgpkgBundle.UseKbldImagesLock
			imgpkgRunner := ImgpkgRunner{
				BundlePath:        fmt.Sprintf("%s:%s", exportStep.ImgpkgBundle.Image, tag),
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
				templateStage.Kbld.Paths = []string{"-", filepath.Join(LockOutputFolder, LockOutputFile)}
			}
		}
	}

	return appSpec, nil
}

func (b *AppSpecBuilder) checkForErrorsAfterReconciliation(app kcv1alpha1.App, fakeClient *fakekc.Clientset) error {
	existingApp, err := fakeClient.KappctrlV1alpha1().Apps(app.Namespace).Get(context.Background(), app.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// TODO: How can we prevent reconciler from trying to prepare kapp?
	if existingApp.Status.UsefulErrorMessage != "" && !strings.Contains(existingApp.Status.UsefulErrorMessage, "Preparing kapp") {
		return fmt.Errorf("Reconciling: %s", existingApp.Status.UsefulErrorMessage)
	}
	return nil
}
