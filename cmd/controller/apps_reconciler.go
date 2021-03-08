// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type AppsReconciler struct {
	kubeclient kubernetes.Interface
	appClient  kcclient.Interface
	log        logr.Logger
	appFactory AppFactory
	appSecrets *reftracker.AppSecrets
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
	if !r.areAppSecretsUpToDate(crdApp.GetSecretRefs(), request.Namespace, existingApp.Name) {
		force = true
	}

	return crdApp.Reconcile(force)
}

// Check whether secrets used by App are latest
// versions. Helps determine whether to force an
// update to App during Reconcile.

// WHY this func is needed: We don't know where the
// reconcileRequest originates from so we need to determine
// if we should force reconciliation. Also, we need a way to
// have secrets be associated with apps and this is how Apps
// will register with their secretRefs.
func (r *AppsReconciler) areAppSecretsUpToDate(secretNames []string, namespace, appName string) bool {
	for _, secretName := range secretNames {
		secret, err := r.kubeclient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
		if err != nil {
			r.log.Info("could not find Secret for App " + appName)
			continue
		}
		appsEntry, err := r.appSecrets.GetSpecificAppForSecret(secret.Name, namespace, appName)
		if err != nil {
			r.log.Info("could not find App for Secret " + secret.Name + ". Adding to map.")
			r.appSecrets.AddAppToMap(secret.Name, namespace, appName, secret.ResourceVersion)
		}
		// Get app entry again since it may have just been added
		appsEntry, _ = r.appSecrets.GetSpecificAppForSecret(secret.Name, namespace, appName)
		if secret.ResourceVersion != appsEntry.GetResourceVersion() {
			// Get app entry again since it may have just been added
			appsEntry, _ = r.appSecrets.GetSpecificAppForSecret(secret.Name, namespace, appName)
			// Reflect changes to secret version in map
			r.appSecrets.UpdateAppInMap(secret.Name, namespace, appName, secret.ResourceVersion)
			return false
		}
	}
	return true
}
