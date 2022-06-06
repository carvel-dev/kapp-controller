// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ErrReconciler struct {
	delegate reconcile.Reconciler
	log      logr.Logger
}

var _ reconcile.Reconciler = &ErrReconciler{}

func (r *ErrReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	res, err := r.delegate.Reconcile(ctx, request)

	switch {
	// Check before error check, since error check does not queue immediately
	case res.Requeue && res.RequeueAfter == 0:
		log.Info("Requeue immediately", "after", 0)

	// This indicates general reconciliation error such as
	// failing to update status, or fetching necessary k8s resource.
	// Note, it does not cover cases such as deploy failing.
	case err != nil:
		res.RequeueAfter = 3*time.Second + wait.Jitter(5*time.Second, 1.0)
		log.Info("Requeue quickly due to err", "after", res.RequeueAfter)

	case res.RequeueAfter > 0:
		log.Info("Requeue after given time", "after", res.RequeueAfter)

	default:
		log.Info("Requeue stopped")
	}

	return res, err
}
