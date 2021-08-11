// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func TestHTTPProxy(t *testing.T) {
	logger := e2e.Logger{}
	kubectl := e2e.Kubectl{t, "kapp-controller", logger}

	// These two variables must match their respective values in config-test/confi-map.yml
	proxyURL := "proxy-svc.proxy-server.svc.cluster.local:80"
	noProxyDomains := "github.com,docker.io"

	logger.Section("inspect controller logs for propogation of proxy env vars", func() {
		// app name must match the app name being deployed in hack/deploy-test.sh
		out := kubectl.Run([]string{"logs", "deployment/kapp-controller"})

		if !strings.Contains(out, fmt.Sprintf("Using http proxy '%s'", proxyURL)) {
			t.Fatalf("expected log line detailing http_proxy settings")
		}

		if !strings.Contains(out, fmt.Sprintf("No proxy set for: %s", noProxyDomains)) {
			t.Fatalf("expected log line detailing no_proxy settings")
		}
	})

}
