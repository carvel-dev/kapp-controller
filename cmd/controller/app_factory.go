// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	ctlapp "github.com/vmware-tanzu/carvel-kapp-controller/pkg/app"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
)

type AppFactory struct {
	coreClient kubernetes.Interface
	appClient  kcclient.Interface
}

func (f *AppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) *ctlapp.CRDApp {
	fetchFactory := fetch.NewFactory(f.coreClient)
	templateFactory := template.NewFactory(f.coreClient, fetchFactory)
	deployFactory := deploy.NewFactory(f.coreClient)
	return ctlapp.NewCRDApp(app, log, f.appClient, fetchFactory, templateFactory, deployFactory)
}
