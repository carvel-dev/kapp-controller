// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/k14s/semver/v4"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/deploy"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/metrics"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// ComponentInfo provides information about components of the system required by templating stage
type ComponentInfo interface {
	KappControllerVersion() (semver.Version, error)
	KubernetesVersion(serviceAccountName string, specCluster *v1alpha1.AppCluster, objMeta *metav1.ObjectMeta) (semver.Version, error)
	KubernetesAPIs() ([]string, error)
}

type Hooks struct {
	BlockDeletion   func() error
	UnblockDeletion func() error
	UpdateStatus    func(string) error
	WatchChanges    func(func(v1alpha1.App), chan struct{}) error
}

// Opts keeps App reconciliation options
type Opts struct {
	DefaultSyncPeriod time.Duration
	MinimumSyncPeriod time.Duration
}

type App struct {
	app     v1alpha1.App
	appPrev v1alpha1.App
	hooks   Hooks

	fetchFactory    fetch.Factory
	templateFactory template.Factory
	deployFactory   deploy.Factory
	compInfo        ComponentInfo

	memoizedKubernetesVersion string
	memoizedKubernetesAPIs    []string

	log        logr.Logger
	opts       Opts
	appMetrics *metrics.AppMetrics

	pendingStatusUpdate   bool
	flushAllStatusUpdates bool
	metadata              *deploy.Meta
}

func NewApp(app v1alpha1.App, hooks Hooks,
	fetchFactory fetch.Factory, templateFactory template.Factory,
	deployFactory deploy.Factory, log logr.Logger, opts Opts, appMetrics *metrics.AppMetrics, compInfo ComponentInfo) *App {

	return &App{app: app, appPrev: *(app.DeepCopy()), hooks: hooks,
		fetchFactory: fetchFactory, templateFactory: templateFactory,
		deployFactory: deployFactory, log: log, opts: opts, appMetrics: appMetrics, compInfo: compInfo}
}

func (a *App) Name() string      { return a.app.Name }
func (a *App) Namespace() string { return a.app.Namespace }

func (a *App) Status() v1alpha1.AppStatus { return a.app.Status }

func (a *App) StatusAsYAMLBytes() ([]byte, error) {
	return yaml.Marshal(a.Status())
}

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

func (a *App) newCancelCh(conditions ...cancelCondition) (chan struct{}, func()) {
	var cancelOnce sync.Once
	cancelCh := make(chan struct{})

	// Ends watching for app changes
	cancelFunc := func() {
		cancelOnce.Do(func() { close(cancelCh) })
	}

	cancelFuncOnApp := func(app v1alpha1.App) {
		for _, condition := range conditions {
			if condition(app) {
				cancelFunc()
				return
			}
		}
	}

	go func() {
		if a.hooks.WatchChanges == nil {
			// do nothing when cannot watch for changes
			return
		}

		err := a.hooks.WatchChanges(cancelFuncOnApp, cancelCh)
		if err != nil {
			a.log.Error(err, "Watching changes") // TODO remove
		}
	}()

	return cancelCh, cancelFunc
}

// Get all SecretRefs from App spec
func (a *App) SecretRefs() map[reftracker.RefKey]struct{} {
	secrets := map[reftracker.RefKey]struct{}{}

	// Fetch SecretRefs
	for _, fetch := range a.app.Spec.Fetch {
		switch {
		case fetch.Inline != nil && fetch.Inline.PathsFrom != nil:
			for _, pathsFrom := range fetch.Inline.PathsFrom {
				if pathsFrom.SecretRef != nil {
					refKey := reftracker.NewSecretKey(pathsFrom.SecretRef.Name, a.app.Namespace)
					secrets[refKey] = struct{}{}
				}
			}
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
		case fetch.HelmChart != nil && fetch.HelmChart.Repository != nil:
			if fetch.HelmChart.Repository.SecretRef != nil {
				refKey := reftracker.NewSecretKey(fetch.HelmChart.Repository.SecretRef.Name, a.app.Namespace)
				secrets[refKey] = struct{}{}
			}
		default:
		}
	}

	// Templating SecretRefs
	for _, tpl := range a.app.Spec.Template {
		switch {
		case tpl.Ytt != nil:
			if tpl.Ytt.Inline != nil {
				for _, pathsFrom := range tpl.Ytt.Inline.PathsFrom {
					if pathsFrom.SecretRef != nil {
						refKey := reftracker.NewSecretKey(pathsFrom.SecretRef.Name, a.app.Namespace)
						secrets[refKey] = struct{}{}
					}
				}
			}
			for _, valuesFrom := range tpl.Ytt.ValuesFrom {
				if valuesFrom.SecretRef != nil {
					refKey := reftracker.NewSecretKey(valuesFrom.SecretRef.Name, a.app.Namespace)
					secrets[refKey] = struct{}{}
				}
			}
		case tpl.HelmTemplate != nil && tpl.HelmTemplate.ValuesFrom != nil:
			for _, valsFrom := range tpl.HelmTemplate.ValuesFrom {
				if valsFrom.SecretRef != nil {
					refKey := reftracker.NewSecretKey(valsFrom.SecretRef.Name, a.app.Namespace)
					secrets[refKey] = struct{}{}
				}
			}
		case tpl.Cue != nil && tpl.Cue.ValuesFrom != nil:
			for _, valsFrom := range tpl.Cue.ValuesFrom {
				if valsFrom.SecretRef != nil {
					refKey := reftracker.NewSecretKey(valsFrom.SecretRef.Name, a.app.Namespace)
					secrets[refKey] = struct{}{}
				}
			}
		default:
		}
	}

	return secrets
}

// Get all ConfigMapRefs from App spec
func (a *App) ConfigMapRefs() map[reftracker.RefKey]struct{} {
	configMaps := map[reftracker.RefKey]struct{}{}

	// Fetch ConfigMapRefs
	for _, fetch := range a.app.Spec.Fetch {
		switch {
		case fetch.Inline != nil && fetch.Inline.PathsFrom != nil:
			for _, pathsFrom := range fetch.Inline.PathsFrom {
				if pathsFrom.ConfigMapRef != nil {
					refKey := reftracker.NewConfigMapKey(pathsFrom.ConfigMapRef.Name, a.app.Namespace)
					configMaps[refKey] = struct{}{}
				}
			}
		default:
		}
	}

	// Templating ConfigMapRefs
	for _, tpl := range a.app.Spec.Template {
		switch {
		case tpl.Ytt != nil:
			if tpl.Ytt.Inline != nil {
				for _, pathsFrom := range tpl.Ytt.Inline.PathsFrom {
					if pathsFrom.ConfigMapRef != nil {
						refKey := reftracker.NewConfigMapKey(pathsFrom.ConfigMapRef.Name, a.app.Namespace)
						configMaps[refKey] = struct{}{}
					}
				}
			}
			for _, valuesFrom := range tpl.Ytt.ValuesFrom {
				if valuesFrom.ConfigMapRef != nil {
					refKey := reftracker.NewConfigMapKey(valuesFrom.ConfigMapRef.Name, a.app.Namespace)
					configMaps[refKey] = struct{}{}
				}
			}
		case tpl.HelmTemplate != nil && tpl.HelmTemplate.ValuesFrom != nil:
			for _, valsFrom := range tpl.HelmTemplate.ValuesFrom {
				if valsFrom.ConfigMapRef != nil {
					refKey := reftracker.NewConfigMapKey(valsFrom.ConfigMapRef.Name, a.app.Namespace)
					configMaps[refKey] = struct{}{}
				}
			}
		case tpl.Cue != nil && tpl.Cue.ValuesFrom != nil:
			for _, valsFrom := range tpl.Cue.ValuesFrom {
				if valsFrom.ConfigMapRef != nil {
					refKey := reftracker.NewConfigMapKey(valsFrom.ConfigMapRef.Name, a.app.Namespace)
					configMaps[refKey] = struct{}{}
				}
			}
		default:
		}
	}

	return configMaps
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
