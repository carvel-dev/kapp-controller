package main

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	ctlapp "github.com/k14s/kapp-controller/pkg/app"
	kcclient "github.com/k14s/kapp-controller/pkg/client/clientset/versioned"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/template"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type AppFactory struct {
	coreClient kubernetes.Interface
	appClient  kcclient.Interface
}

func (f *AppFactory) NewConfigMapAppFromName(request reconcile.Request, log logr.Logger) *ctlapp.ConfigMapApp {
	return ctlapp.NewConfigMapAppFromName(request.NamespacedName, log, f.coreClient)
}

func (f *AppFactory) NewConfigMapApp(appCfgMap *corev1.ConfigMap, log logr.Logger) (*ctlapp.ConfigMapApp, error) {
	fetchFactory := fetch.NewFactory(f.coreClient)
	return ctlapp.NewConfigMapApp(appCfgMap, log, f.coreClient,
		fetchFactory, template.NewFactory(f.coreClient, fetchFactory))
}

func (f *AppFactory) NewCRDAppFromName(request reconcile.Request, log logr.Logger) *ctlapp.CRDApp {
	return ctlapp.NewCRDAppFromName(request.NamespacedName, log, f.appClient)
}

func (f *AppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) (*ctlapp.CRDApp, error) {
	fetchFactory := fetch.NewFactory(f.coreClient)
	return ctlapp.NewCRDApp(app, log, f.appClient,
		fetchFactory, template.NewFactory(f.coreClient, fetchFactory))
}
