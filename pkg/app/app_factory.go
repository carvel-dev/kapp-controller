// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	vendirconf "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/config"
	"k8s.io/client-go/kubernetes"
)

// CRDAppFactory allows to create CRDApps.
type CRDAppFactory struct {
	CoreClient       kubernetes.Interface
	AppClient        kcclient.Interface
	KcConfig         *config.Config
	AppMetrics       *metrics.AppMetrics
	VendirConfigHook func(vendirconf.Config) vendirconf.Config
	KbldAllowBuild   bool
	CmdRunner        exec.CmdRunner
	Kubeconf         *kubeconfig.Kubeconfig
	CompInfo         ComponentInfo
}

// NewCRDApp creates a CRDApp injecting necessary dependencies.
func (f *CRDAppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) *CRDApp {
	vendirOpts := fetch.VendirOpts{
		SkipTLSConfig: f.KcConfig,
		ConfigHook:    f.VendirConfigHook,
	}

	fetchFactory := fetch.NewFactory(f.CoreClient, vendirOpts, f.CmdRunner)
	templateFactory := template.NewFactory(f.CoreClient, fetchFactory, f.KbldAllowBuild, f.CmdRunner)
	deployFactory := deploy.NewFactory(f.CoreClient, f.Kubeconf, f.KcConfig, f.CmdRunner, log)

	return NewCRDApp(app, log, f.AppMetrics, f.AppClient, fetchFactory, templateFactory, deployFactory, f.CompInfo, Opts{
		DefaultSyncPeriod: f.KcConfig.AppDefaultSyncPeriod(),
		MinimumSyncPeriod: f.KcConfig.AppMinimumSyncPeriod(),
	})
}
