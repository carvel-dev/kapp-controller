// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
)

// AppFactory allows to create "hidden" Apps for reconciling PackageRepositories.
type AppFactory struct {
	CoreClient kubernetes.Interface
	AppClient  kcclient.Interface
	KcConfig   *config.Config
}

// NewCRDPackageRepo constructs "hidden" App to reconcile PackageRepository.
func (f *AppFactory) NewCRDPackageRepo(app *kcv1alpha1.App, pkgr *pkgv1alpha1.PackageRepository, log logr.Logger) *CRDApp {
	fetchFactory := fetch.NewFactory(f.CoreClient, fetch.VendirOpts{
		SkipTLSConfig: f.KcConfig,
	})
	templateFactory := template.NewFactory(f.CoreClient, fetchFactory, false)
	deployFactory := deploy.NewFactory(f.CoreClient)
	return NewCRDApp(app, pkgr, log, f.AppClient, fetchFactory, templateFactory, deployFactory)
}
