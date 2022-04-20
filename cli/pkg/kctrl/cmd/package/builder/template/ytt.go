package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"strings"
)

const (
	YttFilesLocation int = iota
	Inline
	ValuesFromConfigMapOrSecret
)

type YttTemplateStep struct {
	Ui             ui.UI
	appTemplateYtt v1alpha1.AppTemplateYtt
}

func NewYttTemplateStep(ui ui.UI) *YttTemplateStep {
	return &YttTemplateStep{
		Ui: ui,
	}
}

func (y YttTemplateStep) PreInteract() error {
	str := `# We need to provide the values to ytt. They can be done in three different ways:
# 1. We can specify the files(including data values) to be used via ytt. Multiple paths can be provided with comma separated values.
# 2. We can enter the values directly i.e. inline`
	y.Ui.PrintBlock([]byte(str))
	return nil
}

func (y YttTemplateStep) PostInteract() error {
	return nil
}

func (y *YttTemplateStep) Interact() error {
	input, err := y.Ui.AskForChoice("Enter how do you prefer to provide values to ytt", []string{"ytt files location(recommended)", "inline", "valuesFrom configMap or secret"})
	if err != nil {
		return err
	}
	switch input {
	case YttFilesLocation:
		paths, err := y.Ui.AskForText("Enter the paths of ytt files")
		if err != nil {
			return err
		}
		y.appTemplateYtt = v1alpha1.AppTemplateYtt{Paths: strings.Split(paths, ",")}
	case Inline:

	case ValuesFromConfigMapOrSecret:
	}
	return nil
}

func (y *YttTemplateStep) Run() error {
	y.PreInteract()
	y.Interact()
	y.PostInteract()
	return nil
}
