// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"path/filepath"

	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/kubeconfig"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
)

// AppFactory allows to create "hidden" Apps for reconciling PackageRepositories.
type AppFactory struct {
	CoreClient   kubernetes.Interface
	AppClient    kcclient.Interface
	TimeMetrics  *metrics.ReconcileTimeMetrics
	CountMetrics *metrics.ReconcileCountMetrics
	KcConfig     *config.Config
	CmdRunner    exec.CmdRunner
	Kubeconf     *kubeconfig.Kubeconfig
	CacheFolder  *memdir.TmpDir
}

// NewCRDPackageRepo constructs "hidden" App to reconcile PackageRepository.
func (f *AppFactory) NewCRDPackageRepo(app *kcv1alpha1.App, pkgr *pkgv1alpha1.PackageRepository, log logr.Logger) *CRDApp {
	vendirOpts := fetch.VendirOpts{
		SkipTLSConfig:   f.KcConfig,
		BaseCacheFolder: filepath.Join(f.CacheFolder.Path(), "pkg-repo"),
	}
	fetchFactory := fetch.NewFactory(f.CoreClient, vendirOpts, f.CmdRunner)
	templateFactory := template.NewFactory(f.CoreClient, fetchFactory, false, f.CmdRunner)
	deployFactory := deploy.NewFactory(f.CoreClient, f.Kubeconf, nil, f.CmdRunner, log)
	return NewCRDApp(app, pkgr, log, f.AppClient, f.TimeMetrics, f.CountMetrics, fetchFactory, templateFactory, deployFactory)
}
