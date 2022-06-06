package template

import (
	pkgbuilder "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/build"

	"github.com/cppforlife/go-cli-ui/ui"
	pkgui "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/builder/ui"
)

const (
	YttFilesLocation int = iota
	Inline
)

type YttTemplateStep struct {
	pkgAuthoringUI pkgui.IPkgAuthoringUI
	pkgBuild       *pkgbuilder.PackageBuild
	pkgLocation    string
}

func NewYttTemplateStep(ui pkgui.IPkgAuthoringUI, pkgLocation string, pkgBuild *pkgbuilder.PackageBuild) *YttTemplateStep {
	return &YttTemplateStep{
		pkgAuthoringUI: ui,
		pkgLocation:    pkgLocation,
		pkgBuild:       pkgBuild,
	}
}

func (yttTemplateStep YttTemplateStep) PreInteract() error {
	yttTemplateStep.pkgAuthoringUI.PrintInformationalText(`# We need to provide the values to ytt. They can be done in two different ways:
# 1. We can specify the files(including data values) to be used via ytt. Multiple paths can be provided with comma separated values.
# 2. We can enter the values directly i.e. inline`)
	return nil
}

func (yttTemplateStep YttTemplateStep) PostInteract() error {
	return nil
}

func (yttTemplateStep *YttTemplateStep) Interact() error {
	input, err := yttTemplateStep.pkgAuthoringUI.AskForChoice(ui.ChoiceOpts{
		Label:   "Enter how do you prefer to provide values to ytt",
		Default: 0,
		Choices: []string{"ytt files location(recommended)", "inline"},
	})
	if err != nil {
		return err
	}
	switch input {
	case YttFilesLocation:
		/*paths, err := yttTemplateStep.pkgAuthoringUI.AskForText(ui.TextOpts{
			Label:        "Enter the paths of ytt files",
			Default:      "",
			ValidateFunc: nil,
		})
		if err != nil {
			return err
		}
		yttTemplateStep.appTemplateYtt = v1alpha1.AppTemplateYtt{Paths: strings.Split(paths, ",")}*/
	case Inline:

	}
	return nil
}
