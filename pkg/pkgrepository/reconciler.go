// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reconciler"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler is responsible for reconciling PackageRepositories.
type Reconciler struct {
	client          kcclient.Interface
	coreClient      kubernetes.Interface
	log             logr.Logger
	appFactory      AppFactory
	appRefTracker   *reftracker.AppRefTracker
	appUpdateStatus *reftracker.AppUpdateStatus
}

var _ reconcile.Reconciler = &Reconciler{}

// NewReconciler is the constructor for the Reconciler struct
func NewReconciler(appClient kcclient.Interface, coreClient kubernetes.Interface,
	log logr.Logger, appFactory AppFactory, appRefTracker *reftracker.AppRefTracker,
	appUpdateStatus *reftracker.AppUpdateStatus) *Reconciler {
	return &Reconciler{appClient, coreClient, log,
		appFactory, appRefTracker, appUpdateStatus}
}

// AttachWatches configures watches needed for reconciler to reconcile PackageRepository.
func (r *Reconciler) AttachWatches(controller controller.Controller, cache cache.Cache) error {
	err := controller.Watch(source.Kind(cache, &pkgv1alpha1.PackageRepository{}), &handler.EnqueueRequestForObject{})
	if err != nil {
		return fmt.Errorf("Watching PackageRepositories: %s", err)
	}

	schRepo := reconciler.NewSecretHandler(r.log, r.appRefTracker, r.appUpdateStatus)

	err = controller.Watch(source.Kind(cache, &corev1.Secret{}), schRepo)
	if err != nil {
		return fmt.Errorf("Watching Secrets: %s", err)
	}

	return nil
}

// Reconcile ensures that Packages/PackageMetadatas are imported
// into the cluster from given PackageRepository.
func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", request)

	log.Info("Reconciling")

	existingPkgRepository, err := r.client.PackagingV1alpha1().PackageRepositories(request.Namespace).Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Could not find PackageRepository", "name", request.Name)
			return reconcile.Result{}, nil // No requeue
		}

		log.Error(err, "Could not fetch PackageRepository")
		return reconcile.Result{}, err
	}

	app, err := NewPackageRepoApp(existingPkgRepository)
	if err != nil {
		return reconcile.Result{}, err
	}

	r.ReconcileFetchPlaceholderSecrets(*existingPkgRepository, app)

	crdApp := r.appFactory.NewCRDPackageRepo(app, existingPkgRepository, log)
	r.UpdatePackageRepoRefs(crdApp.ResourceRefs(), app)

	force := false
	pkgrKey := reftracker.NewPackageRepositoryKey(app.Name, app.Namespace)
	if r.appUpdateStatus.IsUpdateNeeded(pkgrKey) {
		r.appUpdateStatus.MarkUpdated(pkgrKey)
		force = true
	}

	return crdApp.Reconcile(force)
}

// ReconcileFetchPlaceholderSecrets helps determine if a placeholder secret
// needs to be created for the PackageRepository. This placeholder secret is
// populated by secretgen-controller so PackageRepositories can authenticate
// to private registries without needing to explicitly declare a secretRef.
// If no secretRef is specified for the PackageRepository, the placeholder
// is created and used by the PackageRepository.
func (r *Reconciler) ReconcileFetchPlaceholderSecrets(pkgr pkgv1alpha1.PackageRepository, app *v1alpha1.App) error {
	for i, fetch := range app.Spec.Fetch {
		if fetch.ImgpkgBundle != nil && fetch.ImgpkgBundle.SecretRef == nil {
			secretName, err := r.createSecretForSecretgenController(pkgr, i)
			if err != nil {
				return err
			}
			app.Spec.Fetch[i].ImgpkgBundle.SecretRef = &v1alpha1.AppFetchLocalRef{secretName}
		}

		if fetch.Image != nil && fetch.Image.SecretRef == nil {
			secretName, err := r.createSecretForSecretgenController(pkgr, i)
			if err != nil {
				return err
			}
			app.Spec.Fetch[i].Image.SecretRef = &v1alpha1.AppFetchLocalRef{secretName}
		}
	}
	return nil
}

func (r *Reconciler) createSecretForSecretgenController(pkgr pkgv1alpha1.PackageRepository, i int) (string, error) {
	secretName := fmt.Sprintf("%s-fetch-%d", pkgr.Name, i)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: pkgr.Namespace,
			Annotations: map[string]string{
				"secretgen.carvel.dev/image-pull-secret": "",
			},
		},
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: []byte(`{"auths":{}}`),
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}

	controllerutil.SetOwnerReference(&pkgr, secret, scheme.Scheme)

	_, err := r.coreClient.CoreV1().Secrets(pkgr.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return "", err
		}
	}
	return secretName, nil
}

// nolint: revive
func (r *Reconciler) UpdatePackageRepoRefs(refKeys map[reftracker.RefKey]struct{}, app *v1alpha1.App) {
	pkgRepoKey := reftracker.NewPackageRepositoryKey(app.Name, app.Namespace)
	// If PackageRepo is being deleted, remove
	// from all its associated references.
	if app.DeletionTimestamp != nil {
		r.appRefTracker.RemoveAppFromAllRefs(pkgRepoKey)
		return
	}

	// Add new refs for PackageRepo to AppRefTracker/remove
	// any formerly but now unused refs for PackageRepo.
	r.appRefTracker.ReconcileRefs(refKeys, pkgRepoKey)
}
