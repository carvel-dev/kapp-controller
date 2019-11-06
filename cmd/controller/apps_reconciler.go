package main

import (
	"context"

	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type AppsReconciler struct {
	client     client.Client
	log        logr.Logger
	appFactory AppFactory
}

var _ reconcile.Reconciler = &AppsReconciler{}

func (r *AppsReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	app := &kcv1alpha1.App{}

	err := r.client.Get(context.TODO(), request.NamespacedName, app)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find App") // TODO
			return r.appFactory.NewCRDAppFromName(request, log).Delete()
		}

		log.Error(err, "Could not fetch App")
		return reconcile.Result{}, err
	}

	crdApp, err := r.appFactory.NewCRDApp(app, log)
	if crdApp == nil || err != nil {
		return reconcile.Result{}, err
	}

	return crdApp.Reconcile()
}
