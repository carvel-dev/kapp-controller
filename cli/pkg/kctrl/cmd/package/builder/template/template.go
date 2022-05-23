package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	Ytt int = iota
	HelmTemplate
)

type TemplateStep struct {
	PkgAuthoringUI pkgui.IPkgAuthoringUI
	AppTemplate    []v1alpha1.AppTemplate
}

func NewTemplateStep(ui pkgui.IPkgAuthoringUI) *TemplateStep {
	templateStep := TemplateStep{
		PkgAuthoringUI: ui,
	}
	return &templateStep
}

func (t TemplateStep) PreInteract() error {
	t.PkgAuthoringUI.PrintInformationalText(`
# Next step is to define which templating tool to be used to render the package template. 
# A package template can be rendered by various different tools.`)
	return nil
}

func (t TemplateStep) PostInteract() error {
	return nil
}

func (t *TemplateStep) Interact() error {
	templateType, err := t.PkgAuthoringUI.AskForChoice(ui.ChoiceOpts{
		Label:   "Enter the templating tool to be used",
		Default: 0,
		Choices: []string{"ytt(recommended)", "helmTemplate"},
	})
	if err != nil {
		return err
	}
	var appTemplateList []v1alpha1.AppTemplate
	switch templateType {
	case Ytt:
		yttTemplateStep := NewYttTemplateStep(t.PkgAuthoringUI)
		yttTemplateStep.Run()
		yttAppTemplate := v1alpha1.AppTemplate{
			Ytt: &yttTemplateStep.appTemplateYtt,
		}
		appTemplateList = append(appTemplateList, yttAppTemplate)
		t.AppTemplate = appTemplateList
	case HelmTemplate:
	}
	return nil
}
