package ui

import (
	"fmt"
	"strings"

	"github.com/cppforlife/color"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/mitchellh/go-wordwrap"
)

type IPkgAuthoringUI interface {
	PrintInformationalText(text string)
	PrintCmdExecutionText(text string)
	PrintActionableText(text string)
	AskForText(textOpts ui.TextOpts) (string, error)
	AskForChoice(opts ui.ChoiceOpts) (int, error)
	PrintCmdExecutionOutput(text string)
}

type PackageAuthoringUIImpl struct {
	ui ui.UI
}

func NewPackageAuthoringUI(ui ui.UI) IPkgAuthoringUI {
	return PackageAuthoringUIImpl{
		ui: ui,
	}
}

func (c PackageAuthoringUIImpl) PrintInformationalText(text string) {
	c.ui.BeginLinef(color.New(color.Faint).Sprint("\n" + wordwrap.WrapString(text, 70)))
}

func (c PackageAuthoringUIImpl) PrintCmdExecutionText(text string) {
	c.ui.BeginLinef(fmt.Sprintf("\n\t    | $ %s\n", text))
}

func (c PackageAuthoringUIImpl) PrintCmdExecutionOutput(output string) {
	lines := strings.Split(output, "\n")
	for ind, line := range lines {
		if line != "" {
			lines[ind] = fmt.Sprintf("\t    | %s", lines[ind])
		}
	}

	indentedBlock := strings.Join(lines, "\n")
	if strings.LastIndex(indentedBlock, "\n") != (len(indentedBlock) - 1) {
		indentedBlock += "\n"
	}
	c.ui.PrintBlock([]byte(indentedBlock))
}

func (c PackageAuthoringUIImpl) PrintActionableText(text string) {
	c.ui.BeginLinef(color.New(color.Bold).Sprintf("\n%s", text))
}

func (c PackageAuthoringUIImpl) AskForText(textOpts ui.TextOpts) (string, error) {
	col := color.New(color.Bold)
	textOpts.Label = fmt.Sprintf(col.Sprint("> ")) + textOpts.Label
	return c.ui.AskForText(textOpts)
}

func (c PackageAuthoringUIImpl) AskForChoice(choiceOpts ui.ChoiceOpts) (int, error) {
	col := color.New(color.Bold)
	choiceOpts.Label = fmt.Sprintf(col.Sprint("> ")) + choiceOpts.Label
	return c.ui.AskForChoice(choiceOpts)
}
