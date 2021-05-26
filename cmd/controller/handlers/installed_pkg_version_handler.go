// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"context"

	"github.com/go-logr/logr"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
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
type InstalledPkgVersionHandler struct {
	client kcclient.Interface
	log    logr.Logger
}

var _ handler.EventHandler = &InstalledPkgVersionHandler{}

func NewInstalledPkgVersionHandler(c kcclient.Interface, log logr.Logger) *InstalledPkgVersionHandler {
	return &InstalledPkgVersionHandler{c, log}
}

func (ipvh *InstalledPkgVersionHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueEligibleInstalledPackages(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueEligibleInstalledPackages(q, evt.ObjectNew)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueEligibleInstalledPackages(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueEligibleInstalledPackages(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) enqueueEligibleInstalledPackages(q workqueue.RateLimitingInterface, obj runtime.Object) error {
	pv := obj.(*datapkgingv1alpha1.PackageVersion)
	installedPkgList, err := ipvh.client.PackagingV1alpha1().InstalledPackages("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ip := range installedPkgList.Items {
		if ip.Spec.PackageVersionRef == nil {
			continue
		}

		if ip.Spec.PackageVersionRef.PackageName == pv.Spec.PackageName && ipvh.isEligibleForVersionUpgrade(pv.Spec.Version, ip) {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      ip.Name,
				Namespace: ip.Namespace,
			}})
		}
	}

	return nil
}

func (ipvh *InstalledPkgVersionHandler) isEligibleForVersionUpgrade(version string, installedPkg pkgingv1alpha1.InstalledPackage) bool {
	if installedPkg.Spec.PackageVersionRef == nil {
		return false
	}

	semverConfig := installedPkg.Spec.PackageVersionRef.VersionSelection

	selectedVersion, err := versions.HighestConstrainedVersion([]string{version}, versions.VersionSelection{Semver: semverConfig})
	if selectedVersion == "" || err != nil {
		return false
	}

	return true
}
