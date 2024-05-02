// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"

	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	kcclient "carvel.dev/kapp-controller/pkg/client/clientset/versioned"
	"carvel.dev/vendir/pkg/vendir/versions"
	verv1alpha1 "carvel.dev/vendir/pkg/vendir/versions/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TODO For this PoC, we are simply going to add all packages to
// the workqueue, but in the future, we will only enqueue the
// packages that are eligible for upgrade based on the new
// packages
type PackageInstallVersionHandler struct {
	client   kcclient.Interface
	globalNS string
	log      logr.Logger
}

var _ handler.EventHandler = &PackageInstallVersionHandler{}

func NewPackageInstallVersionHandler(c kcclient.Interface, globalNS string, log logr.Logger) *PackageInstallVersionHandler {
	return &PackageInstallVersionHandler{c, globalNS, log}
}

// Create is called in response to an create event
func (ipvh *PackageInstallVersionHandler) Create(_ context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing PackageInstallList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

// Update is called in response to an update event
func (ipvh *PackageInstallVersionHandler) Update(_ context.Context, evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing PackageInstallList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.ObjectNew)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

// Delete is called in response to a delete event
func (ipvh *PackageInstallVersionHandler) Delete(_ context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing PackageInstallList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

// Generic is called in response to an event of an unknown type or a synthetic event triggered as a cron or
// external trigger request - e.g. reconcile Autoscaling, or a Webhook.
func (ipvh *PackageInstallVersionHandler) Generic(_ context.Context, evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

func (ipvh *PackageInstallVersionHandler) enqueueEligiblePackageInstalls(q workqueue.RateLimitingInterface, obj runtime.Object) error {
	pv := obj.(*datapkgingv1alpha1.Package)

	namespace := ""
	if pv.Namespace != ipvh.globalNS {
		namespace = pv.Namespace
	}

	installedPkgList, err := ipvh.client.PackagingV1alpha1().PackageInstalls(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ip := range installedPkgList.Items {
		if ip.Spec.PackageRef == nil {
			continue
		}

		if ip.Spec.PackageRef.RefName == pv.Spec.RefName && ipvh.isEligibleForVersionUpgrade(pv.Spec.Version, ip) {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      ip.Name,
				Namespace: ip.Namespace,
			}})
		}
	}

	return nil
}

func (ipvh *PackageInstallVersionHandler) isEligibleForVersionUpgrade(version string, installedPkg pkgingv1alpha1.PackageInstall) bool {
	if installedPkg.Spec.PackageRef == nil {
		return false
	}

	semverConfig := installedPkg.Spec.PackageRef.VersionSelection

	selectedVersion, err := versions.HighestConstrainedVersion([]string{version}, verv1alpha1.VersionSelection{Semver: semverConfig})
	if selectedVersion == "" || err != nil {
		return false
	}

	return true
}
