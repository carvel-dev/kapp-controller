// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/go-logr/logr"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SecretHandler struct {
	client     kcclient.Interface
	log        logr.Logger
	appSecrets *reftracker.AppSecrets
}

var _ handler.EventHandler = &SecretHandler{}

func NewSecretHandler(kc kcclient.Interface, log logr.Logger, as *reftracker.AppSecrets) *SecretHandler {
	return &SecretHandler{kc, log, as}
}

func (sch *SecretHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {}

func (sch *SecretHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	sch.enqueueAppsForUpdate(evt.MetaNew.GetName(), evt.MetaNew.GetNamespace(), q)
}

func (sch *SecretHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	sch.appSecrets.RemoveSecretFromMap(evt.Meta.GetName(), evt.Meta.GetNamespace())
}

func (sch *SecretHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {}

func (sch *SecretHandler) enqueueAppsForUpdate(secretName, secretNamespace string, q workqueue.RateLimitingInterface) {
	apps := sch.appSecrets.GetAppsForSecret(secretName, secretNamespace)
	for _, app := range apps {
		sch.log.Info("enqueueing App " + app.GetAppName() + " from update to secret " + secretName)
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      app.GetAppName(),
			Namespace: secretNamespace,
		}})
	}
}
