// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/deploy"
	"carvel.dev/kapp-controller/pkg/fetch"
	"carvel.dev/kapp-controller/pkg/metrics"
	"carvel.dev/kapp-controller/pkg/reftracker"
	"carvel.dev/kapp-controller/pkg/template"
	"github.com/go-logr/logr"
	types "k8s.io/apimachinery/pkg/types"
)

type Hooks struct {
	BlockDeletion   func() error
	UnblockDeletion func() error
	UpdateStatus    func(string) error
}

type App struct {
	app        v1alpha1.App
	appPrev    v1alpha1.App
	pkgRepoUID types.UID
	hooks      Hooks

	fetchFactory    fetch.Factory
	templateFactory template.Factory
	deployFactory   deploy.Factory

	appMetrics *metrics.Metrics

	log logr.Logger

	pendingStatusUpdate   bool
	flushAllStatusUpdates bool
}

// NewApp creates a new instance of an App based on v1alpha1.App
func NewApp(app v1alpha1.App, hooks Hooks, fetchFactory fetch.Factory, templateFactory template.Factory, deployFactory deploy.Factory,
	log logr.Logger, appMetrics *metrics.Metrics, pkgRepoUID types.UID) *App {
	return &App{app: app, appPrev: *(app.DeepCopy()), hooks: hooks,
		fetchFactory: fetchFactory, templateFactory: templateFactory,
		deployFactory: deployFactory, appMetrics: appMetrics, log: log, pkgRepoUID: pkgRepoUID}
}

func (a *App) Name() string      { return a.app.Name }
func (a *App) Namespace() string { return a.app.Namespace }

// Kind return kind of pkg repo
func (a *App) Kind() string { return "PackageRepository" }

func (a *App) Status() v1alpha1.AppStatus { return a.app.Status }

func (a *App) blockDeletion() error   { return a.hooks.BlockDeletion() }
func (a *App) unblockDeletion() error { return a.hooks.UnblockDeletion() }

func (a *App) updateStatus(desc string) error {
	a.pendingStatusUpdate = true

	if !a.flushAllStatusUpdates {
		// If there is no direct changes to the CRD, throttle status update
		if a.app.Generation == a.appPrev.Status.ObservedGeneration {
			return nil
		}
	}

	a.pendingStatusUpdate = false
	return a.hooks.UpdateStatus(desc)
}

func (a *App) startFlushingAllStatusUpdates() {
	a.flushAllStatusUpdates = true
	a.flushUpdateStatus("flush all")
}

func (a *App) flushUpdateStatus(desc string) error {
	// Last possibility to save any pending status changes
	if a.pendingStatusUpdate {
		a.pendingStatusUpdate = false
		return a.hooks.UpdateStatus("flushing: " + desc)
	}
	return nil
}

func (a *App) SecretRefs() map[reftracker.RefKey]struct{} {
	secrets := map[reftracker.RefKey]struct{}{}

	// Fetch SecretRefs
	for _, fetch := range a.app.Spec.Fetch {
		switch {
		case fetch.Image != nil && fetch.Image.SecretRef != nil:
			refKey := reftracker.NewSecretKey(fetch.Image.SecretRef.Name, a.app.Namespace)
			secrets[refKey] = struct{}{}
		case fetch.ImgpkgBundle != nil && fetch.ImgpkgBundle.SecretRef != nil:
			refKey := reftracker.NewSecretKey(fetch.ImgpkgBundle.SecretRef.Name, a.app.Namespace)
			secrets[refKey] = struct{}{}
		case fetch.HTTP != nil && fetch.HTTP.SecretRef != nil:
			refKey := reftracker.NewSecretKey(fetch.HTTP.SecretRef.Name, a.app.Namespace)
			secrets[refKey] = struct{}{}
		case fetch.Git != nil && fetch.Git.SecretRef != nil:
			refKey := reftracker.NewSecretKey(fetch.Git.SecretRef.Name, a.app.Namespace)
			secrets[refKey] = struct{}{}
		}
	}

	return secrets
}

// HasImageOrImgpkgBundle is used to determine if the
// App's spec contains a fetch stage for an image or
// imgpkgbundle. It is mainly used to determine whether
// to retry a fetch attempt when placeholder secrets are
// involved with authenticating to private registries. Placeholder
// secrets are not always populated quick enough for Apps to use
// the secret, and private auth is only supported for images/bundles,
// so this helps to narrow down when to retry a fetch attempt.
func (a App) HasImageOrImgpkgBundle() bool {
	for _, fetch := range a.app.Spec.Fetch {
		if fetch.ImgpkgBundle != nil || fetch.Image != nil {
			return true
		}
	}
	return false
}
