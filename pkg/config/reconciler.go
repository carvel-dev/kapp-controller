// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler is responsible for reconciling Apps.
type Reconciler struct {
	coreClient kubernetes.Interface
	log        logr.Logger
}

// NewReconciler constructs new Reconciler.
func NewReconciler(coreClient kubernetes.Interface, log logr.Logger) *Reconciler {
	return &Reconciler{coreClient, log}
}

var _ reconcile.Reconciler = &Reconciler{}

// AttachWatches configures watches needed for reconciler to reconcile Apps.
func (r *Reconciler) AttachWatches(controller controller.Controller) error {
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetName() == "kapp-controller-config"
		},
	}

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

// nolint: revive
func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	kcConfig, err := GetConfig(r.coreClient)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("getting kapp-controller config: %s", err)
	}

	log.Info("Applying new config")
	err = kcConfig.Apply()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("Applying configuration: %s", err)
	}

	return reconcile.Result{}, nil
}
