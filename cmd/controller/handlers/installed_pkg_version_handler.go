package main

import "sigs.k8s.io/controller-runtime/pkg/handler"

type InstalledPkgVersionHandler struct{}

var _ handler.EventHandler = InstalledPkgVersionHandler{}

func (ipvh *InstalledPkgVersionHandler) Create(event.CreateEvent, workqueue.RateLimitingInterface)

func (ipvh *InstalledPkgVersionHandler) Update(event.UpdateEvent, workqueue.RateLimitingInterface)

func (ipvh *InstalledPkgVersionHandler) Delete(event.DeleteEvent, workqueue.RateLimitingInterface)

func (ipvh *InstalledPkgVersionHandler) Generic(event.GenericEvent, workqueue.RateLimitingInterface)
