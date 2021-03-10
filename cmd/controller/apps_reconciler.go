// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	secret    = "secret"
	configmap = "configmap"
)

type AppsReconciler struct {
	appClient     kcclient.Interface
	log           logr.Logger
	appFactory    AppFactory
	appRefTracker *reftracker.AppRefTracker
}

var _ reconcile.Reconciler = &AppsReconciler{}

func (r *AppsReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	// TODO currently we've decided to get a fresh copy of app so
	// that we do not operate on stale copy for efficiency reasons
	existingApp, err := r.appClient.KappctrlV1alpha1().Apps(request.Namespace).Get(request.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find App")
			return reconcile.Result{}, nil // No requeue
		}

		log.Error(err, "Could not fetch App")
		return reconcile.Result{}, err
	}

	force := false
	crdApp := r.appFactory.NewCRDApp(existingApp, log)
	if !r.areAppRefsUpToDate(crdApp.GetSecretRefs(), secret, existingApp) ||
		!r.areAppRefsUpToDate(crdApp.GetConfigMapRefs(), configmap, existingApp) {
		force = true
	}

	result, err := crdApp.Reconcile(force)
	if err != nil {
		return result, err
	}

	// Only update if reconcile is successful.
	// Leave as not updated if error occurs
	// in case follow up reconcile request may
	// address issue for App.
	r.appRefTracker.MarkAppUpdated(existingApp.Name, existingApp.Namespace)

	return result, err
}

func (r *AppsReconciler) areAppRefsUpToDate(refNames map[string]struct{}, kind string, app *v1alpha1.App) bool {
	// No updates if paused or cancelled
	if app.Spec.Canceled || app.Spec.Paused {
		return true
	}

	// If App is being deleted, remove App from appRefTracker.
	if app.DeletionTimestamp != nil {
		for refName := range refNames {
			r.appRefTracker.RemoveAppFromRefMap(kind, refName, app.Namespace, app.Name)
			r.appRefTracker.RemoveAppFromUpdateMap(app.Name, app.Namespace)
		}
		return true
	}

	// Make sure refs for App are always up to date
	// in appRefTracker.
	for refName := range refNames {
		if !r.appRefTracker.CheckAppExistsForRef(kind, refName, app.Namespace, app.Name) {
			r.appRefTracker.AddAppToRefMap(kind, refName, app.Namespace, app.Name)
		}
	}

	if r.appRefTracker.GetAppUpdateStatus(app.Name, app.Namespace) {
		return false
	}

	return true
}
