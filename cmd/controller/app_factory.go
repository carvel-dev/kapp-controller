// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	ctlapp "github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	ctlpkgr "github.com/vmware-tanzu/carvel-kapp-controller/pkg/pkgrepository"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
)

type AppFactory struct {
	coreClient kubernetes.Interface
	appClient  kcclient.Interface
	kcConfig   *config.Config
	appMetrics *metrics.AppMetrics
}

// NewCRDApp creates new CRD app
func (f *AppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) *ctlapp.CRDApp {
	fetchFactory := fetch.NewFactory(f.coreClient, f.kcConfig)
	templateFactory := template.NewFactory(f.coreClient, fetchFactory)
	deployFactory := deploy.NewFactory(f.coreClient)

	return ctlapp.NewCRDApp(app, log, f.appMetrics, f.appClient, fetchFactory, templateFactory, deployFactory)
}

// TODO: Create a PackageRepo factory for this func
func (f *AppFactory) NewCRDPackageRepo(app *kcv1alpha1.App, pkgr *pkgingv1alpha1.PackageRepository, log logr.Logger) *ctlpkgr.CRDApp {
	fetchFactory := fetch.NewFactory(f.coreClient, f.kcConfig)
	templateFactory := template.NewFactory(f.coreClient, fetchFactory)
	deployFactory := deploy.NewFactory(f.coreClient)
	return ctlpkgr.NewCRDApp(app, pkgr, log, f.appClient, fetchFactory, templateFactory, deployFactory)
}
