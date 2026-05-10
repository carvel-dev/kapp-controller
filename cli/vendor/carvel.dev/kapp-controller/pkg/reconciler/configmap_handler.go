// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"context"

	"carvel.dev/kapp-controller/pkg/reftracker"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
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

var _ handler.TypedEventHandler[*corev1.ConfigMap, reconcile.Request] = &ConfigMapHandler{}

func NewConfigMapHandler(log logr.Logger, as *reftracker.AppRefTracker, aus *reftracker.AppUpdateStatus) *ConfigMapHandler {
	return &ConfigMapHandler{log, as, aus}
}

// Create is called in response to create event
func (e *ConfigMapHandler) Create(ctx context.Context, evt event.TypedCreateEvent[*corev1.ConfigMap], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
}

// Update is called in response to an update event
func (e *ConfigMapHandler) Update(ctx context.Context, evt event.TypedUpdateEvent[*corev1.ConfigMap], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueueAppsForUpdate(evt.ObjectNew.GetName(), evt.ObjectNew.GetNamespace(), q)
}

// Delete is called in response to a delete event
func (e *ConfigMapHandler) Delete(ctx context.Context, evt event.TypedDeleteEvent[*corev1.ConfigMap], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
	e.appRefTracker.RemoveRef(reftracker.NewConfigMapKey(evt.Object.GetName(), evt.Object.GetNamespace()))
}

// Generic is called in response to an event of an unknown type or a synthetic event triggered as a cron or
// external trigger request - e.g. reconcile Autoscaling, or a Webhook.
func (e *ConfigMapHandler) Generic(ctx context.Context, evt event.TypedGenericEvent[*corev1.ConfigMap], q workqueue.TypedRateLimitingInterface[reconcile.Request]) {
}

func (e *ConfigMapHandler) enqueueAppsForUpdate(cfgmName, cfgmNamespace string, q workqueue.TypedRateLimitingInterface[reconcile.Request]) error {
	apps, err := e.appRefTracker.AppsForRef(reftracker.NewConfigMapKey(cfgmName, cfgmNamespace))
	if err != nil {
		return err
	}

	for refKey := range apps {
		e.log.Info("enqueueing " + refKey.Description() + " from update to configmap " + cfgmName)
		e.appUpdateStatus.MarkNeedsUpdate(refKey)
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      refKey.RefName(),
			Namespace: cfgmNamespace,
		}})
	}

	return nil
}
