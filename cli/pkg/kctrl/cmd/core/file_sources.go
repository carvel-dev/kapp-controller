// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"io/ioutil"
	"os"
)

type InputFile struct {
	source FileSource
}

func NewInputFile(file string) InputFile {
	switch file {
	case "-":
		return InputFile{source: NewStdinSource()}
	default:
		return InputFile{source: NewLocalFileSource(file)}
	}
}

func (f InputFile) Bytes() ([]byte, error) {
	return f.source.Bytes()
}

type FileSource interface {
	Bytes() ([]byte, error)
}

type StdinSource struct{}

var _ FileSource = StdinSource{}

func NewStdinSource() StdinSource            { return StdinSource{} }
func (s StdinSource) Bytes() ([]byte, error) { return ioutil.ReadAll(os.Stdin) }

type LocalFileSource struct {
	path string
}

var _ FileSource = LocalFileSource{}

func NewLocalFileSource(path string) LocalFileSource { return LocalFileSource{path} }
func (s LocalFileSource) Bytes() ([]byte, error)     { return ioutil.ReadFile(s.path) }
