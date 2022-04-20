package common

import (
	goCLIUI "github.com/cppforlife/go-cli-ui/ui"
)

var ui goCLIUI.UI

type Step interface {
	PreInteract() error
	PostInteract() error
	Interact() error
	Run() error
}

func ValidateInputYesOrNo(input string) (bool, bool) {
	if input == "y" || input == "Y" {
		return true, true
	} else if input == "n" || input == "N" {
		return false, true
	}
	return false, false
}