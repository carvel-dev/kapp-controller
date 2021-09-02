// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"context"

	"github.com/go-logr/logr"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions"
	verv1alpha1 "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"
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

func (ipvh *PackageInstallVersionHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing PackageInstallList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

func (ipvh *PackageInstallVersionHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing PackageInstallList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.ObjectNew)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

func (ipvh *PackageInstallVersionHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing PackageInstallList")
	err := ipvh.enqueueEligiblePackageInstalls(q, evt.Object)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all PackageInstalls")
	}
}

func (ipvh *PackageInstallVersionHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
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
