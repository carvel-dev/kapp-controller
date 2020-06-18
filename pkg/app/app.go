package app

import (
	"sync"

	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/deploy"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/template"
)

type AppHooks struct {
	BlockDeletion   func() error
	UnblockDeletion func() error
	UpdateStatus    func() error
	WatchChanges    func(func(v1alpha1.App), chan struct{}) error
}

type App struct {
	app     v1alpha1.App
	appPrev v1alpha1.App
	hooks   AppHooks

	fetchFactory    fetch.Factory
	templateFactory template.Factory
	deployFactory   deploy.Factory

	log logr.Logger

	pendingStatusUpdate   bool
	flushAllStatusUpdates bool
}

func NewApp(app v1alpha1.App, hooks AppHooks,
	fetchFactory fetch.Factory, templateFactory template.Factory,
	deployFactory deploy.Factory, log logr.Logger) *App {

	return &App{app: app, appPrev: *(app.DeepCopy()), hooks: hooks,
		fetchFactory: fetchFactory, templateFactory: templateFactory,
		deployFactory: deployFactory, log: log}
}

func (a *App) Name() string      { return a.app.Name }
func (a *App) Namespace() string { return a.app.Namespace }

func (a *App) Status() v1alpha1.AppStatus { return a.app.Status }

func (a *App) StatusAsYAMLBytes() ([]byte, error) {
	return yaml.Marshal(a.Status())
}

func (a *App) blockDeletion() error   { return a.hooks.BlockDeletion() }
func (a *App) unblockDeletion() error { return a.hooks.UnblockDeletion() }

func (a *App) updateStatus() error {
	a.pendingStatusUpdate = true

	if !a.flushAllStatusUpdates {
		// If there is no direct changes to the CRD, throttle status update
		if a.app.Generation == a.appPrev.Status.ObservedGeneration {
			return nil
		}
	}

	a.pendingStatusUpdate = false
	return a.hooks.UpdateStatus()
}

func (a *App) startFlushingAllStatusUpdates() {
	a.flushAllStatusUpdates = true
	a.flushUpdateStatus()
}

func (a *App) flushUpdateStatus() error {
	// Last possibility to save any pending status changes
	if a.pendingStatusUpdate {
		a.pendingStatusUpdate = false
		return a.hooks.UpdateStatus()
	}
	return nil
}

func (a *App) newCancelCh() (chan struct{}, func()) {
	var cancelOnce sync.Once
	cancelCh := make(chan struct{})

	// Ends watching for app changes
	cancelFunc := func() {
		cancelOnce.Do(func() { close(cancelCh) })
	}

	cancelFuncOnApp := func(app v1alpha1.App) {
		if app.Spec.Canceled {
			cancelFunc()
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
