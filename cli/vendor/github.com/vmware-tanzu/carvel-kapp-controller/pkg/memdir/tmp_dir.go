// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package memdir

import (
	"io/ioutil"
	"os"
)

var (
	tmpDir = os.Getenv("KAPPCTRL_MEM_TMP_DIR")
)

type TmpDir struct {
	id   string
	path string
}

func NewTmpDir(id string) *TmpDir {
	if len(id) == 0 {
		panic("Expected non-empty id")
	}
	return &TmpDir{id: id}
}

func (d *TmpDir) Create() error {
	path, err := ioutil.TempDir(tmpDir, "kapp-controller-"+d.id)
	if err != nil {
		return err
	}
	d.path = path
	return nil
}

func (d *TmpDir) Remove() error {
	if len(d.path) > 0 {
		return os.RemoveAll(d.path)
	}
	return nil
}

func (d *TmpDir) Path() string { return d.path }
