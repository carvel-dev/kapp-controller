// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"github.com/go-logr/logr"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/resourcetracker"
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
	appSecrets *resourcetracker.AppSecrets
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
	if !r.areAppSecretsUpToDate(crdApp.GetApp().GetSecretRefs(), request.Namespace, existingApp.Name) {
		crdApp.GetApp().SetReconcileMarker()
	}

	return crdApp.Reconcile()
}

// Check whether secrets used by App are latest
// versions. Helps determine whether to force an
// update to App during Reconcile.
func (r *AppsReconciler) areAppSecretsUpToDate(secretNames []string, namespace, appName string) bool {
	// TODO: Remove logging used for debugging
	for _, secretName := range secretNames {
		secret, err := r.kubeclient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
		if err != nil {
			r.log.Info("could not find Secret for App " + appName)
			// TODO: Should secret be removed in this case?
			continue
		}
		appsEntry, err := r.appSecrets.GetSpecificAppForSecret(secret.Name, namespace, appName)
		if err != nil {
			r.log.Info("could not find App for Secret " + secret.Name + ". Adding to map.")
			r.appSecrets.AddAppToMap(secret.Name, namespace, appName, secret.ResourceVersion)
		}
		r.log.Info("CHECKING SECRET DIFF")
		// Get app entry again since it may have just been added
		appsEntry, _ = r.appSecrets.GetSpecificAppForSecret(secret.Name, namespace, appName)
		if secret.ResourceVersion != appsEntry.GetResourceVersion() {
			// Get app entry again since it may have just been added
			appsEntry, _ = r.appSecrets.GetSpecificAppForSecret(secret.Name, namespace, appName)
			r.log.Info("DIFF: " + secret.ResourceVersion + " " + appsEntry.GetResourceVersion())
			// Reflect changes to secret version in map
			r.appSecrets.UpdateAppInMap(secret.Name, namespace, appName, secret.ResourceVersion)
			return false
		}
		r.log.Info("NO DIFF")
	}
	return true
}
