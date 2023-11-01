// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"path/filepath"

	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "carvel.dev/kapp-controller/pkg/client/clientset/versioned"
	"carvel.dev/kapp-controller/pkg/config"
	"carvel.dev/kapp-controller/pkg/deploy"
	"carvel.dev/kapp-controller/pkg/exec"
	"carvel.dev/kapp-controller/pkg/fetch"
	"carvel.dev/kapp-controller/pkg/kubeconfig"
	"carvel.dev/kapp-controller/pkg/memdir"
	"carvel.dev/kapp-controller/pkg/template"
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

// AppFactory allows to create "hidden" Apps for reconciling PackageRepositories.
type AppFactory struct {
	CoreClient  kubernetes.Interface
	AppClient   kcclient.Interface
	KcConfig    *config.Config
	CmdRunner   exec.CmdRunner
	Kubeconf    *kubeconfig.Kubeconfig
	CacheFolder *memdir.TmpDir
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
	return NewCRDApp(app, pkgr, log, f.AppClient, fetchFactory, templateFactory, deployFactory)
}
