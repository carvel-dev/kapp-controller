package core

import (
	"fmt"
	"strings"

	"github.com/cppforlife/color"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/mitchellh/go-wordwrap"
)

type IAuthoringUI interface {
	PrintInformationalText(text string)
	PrintCmdExecutionText(text string)
	PrintActionableText(text string)
	AskForText(textOpts ui.TextOpts) (string, error)
	AskForChoice(opts ui.ChoiceOpts) (int, error)
	PrintCmdExecutionOutput(text string)
	PrintHeaderText(text string)
}

type AuthoringUIImpl struct {
	ui ui.UI
}

func NewAuthoringUI(ui ui.UI) IAuthoringUI {
	return AuthoringUIImpl{
		ui: ui,
	}
}

func (uiImpl AuthoringUIImpl) PrintInformationalText(text string) {
	uiImpl.ui.BeginLinef(color.New(color.Faint).Sprint(wordwrap.WrapString(text, 80)))
}

func (uiImpl AuthoringUIImpl) PrintCmdExecutionText(text string) {
	uiImpl.ui.BeginLinef(fmt.Sprintf("\n\t    | $ %s\n", text))
}

func (uiImpl AuthoringUIImpl) PrintCmdExecutionOutput(output string) {
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
	uiImpl.ui.PrintBlock([]byte(indentedBlock))
}

func (uiImpl AuthoringUIImpl) PrintActionableText(text string) {
	uiImpl.ui.BeginLinef(color.New(color.Bold).Sprintf("\n%s", text))
}

func (uiImpl AuthoringUIImpl) AskForText(textOpts ui.TextOpts) (string, error) {
	col := color.New(color.Bold)
	textOpts.Label = fmt.Sprintf(col.Sprint("> ")) + textOpts.Label
	return uiImpl.ui.AskForText(textOpts)
}

func (uiImpl AuthoringUIImpl) AskForChoice(choiceOpts ui.ChoiceOpts) (int, error) {
	col := color.New(color.Bold)
	choiceOpts.Label = fmt.Sprintf(col.Sprint("> ")) + choiceOpts.Label
	return uiImpl.ui.AskForChoice(choiceOpts)
}

func (uiImpl AuthoringUIImpl) PrintHeaderText(text string) {
	uiImpl.ui.BeginLinef(color.New(color.Bold).Sprintf("%s\n", text))
}
