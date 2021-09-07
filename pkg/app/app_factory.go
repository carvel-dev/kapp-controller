// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
)

// CRDAppFactory allows to create CRDApps.
type CRDAppFactory struct {
	CoreClient kubernetes.Interface
	AppClient  kcclient.Interface
	KcConfig   *config.Config
}

// NewCRDApp creates a CRDApp injecting necessary dependencies.
func (f *CRDAppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) *CRDApp {
	fetchFactory := fetch.NewFactory(f.CoreClient, f.KcConfig)
	templateFactory := template.NewFactory(f.CoreClient, fetchFactory)
	deployFactory := deploy.NewFactory(f.CoreClient)
	return NewCRDApp(app, log, f.AppClient, fetchFactory, templateFactory, deployFactory)
}
