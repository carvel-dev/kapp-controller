package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	Ytt int = iota
	HelmTemplate
)

type TemplateStep struct {
	Ui          ui.UI
	AppTemplate []v1alpha1.AppTemplate
}

func NewTemplateStep(ui ui.UI) *TemplateStep {
	templateStep := TemplateStep{
		Ui: ui,
	}
	return &templateStep
}

func (t TemplateStep) PreInteract() error {
	str := `
# Next step is to define which templating tool to be used to render the package template. 
# A package template can be rendered by various different tools.`
	t.Ui.PrintBlock([]byte(str))
	return nil
}

func (t TemplateStep) PostInteract() error {
	return nil
}

func (t *TemplateStep) Interact() error {
	templateType, err := t.Ui.AskForChoice("Enter the templating tool to be used", []string{"ytt(recommended)", "helmTemplate"})
	if err != nil {
		return err
	}
	var appTemplateList []v1alpha1.AppTemplate
	switch templateType {
	case Ytt:
		yttTemplateStep := NewYttTemplateStep(t.Ui)
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
