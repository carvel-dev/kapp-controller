/*
 * Copyright 2020 VMware, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type UniqueReconciler struct {
	delegate reconcile.Reconciler

	// Represents currently executing reconcilation requests
	ongoing    map[string]struct{}
	ongoingMux sync.Mutex

	// Represents whether reconilation should happen immediately after
	// due to reconcilcation requests coming in while one was being
	// already executed
	pending    map[string]struct{}
	pendingMux sync.Mutex
}

var _ reconcile.Reconciler = &UniqueReconciler{}

func NewUniqueReconciler(delegate reconcile.Reconciler) *UniqueReconciler {
	return &UniqueReconciler{
		delegate: delegate,
		ongoing:  map[string]struct{}{},
		pending:  map[string]struct{}{},
	}
}

func (r *UniqueReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	resKey := request.NamespacedName.String()

	if r.shouldReconcileAndMarkOngoing(resKey) {
		res, err := r.delegate.Reconcile(request)

		// If there are any pending "reconciliation requests", reconcile immediately
		if r.isPendingAndUnmark(resKey) {
			res.Requeue = true
			res.RequeueAfter = 0
		}

		r.unmarkOngoing(resKey)
		return res, err
	}

	// Mark to be requeued at the end of ongoing reconciliation request
	r.markPending(resKey)
	return reconcile.Result{}, nil
}

func (r *UniqueReconciler) shouldReconcileAndMarkOngoing(resKey string) bool {
	r.ongoingMux.Lock()
	defer r.ongoingMux.Unlock()

	if _, found := r.ongoing[resKey]; !found {
		r.ongoing[resKey] = struct{}{}
		return true
	}

	return false
}

func (r *UniqueReconciler) unmarkOngoing(resKey string) {
	r.ongoingMux.Lock()
	delete(r.ongoing, resKey)
	r.ongoingMux.Unlock()
}

func (r *UniqueReconciler) isPendingAndUnmark(resKey string) bool {
	r.pendingMux.Lock()
	_, found := r.pending[resKey]
	delete(r.pending, resKey)
	r.pendingMux.Unlock()
	return found
}

func (r *UniqueReconciler) markPending(resKey string) {
	r.pendingMux.Lock()
	r.pending[resKey] = struct{}{}
	r.pendingMux.Unlock()
}
