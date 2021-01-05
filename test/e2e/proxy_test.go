package e2e

import (
	"fmt"
	"strings"
	"testing"
)

func TestHTTPProxy(t *testing.T) {
	logger := Logger{}
	kubectl := Kubectl{t, "kapp-controller", logger}

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
