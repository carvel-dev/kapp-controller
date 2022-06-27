// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package common

type Step interface {
	PreInteract() error
	PostInteract() error
	Interact() error
}

func Run(step Step) error {
	err := step.PreInteract()
	if err != nil {
		return err
	}
	err = step.Interact()
	if err != nil {
		return err
	}
	err = step.PostInteract()
	if err != nil {
		return err
	}
	return nil
}
