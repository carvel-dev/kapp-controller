// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec_test

import (
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/sidecarexec"
)

func Test_TrustedCertsCreateConfig(t *testing.T) {
	backup, certs, closeFile, err := createCertTempFiles(t)
	assert.NoError(t, err)
	defer closeFile()

	backup.Write([]byte("existing-cert"))

	osConfig := sidecarexec.NewOSConfig(logr.Discard())
	osConfig.CACertsLoc = sidecarexec.OSConfigCACertsLoc{
		Path:         certs.Name(),
		OrigCopyPath: backup.Name(),
	}

	var result int
	err = osConfig.ApplyCACerts("cert-42", &result)
	assert.NoError(t, err)

	contents, err := os.ReadFile(certs.Name())
	assert.NoError(t, err)

	assert.Contains(t, string(contents), "existing-cert\ncert-42")
}

func Test_TrustedCertsUpdateConfig(t *testing.T) {
	backup, certs, closeFile, err := createCertTempFiles(t)
	assert.NoError(t, err)
	defer closeFile()

	backup.Write([]byte("existing-cert"))

	osConfig := sidecarexec.NewOSConfig(logr.Discard())
	osConfig.CACertsLoc = sidecarexec.OSConfigCACertsLoc{
		Path:         certs.Name(),
		OrigCopyPath: backup.Name(),
	}

	var result int
	err = osConfig.ApplyCACerts("cert-42", &result)
	assert.NoError(t, err)

	contents, err := os.ReadFile(certs.Name())
	assert.NoError(t, err)
	assert.Equal(t, string(contents), "existing-cert\ncert-42")

	// update config
	err = osConfig.ApplyCACerts("cert-43", &result)
	assert.NoError(t, err)

	contents, err = os.ReadFile(certs.Name())
	assert.NoError(t, err)
	assert.Equal(t, string(contents), "existing-cert\ncert-43")
}

func Test_TrustedCertsDeleteConfig(t *testing.T) {
	backup, certs, closeFile, err := createCertTempFiles(t)
	assert.NoError(t, err)
	defer closeFile()

	backup.Write([]byte("existing-cert"))

	osConfig := sidecarexec.NewOSConfig(logr.Discard())
	osConfig.CACertsLoc = sidecarexec.OSConfigCACertsLoc{
		Path:         certs.Name(),
		OrigCopyPath: backup.Name(),
	}

	// no config found
	var result int
	err = osConfig.ApplyCACerts("", &result)
	assert.NoError(t, err)

	contents, err := os.ReadFile(certs.Name())
	assert.NoError(t, err)
	// restored to the backup without any additional data
	assert.Equal(t, string(contents), "existing-cert\n")
}

func createCertTempFiles(_ *testing.T) (backup *os.File, certs *os.File, close func(), err error) {
	backup, err = os.CreateTemp("", "backup.crt")
	if err != nil {
		return nil, nil, nil, err
	}

	certs, err = os.CreateTemp("", "certs.crt")
	if err != nil {
		return nil, nil, nil, err
	}

	return backup, certs, func() {
		backup.Close()
		certs.Close()
	}, nil
}
