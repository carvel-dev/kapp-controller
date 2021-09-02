// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"github.com/go-logr/logr"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ConfigMapHandler struct {
	log             logr.Logger
	appRefTracker   *reftracker.AppRefTracker
	appUpdateStatus *reftracker.AppUpdateStatus
}

var _ handler.EventHandler = &ConfigMapHandler{}

func NewConfigMapHandler(log logr.Logger, as *reftracker.AppRefTracker, aus *reftracker.AppUpdateStatus) *ConfigMapHandler {
	return &ConfigMapHandler{log, as, aus}
}

func (sch *ConfigMapHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
}

func (sch *ConfigMapHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.ObjectNew.GetName(), evt.ObjectNew.GetNamespace(), q)
}

func (sch *ConfigMapHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
	sch.appRefTracker.RemoveRef(reftracker.NewConfigMapKey(evt.Object.GetName(), evt.Object.GetNamespace()))
}

func (sch *ConfigMapHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {}

func (sch *ConfigMapHandler) enqueueAppsForUpdate(cfgmName, cfgmNamespace string, q workqueue.RateLimitingInterface) error {
	apps, err := sch.appRefTracker.AppsForRef(reftracker.NewConfigMapKey(cfgmName, cfgmNamespace))
	if err != nil {
		return err
	}

	for refKey := range apps {
		sch.log.Info("enqueueing " + refKey.Description() + " from update to configmap " + cfgmName)
		sch.appUpdateStatus.MarkNeedsUpdate(refKey)
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      refKey.RefName(),
			Namespace: cfgmNamespace,
		}})
	}

	return nil
}
