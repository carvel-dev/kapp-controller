package main

import (
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PeriodicReconciler struct {
	delegate reconcile.Reconciler
	log      logr.Logger
}

var _ reconcile.Reconciler = &PeriodicReconciler{}

func (r *PeriodicReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	res, err := r.delegate.Reconcile(request)

	switch {
	case res.Requeue && res.RequeueAfter == 0:
		log.Info("Requeue immediately", "after", 0)

	// This indicates general reconcilation error such as
	// failing to update status, or fetching necessary k8s resource.
	// Note, it does not cover cases such as deploy failing.
	case err != nil:
		res.RequeueAfter = 3*time.Second + wait.Jitter(5*time.Second, 1.0)
		log.Info("Requeue quickly due to err", "after", res.RequeueAfter)

	case res.RequeueAfter > 0 && res.RequeueAfter < 30*time.Second:
		log.Info("Requeue after given time", "after", res.RequeueAfter)

	default:
		res.RequeueAfter = 25*time.Second + wait.Jitter(5*time.Second, 1.0)
		log.Info("Requeue soon", "after", res.RequeueAfter)
	}

	return res, err
}
