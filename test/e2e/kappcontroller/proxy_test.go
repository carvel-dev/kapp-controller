// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func TestHTTPProxy(t *testing.T) {
	assert := assert.New(t)
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller", logger}

	// Proxy configured in config-test/secret-config.yml
	logger.Section("inspect controller logs for propogation of proxy env vars", func() {
		// app name must match the app name being deployed in hack/deploy-test.sh
		out := kubectl.Run([]string{"logs", "deployment/kapp-controller"})

		assert.Containsf(out, "http_proxy is enabled.",
			"expected log line detailing http_proxy is enabled")

		assert.Containsf(out, "no_proxy is enabled.",
			"expected log line detailing no_proxy is enabled")
	})

}
