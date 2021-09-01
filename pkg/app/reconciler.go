// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler is responsible for reconciling Apps.
type Reconciler struct {
	appClient       kcclient.Interface
	log             logr.Logger
	crdAppFactory   CRDAppFactory
	appRefTracker   *reftracker.AppRefTracker
	appUpdateStatus *reftracker.AppUpdateStatus
}

// NewReconciler constructs new Reconciler.
func NewReconciler(appClient kcclient.Interface, log logr.Logger, crdAppFactory CRDAppFactory,
	appRefTracker *reftracker.AppRefTracker, appUpdateStatus *reftracker.AppUpdateStatus) *Reconciler {
	return &Reconciler{appClient, log, crdAppFactory, appRefTracker, appUpdateStatus}
}

var _ reconcile.Reconciler = &Reconciler{}

// AttachWatches configures watches needed for reconciler to reconcile Apps.
func (r *Reconciler) AttachWatches(controller controller.Controller) error {
	err := controller.Watch(&source.Kind{Type: &kcv1alpha1.App{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return fmt.Errorf("Watch Apps: %s", err)
	}

	secretHandler := reconciler.NewSecretHandler(r.log, r.appRefTracker, r.appUpdateStatus)
	cmHandler := reconciler.NewConfigMapHandler(r.log, r.appRefTracker, r.appUpdateStatus)

	err = controller.Watch(&source.Kind{Type: &corev1.Secret{}}, secretHandler)
	if err != nil {
		return fmt.Errorf("Watch Secrets: %s", err)
	}

	err = controller.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, cmHandler)
	if err != nil {
		return fmt.Errorf("Watch ConfigMaps: %s", err)
	}

	return nil
}

// nolint: revive
func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	// TODO currently we've decided to get a fresh copy of app so
	// that we do not operate on stale copy for efficiency reasons
	existingApp, err := r.appClient.KappctrlV1alpha1().Apps(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find App")
			return reconcile.Result{}, nil // No requeue
		}

		log.Error(err, "Could not fetch App")
		return reconcile.Result{}, err
	}

	crdApp := r.crdAppFactory.NewCRDApp(existingApp, log)
	r.UpdateAppRefs(crdApp.ResourceRefs(), existingApp)

	force := false
	appKey := reftracker.NewAppKey(existingApp.Name, existingApp.Namespace)
	if r.appUpdateStatus.IsUpdateNeeded(appKey) {
		r.appUpdateStatus.MarkUpdated(appKey)
		force = true
	}

	return crdApp.Reconcile(force)
}

// nolint: revive
// TODO should be made private
func (r *Reconciler) UpdateAppRefs(refKeys map[reftracker.RefKey]struct{}, app *v1alpha1.App) {
	appKey := reftracker.NewAppKey(app.Name, app.Namespace)
	// If App is being deleted, remove the App
	// from all its associated references.
	if app.DeletionTimestamp != nil {
		r.appRefTracker.RemoveAppFromAllRefs(appKey)
		return
	}

	// Add new refs for App to AppRefTracker/remove
	// any formerly but now unused refs for App.
	r.appRefTracker.ReconcileRefs(refKeys, appKey)
}

// nolint: revive
// TODO should be removed since AppRefTracker is passed into this object
func (r *Reconciler) AppRefTracker() *reftracker.AppRefTracker {
	return r.appRefTracker
}
