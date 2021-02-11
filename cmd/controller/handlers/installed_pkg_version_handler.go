// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/go-logr/logr"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	err := ipvh.enqueueAllPackages(q)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueAllPackages(q)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueAllPackages(q)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	ipvh.log.Info("enqueueing installedPkgList")
	err := ipvh.enqueueAllPackages(q)
	if err != nil {
		ipvh.log.Error(err, "enqueueing all installed pakcages")
	}
}

func (ipvh *InstalledPkgVersionHandler) enqueueAllPackages(q workqueue.RateLimitingInterface) error {
	installedPkgList, err := ipvh.client.InstallV1alpha1().InstalledPackages("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ip := range installedPkgList.Items {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      ip.Name,
			Namespace: ip.Namespace,
		}})
	}

	return nil
}
