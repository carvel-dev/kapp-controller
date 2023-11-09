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

type SecretHandler struct {
	log             logr.Logger
	appRefTracker   *reftracker.AppRefTracker
	appUpdateStatus *reftracker.AppUpdateStatus
}

var _ handler.EventHandler = &SecretHandler{}

func NewSecretHandler(log logr.Logger, as *reftracker.AppRefTracker, aus *reftracker.AppUpdateStatus) *SecretHandler {
	return &SecretHandler{log, as, aus}
}

// Create is called in response to an create event
func (sch *SecretHandler) Create(_ context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
}

// Update is called in response to an update event
func (sch *SecretHandler) Update(_ context.Context, evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.ObjectNew.GetName(), evt.ObjectNew.GetNamespace(), q)
}

// Delete is called in response to a delete event
func (sch *SecretHandler) Delete(_ context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.Object.GetName(), evt.Object.GetNamespace(), q)
	sch.appRefTracker.RemoveRef(reftracker.NewSecretKey(evt.Object.GetName(), evt.Object.GetNamespace()))
}

// Generic is called in response to an event of an unknown type or a synthetic event triggered as a cron or
// external trigger request - e.g. reconcile Autoscaling, or a Webhook.
func (sch *SecretHandler) Generic(_ context.Context, _ event.GenericEvent, _ workqueue.RateLimitingInterface) {
}

func (sch *SecretHandler) enqueueAppsForUpdate(secretName, secretNamespace string, q workqueue.RateLimitingInterface) error {
	apps, err := sch.appRefTracker.AppsForRef(reftracker.NewSecretKey(secretName, secretNamespace))
	if err != nil {
		return err
	}

	for refKey := range apps {
		sch.log.Info("enqueueing " + refKey.Description() + " from update to secret " + secretName)
		sch.appUpdateStatus.MarkNeedsUpdate(refKey)
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      refKey.RefName(),
			Namespace: secretNamespace,
		}})
	}

	return nil
}
