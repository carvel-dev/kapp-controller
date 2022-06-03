package template

import (
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	YttFilesLocation int = iota
	Inline
)

type YttTemplateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	appTemplateYtt v1alpha1.AppTemplateYtt
}

func NewYttTemplateStep(ui pkgui.IPkgAuthoringUI) *YttTemplateStep {
	return &YttTemplateStep{
		pkgAuthoringUI: ui,
	}
}

func (y YttTemplateStep) PreInteract() error {
	y.pkgAuthoringUI.PrintInformationalText(`# We need to provide the values to ytt. They can be done in two different ways:
# 1. We can specify the files(including data values) to be used via ytt. Multiple paths can be provided with comma separated values.
# 2. We can enter the values directly i.e. inline`)
	return nil
}

func (y YttTemplateStep) PostInteract() error {
	return nil
}

func (y *YttTemplateStep) Interact() error {
	input, err := y.pkgAuthoringUI.AskForChoice(ui.ChoiceOpts{
		Label:   "Enter how do you prefer to provide values to ytt",
		Default: 0,
		Choices: []string{"ytt files location(recommended)", "inline"},
	})
	if err != nil {
		return err
	}
	switch input {
	case YttFilesLocation:
		paths, err := y.pkgAuthoringUI.AskForText(ui.TextOpts{
			Label:        "Enter the paths of ytt files",
			Default:      "",
			ValidateFunc: nil,
		})
		if err != nil {
			return err
		}
		y.appTemplateYtt = v1alpha1.AppTemplateYtt{Paths: strings.Split(paths, ",")}
	case Inline:

	}
	return nil
}

func (y *YttTemplateStep) Run() error {
	y.PreInteract()
	y.Interact()
	y.PostInteract()
	return nil
}
