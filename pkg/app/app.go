package app

import (
	"github.com/ghodss/yaml"
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/template"
)

type AppHooks struct {
	BlockDeletion   func() error
	UnblockDeletion func() error
	UpdateStatus    func() error
}

type App struct {
	app             v1alpha1.App
	hooks           AppHooks
	fetchFactory    fetch.Factory
	templateFactory template.Factory
}

func NewApp(app v1alpha1.App, hooks AppHooks,
	fetchFactory fetch.Factory, templateFactory template.Factory) *App {
	return &App{app, hooks, fetchFactory, templateFactory}
}

func (a *App) Name() string      { return a.app.Name }
func (a *App) Namespace() string { return a.app.Namespace }

func (a *App) Status() v1alpha1.AppStatus { return a.app.Status }

func (a *App) StatusAsYAMLBytes() ([]byte, error) {
	return yaml.Marshal(a.Status())
}
