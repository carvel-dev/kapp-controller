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

type AppsReconciler struct {
	appClient       kcclient.Interface
	log             logr.Logger
	appFactory      AppFactory
	appRefTracker   *reftracker.AppRefTracker
	appUpdateStatus *reftracker.AppUpdateStatus
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

	crdApp := r.appFactory.NewCRDApp(existingApp, log)
	//r.UpdateAppRefs(crdApp.ConfigMapRefs(), existingApp)
	//r.UpdateAppRefs(crdApp.SecretRefs(), existingApp)
	r.UpdateAppRefs(crdApp.ResourceRefs(), existingApp)

	force := false
	if r.appUpdateStatus.IsUpdateNeeded(existingApp.Name, existingApp.Namespace) {
		r.appUpdateStatus.MarkUpdated(existingApp.Name, existingApp.Namespace)
		force = true
	}

	return crdApp.Reconcile(force)
}

func (r *AppsReconciler) UpdateAppRefs(refKeys map[reftracker.RefKey]struct{}, app *v1alpha1.App) {
	// If App is being deleted, remove the App
	// from all its associated references.
	if app.DeletionTimestamp != nil {
		r.appRefTracker.RemoveAppFromAllRefs(app.Name, app.Namespace)
		return
	}

	// Make sure refs for App are always up to date
	// in appRefTracker.
	for refKey := range refKeys {
		r.appRefTracker.AddAppForRef(refKey, app.Name)
	}

	// Remove any reference associations for App
	// if the App no longer uses refs it once did.
	r.appRefTracker.PruneAppFromRefs(refKeys, app.Name, app.Namespace)
}

func (r *AppsReconciler) AppRefTracker() *reftracker.AppRefTracker {
	return r.appRefTracker
}

func (r *AppsReconciler) SetAppRefTracker(appRefTracker *reftracker.AppRefTracker) {
	r.appRefTracker = appRefTracker
}
