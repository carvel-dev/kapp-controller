 // Copyright 2020 VMware, Inc.
 // SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	ctlapp "github.com/k14s/kapp-controller/pkg/app"
	kcclient "github.com/k14s/kapp-controller/pkg/client/clientset/versioned"
	"github.com/k14s/kapp-controller/pkg/deploy"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
)

type AppFactory struct {
	coreClient                kubernetes.Interface
	appClient                 kcclient.Interface
	allowSharedServiceAccount bool
}

func (f *AppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) *ctlapp.CRDApp {
	fetchFactory := fetch.NewFactory(f.coreClient)
	templateFactory := template.NewFactory(f.coreClient, fetchFactory)
	deployFactory := deploy.NewFactory(f.coreClient, allowSharedServiceAccount)
	return ctlapp.NewCRDApp(app, log, f.appClient, fetchFactory, templateFactory, deployFactory)
}
