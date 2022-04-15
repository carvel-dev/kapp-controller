// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler is responsible for reconciling Kapp-controllers config.
type Reconciler struct {
	coreClient kubernetes.Interface
	log        logr.Logger
}

// NewReconciler constructs new Reconciler.
func NewReconciler(coreClient kubernetes.Interface, log logr.Logger) *Reconciler {
	return &Reconciler{coreClient, log}
}

var _ reconcile.Reconciler = &Reconciler{}

// AttachWatches configures watches needed for reconciler to reconcile the kapp-controller Config.
func (r *Reconciler) AttachWatches(controller controller.Controller, ns string) error {
	// only reconcile on the KC's config
	p := predicate.NewPredicateFuncs(func(o client.Object) bool {
		return o.GetNamespace() == ns && o.GetName() == kcConfigName
	})

	err := controller.Watch(&source.Kind{Type: &v1.ConfigMap{}}, &handler.EnqueueRequestForObject{}, p)
	if err != nil {
		return fmt.Errorf("Watching Configmaps: %s", err)
	}

	err = controller.Watch(&source.Kind{Type: &v1.Secret{}}, &handler.EnqueueRequestForObject{}, p)
	if err != nil {
		return fmt.Errorf("Watching Secrets: %s", err)
	}

	return nil
}

// Reconcile gets the current config from the cluster and applies any changes.
func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	kcConfig, err := GetConfig(r.coreClient)
	if err != nil {
		log.Error(err, "getting kapp-controller config")
		return reconcile.Result{}, nil // no re-queue
	}

	log.Info("Applying new config")
	err = kcConfig.Apply()
	if err != nil {
		log.Error(err, "applying kapp-controller config")
		return reconcile.Result{}, nil // no re-queue
	}

	return reconcile.Result{}, nil
}
