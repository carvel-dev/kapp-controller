// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTkgpackageclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tkgpackageclient Suite")
}
