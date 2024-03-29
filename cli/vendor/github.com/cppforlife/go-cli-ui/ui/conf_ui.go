package ui

import (
	. "github.com/cppforlife/go-cli-ui/ui/table"
)

type ConfUI struct {
	parent      UI
	isTTY       bool
	logger      ExternalLogger
	showColumns []Header
}

func NewConfUI(logger ExternalLogger) *ConfUI {
	var ui UI

	writerUI := NewConsoleUI(logger)
	ui = NewPaddingUI(writerUI)

	return &ConfUI{
		parent: ui,
		isTTY:  writerUI.IsTTY(),
		logger: logger,
	}
}

func NewWrappingConfUI(parent UI, logger ExternalLogger) *ConfUI {
	return &ConfUI{
		parent: parent,
		isTTY:  true,
		logger: logger,
	}
}

func (ui *ConfUI) EnableTTY(force bool) {
	if !ui.isTTY && !force {
		ui.parent = NewNonTTYUI(ui.parent)
	}
}

func (ui *ConfUI) EnableColor() {
	ui.parent = NewColorUI(ui.parent)
}

func (ui *ConfUI) EnableJSON() {
	ui.parent = NewJSONUI(ui.parent, ui.logger)
}

func (ui *ConfUI) ShowColumns(columns []Header) {
	ui.showColumns = columns
}

func (ui *ConfUI) EnableNonInteractive() {
	ui.parent = NewNonInteractiveUI(ui.parent)
}

func (ui *ConfUI) ErrorLinef(pattern string, args ...interface{}) {
	ui.parent.ErrorLinef(pattern, args...)
}

func (ui *ConfUI) PrintLinef(pattern string, args ...interface{}) {
	ui.parent.PrintLinef(pattern, args...)
}

func (ui *ConfUI) BeginLinef(pattern string, args ...interface{}) {
	ui.parent.BeginLinef(pattern, args...)
}

func (ui *ConfUI) EndLinef(pattern string, args ...interface{}) {
	ui.parent.EndLinef(pattern, args...)
}

func (ui *ConfUI) PrintBlock(block []byte) {
	ui.parent.PrintBlock(block)
}

func (ui *ConfUI) PrintErrorBlock(block string) {
	ui.parent.PrintErrorBlock(block)
}

func (ui *ConfUI) PrintTable(table Table) {
	if len(ui.showColumns) > 0 {
		err := table.SetColumnVisibility(ui.showColumns)
		if err != nil {
			panic(err)
		}
	}

	ui.parent.PrintTable(table)
}

func (ui *ConfUI) AskForText(opts TextOpts) (string, error) {
	return ui.parent.AskForText(opts)
}

func (ui *ConfUI) AskForChoice(opts ChoiceOpts) (int, error) {
	return ui.parent.AskForChoice(opts)
}

func (ui *ConfUI) AskForPassword(label string) (string, error) {
	return ui.parent.AskForPassword(label)
}

func (ui *ConfUI) AskForConfirmation() error {
	return ui.parent.AskForConfirmation()
}

func (ui *ConfUI) IsInteractive() bool {
	return ui.parent.IsInteractive()
}

func (ui *ConfUI) Flush() {
	ui.parent.Flush()
}
