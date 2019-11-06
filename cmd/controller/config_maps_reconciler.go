package main

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ConfigMapsReconciler struct {
	client     client.Client
	log        logr.Logger
	appFactory AppFactory
}

var _ reconcile.Reconciler = &ConfigMapsReconciler{}

func (r *ConfigMapsReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	appCfgMap := &corev1.ConfigMap{}

	err := r.client.Get(context.TODO(), request.NamespacedName, appCfgMap)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find app ConfigMap") // TODO
			return r.appFactory.NewConfigMapAppFromName(request, log).Delete()
		}

		log.Error(err, "Could not fetch app ConfigMap")
		return reconcile.Result{}, err
	}

	app, err := r.appFactory.NewConfigMapApp(appCfgMap, log)
	if app == nil || err != nil {
		return reconcile.Result{}, err
	}

	return app.Reconcile()
}
