package main

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	periodicRequeueDur = 30 * time.Second
)

type PeriodicReconciler struct {
	delegate reconcile.Reconciler
}

var _ reconcile.Reconciler = &PeriodicReconciler{}

func (r *PeriodicReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	res, err := r.delegate.Reconcile(request)
	res.RequeueAfter = periodicRequeueDur // implies Requeue=true
	return res, err
}
