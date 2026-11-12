// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreflightFlagAllowedOnlyForDeploy(t *testing.T) {
	kapp := &Kapp{}
	opts := []string{"--preflight=PermissionValidation=true,CRDUpgradeSafety=true"}

	args, err := kapp.addRawOpts([]string{"deploy"}, opts, kappAllowedDeployFlagSet)
	require.NoError(t, err)
	require.Equal(t, append([]string{"deploy"}, opts...), args)

	_, err = kapp.addRawOpts(nil, opts, kappAllowedDeleteFlagSet)
	require.ErrorContains(t, err, "Unexpected flag '--preflight'")

	_, err = kapp.addRawOpts(nil, opts, kappAllowedInspectFlagSet)
	require.ErrorContains(t, err, "Unexpected flag '--preflight'")
}
