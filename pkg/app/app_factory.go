// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"path/filepath"

	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "carvel.dev/kapp-controller/pkg/client/clientset/versioned"
	"carvel.dev/kapp-controller/pkg/config"
	"carvel.dev/kapp-controller/pkg/deploy"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/fetch"
	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"carvel.dev/kapp-controller/pkg/memdir"
	"carvel.dev/kapp-controller/pkg/metrics"
	"carvel.dev/kapp-controller/pkg/template"
	vendirconf "carvel.dev/vendir/pkg/vendir/config"
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

// CRDAppFactory allows to create CRDApps.
type CRDAppFactory struct {
	CoreClient       kubernetes.Interface
	AppClient        kcclient.Interface
	KcConfig         *config.Config
	AppMetrics       *metrics.Metrics
	VendirConfigHook func(vendirconf.Config) vendirconf.Config
	KbldAllowBuild   bool
	CmdRunner        exec.CmdRunner
	Kubeconf         *kubeconfig.Kubeconfig
	CompInfo         ComponentInfo
	DeployFactory    deploy.Factory
	CacheFolder      *memdir.TmpDir
}

// NewCRDApp creates a CRDApp injecting necessary dependencies.
func (f *CRDAppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) *CRDApp {
	vendirOpts := fetch.VendirOpts{
		SkipTLSConfig:   f.KcConfig,
		ConfigHook:      f.VendirConfigHook,
		BaseCacheFolder: filepath.Join(f.CacheFolder.Path(), "apps"),
	}

	fetchFactory := fetch.NewFactory(f.CoreClient, vendirOpts, f.CmdRunner)
	templateFactory := template.NewFactory(f.CoreClient, fetchFactory, f.KbldAllowBuild, f.CmdRunner)
	deployFactory := deploy.NewFactory(f.CoreClient, f.Kubeconf, f.KcConfig, f.CmdRunner, log)

	return NewCRDApp(app, log, f.AppMetrics, f.AppClient, fetchFactory, templateFactory, deployFactory, f.CompInfo, Opts{
		DefaultSyncPeriod: f.KcConfig.AppDefaultSyncPeriod(),
		MinimumSyncPeriod: f.KcConfig.AppMinimumSyncPeriod(),
	})
}
