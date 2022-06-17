// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
)

// OSConfig provides RPC interface system configuration.
type OSConfig struct {
	log logr.Logger

	// Mostly used for tests
	CACertsLoc   OSConfigCACertsLoc
	SetenvFunc   func(key, value string) error
	UnsetenvFunc func(string) error
}

// OSConfigCACertsLoc is a set of CA cert paths needed for cert management.
type OSConfigCACertsLoc struct {
	Path         string
	OrigCopyPath string
}

// NewOSConfig returns new OSConfig.
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

// ApplyCACerts atomically updates existing CA certs file
// with additional CA certs provided.
func (r OSConfig) ApplyCACerts(chain string, unusedResult *int) error {
	r.log.Info("Applying CA certs")

	origCopyFile, err := os.Open(r.CACertsLoc.OrigCopyPath)
	if err != nil {
		return fmt.Errorf("Opening original certs file: %s", err)
	}
	defer origCopyFile.Close()

	// Place tmp file into dst directory as Rename call below
	// cannot succeed if tmp file is on a different mount
	// (no guarantee that /tmp and /etc are on same fs).
	tmpFile, err := os.CreateTemp(filepath.Dir(r.CACertsLoc.Path), "tmp-ca-bundle-")
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

// ProxyInput describes proxy configuration.
type ProxyInput struct {
	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string
}

// ApplyProxy sets proxy related environment variables.
func (r OSConfig) ApplyProxy(in ProxyInput, unusedResult *int) error {
	vals := map[string]string{
		"http_proxy":  in.HTTPProxy,
		"https_proxy": in.HTTPSProxy,
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
