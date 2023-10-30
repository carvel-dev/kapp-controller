// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"context"

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

// Create is called in response to create event
func (sch *ConfigMapHandler) Create(_ context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
}

// Update is called in response to an update event
func (sch *ConfigMapHandler) Update(_ context.Context, evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.ObjectNew.GetName(), evt.ObjectNew.GetNamespace(), q)
}

// Delete is called in response to a delete event
func (sch *ConfigMapHandler) Delete(_ context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
	sch.appRefTracker.RemoveRef(reftracker.NewConfigMapKey(evt.Object.GetName(), evt.Object.GetNamespace()))
}

// Generic is called in response to an event of an unknown type or a synthetic event triggered as a cron or
// external trigger request - e.g. reconcile Autoscaling, or a Webhook.
func (sch *ConfigMapHandler) Generic(_ context.Context, _ event.GenericEvent, _ workqueue.RateLimitingInterface) {
}

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
