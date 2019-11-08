package main

import (
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PeriodicReconciler struct {
	delegate reconcile.Reconciler
}

var _ reconcile.Reconciler = &PeriodicReconciler{}

func (r *PeriodicReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	res, err := r.delegate.Reconcile(request)
	res.RequeueAfter = 25*time.Second + wait.Jitter(5*time.Second, 1.0) // implies Requeue=true
	return res, err
}
