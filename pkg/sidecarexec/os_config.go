// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-logr/logr"
)

type OSConfig struct {
	log logr.Logger

	// Mostly used for tests
	CACertsLoc   OSConfigCACertsLoc
	SetenvFunc   func(key, value string) error
	UnsetenvFunc func(string) error
}

type OSConfigCACertsLoc struct {
	Path         string
	OrigCopyPath string
}

func NewOSConfig(log logr.Logger) OSConfig {
	return OSConfig{
		log: log,
		CACertsLoc: OSConfigCACertsLoc{
			Path:         "/etc/pki/tls/certs/ca-bundle.crt",
			OrigCopyPath: "/etc/pki/tls/certs/ca-bundle.crt.orig",
		},
		SetenvFunc:   os.Setenv,
		UnsetenvFunc: os.Unsetenv,
	}
}

func (r OSConfig) ApplyCACerts(chain string, unusedResult *int) error {
	r.log.Info("Applying CA certs")

	origCopyFile, err := os.Open(r.CACertsLoc.OrigCopyPath)
	if err != nil {
		return fmt.Errorf("Opening original certs file: %s", err)
	}
	defer origCopyFile.Close()

	tmpFile, err := os.CreateTemp(os.TempDir(), "tmp-ca-bundle-")
	if err != nil {
		return fmt.Errorf("Creating tmp certs file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, origCopyFile)
	if err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("Copying certs file: %s", err)
	}

	_, err = tmpFile.Write([]byte("\n" + chain))
	if err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err = tmpFile.Close(); err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), r.CACertsLoc.Path)
	if err != nil {
		return fmt.Errorf("Renaming certs file: %s", err)
	}

	return nil
}

type ProxyInput struct {
	HTTPProxy  string
	HTTPsProxy string
	NoProxy    string
}

func (r OSConfig) ApplyProxy(in ProxyInput, unusedResult *int) error {
	vals := map[string]string{
		"http_proxy":  in.HTTPProxy,
		"https_proxy": in.HTTPsProxy,
		"no_proxy":    in.NoProxy,
	}

	for envVar, val := range vals {
		if val == "" {
			r.UnsetenvFunc(envVar)
			r.UnsetenvFunc(strings.ToUpper(envVar))
			r.log.Info("Clearing " + envVar)
		} else {
			r.SetenvFunc(envVar, val)
			r.SetenvFunc(strings.ToUpper(envVar), val)
			r.log.Info("Setting " + envVar)
		}
	}
	return nil
}
