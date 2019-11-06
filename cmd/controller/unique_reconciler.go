package main

import (
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type UniqueReconciler struct {
	delegate reconcile.Reconciler

	ongoing    map[string]struct{}
	ongoingMux sync.Mutex

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

	r.allowRequeues(resKey)

	if r.shouldReconcileAndMarkOngoing(resKey) {
		res, err := r.delegate.Reconcile(request)
		r.markReconcileFinished(resKey)
		return res, err
	}

	return r.requeueIfAllowed(resKey)
}

func (r *UniqueReconciler) allowRequeues(resKey string) {
	r.pendingMux.Lock()
	delete(r.pending, resKey)
	r.pendingMux.Unlock()
}

func (r *UniqueReconciler) requeueIfAllowed(resKey string) (reconcile.Result, error) {
	// Avoid excessive requeueing of same resource
	r.pendingMux.Lock()
	defer r.pendingMux.Unlock()

	if _, found := r.pending[resKey]; !found {
		r.pending[resKey] = struct{}{}
		return reconcile.Result{Requeue: true}, nil
	}

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

func (r *UniqueReconciler) markReconcileFinished(resKey string) {
	r.ongoingMux.Lock()
	delete(r.ongoing, resKey)
	r.ongoingMux.Unlock()
}
