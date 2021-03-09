// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	if !r.areAppSecretsUpToDate(crdApp.GetSecretRefs(), existingApp) {
		force = true
	}
	//if !force && r.areAppConfigMapsUpToDate(crdApp.GetConfigMapRefs(), request.Namespace, existingApp.Name)

	result, err := crdApp.Reconcile(force)
	if err != nil {
		return result, err
	}

	r.appRefTracker.MarkAppUpdated(existingApp.Name, existingApp.Namespace)
	return result, err

}

// Check whether secrets used by App are latest
// versions. Helps determine whether to force an
// update to App during Reconcile.

// WHY this func is needed: We don't know where the
// reconcileRequest originates from so we need to determine
// if we should force reconciliation. Also, we need a way to
// have secrets be associated with apps and this is how Apps
// will register with their secretRefs.
func (r *AppsReconciler) areAppSecretsUpToDate(secretNames []string, app *v1alpha1.App) bool {
	// No updates if paused or cancelled
	if app.Spec.Canceled || app.Spec.Paused {
		return true
	}

	// If App is being deleted, remove App from appRefTracker.
	if app.DeletionTimestamp != nil {
		for _, secretName := range secretNames {
			r.appRefTracker.RemoveAppFromRefMap(v1.Secret{}.Kind, secretName, app.Namespace, app.Name)
			r.appRefTracker.RemoveAppFromUpdateMap(app.Name, app.Namespace)
		}
		return true
	}

	// Make sure secrets for App are always up to date
	// in appRefTracker.
	for _, secretName := range secretNames {
		if !r.appRefTracker.CheckAppExistsForRef(v1.Secret{}.Kind, secretName, app.Namespace, app.Name) {
			r.appRefTracker.AddAppToRefMap(v1.Secret{}.Kind, secretName, app.Namespace, app.Name)
			if r.appRefTracker.CheckAppExistsForRef(v1.Secret{}.Kind, secretName, app.Namespace, app.Name) {
			}
		}
	}

	if r.appRefTracker.GetAppUpdateStatus(app.Name, app.Namespace) {
		return false
	}

	return true
}
