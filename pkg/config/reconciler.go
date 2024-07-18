// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler is responsible for reconciling kapp-controllers config.
type Reconciler struct {
	coreClient kubernetes.Interface
	config     *Config
	osConfig   OSConfig
	log        logr.Logger
}

// NewReconciler constructs new Reconciler.
func NewReconciler(coreClient kubernetes.Interface,
	config *Config, osConfig OSConfig, log logr.Logger) *Reconciler {

	return &Reconciler{coreClient, config, osConfig, log}
}

var _ reconcile.Reconciler = &Reconciler{}

// AttachWatches configures watches needed for reconciler to reconcile the kapp-controller Config.
func (r *Reconciler) AttachWatches(controller controller.Controller, ns string, mgr manager.Manager) error {
	// only reconcile on the KC's config
	configMapPredicate := []predicate.TypedPredicate[*v1.ConfigMap]{
		predicate.NewTypedPredicateFuncs[*v1.ConfigMap](func(cm *v1.ConfigMap) bool {
			return cm.GetNamespace() == ns && cm.GetName() == kcConfigName
		}),
	}

	err := controller.Watch(
		source.Kind(mgr.GetCache(), &v1.ConfigMap{}, &handler.TypedEnqueueRequestForObject[*v1.ConfigMap]{}, configMapPredicate...),
	)
	if err != nil {
		return fmt.Errorf("Watching Configmaps: %s", err)
	}

	secretPredicate := []predicate.TypedPredicate[*v1.Secret]{
		predicate.NewTypedPredicateFuncs[*v1.Secret](func(s *v1.Secret) bool {
			return s.GetNamespace() == ns && s.GetName() == kcConfigName
		}),
	}

	err = controller.Watch(source.Kind(mgr.GetCache(), &v1.Secret{}, &handler.TypedEnqueueRequestForObject[*v1.Secret]{}, secretPredicate...))
	if err != nil {
		return fmt.Errorf("Watching Secrets: %s", err)
	}

	return nil
}

// Reconcile gets the current config from the cluster and applies any changes.
func (r *Reconciler) Reconcile(_ context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	err := r.config.Reload()
	if err != nil {
		log.Error(err, "Getting kapp-controller config")
		return reconcile.Result{}, nil // no re-queue
	}

	log.Info("Applying new config")

	err = r.osConfig.ApplyCACerts(r.config.CACerts())
	if err != nil {
		log.Error(err, "Failed applying CA certificates")
		// continue on
	}

	err = r.osConfig.ApplyProxy(r.config.ProxyOpts())
	if err != nil {
		log.Error(err, "Failed applying proxy opts")
		// continue on
	}

	return reconcile.Result{}, nil // no re-queue
}
