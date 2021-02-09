package e2e

import (
	"strings"
	"testing"
)

func Test_PackageRepoBundle_PackagesAvailable(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kubectl := Kubectl{t, env.Namespace, logger}
	// contents of this bundle (ewrenn/repo-bundle:v1.0.0)
	// under examples/packaging-demo/repo-bundle
	yamlRepo := `---
apiVersion: kappctrl.k14s.io/v1alpha1
kind: PkgRepository
metadata:
  name: basic.test.carvel.dev
  # cluster scoped
spec:
  fetch:
    bundle:
      image: ewrenn/repo-bundle:v1.0.0`

	cleanUp := func() {
		kubectl.RunWithOpts([]string{"delete", "pkgrepository/basic.test.carvel.dev"}, RunOpts{NoNamespace: true})
	}
	defer cleanUp()

	kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(yamlRepo)})

	retryFunc := func() error {
		_, err := kubectl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.1.0.0"}, RunOpts{NoNamespace: true, AllowError: true})
		if err != nil {
			return err
		}
		_, err = kubectl.RunWithOpts([]string{"get", "pkg/pkg2.test.carvel.dev.2.0.0"}, RunOpts{NoNamespace: true, AllowError: true})
		if err != nil {
			return err
		}
		return nil
	}

	err := retry(5, retryFunc)
	if err != nil {
		t.Fatalf("Expected to find pkgs (pkg2.test.carvel.dev.1.0.0, pkg2.test.carvel.dev.2.0.0) but couldn't: %v", err)
	}
}

func retry(iterations int, f func() error) error {
	var err error
	for i := 0; i < iterations; i++ {
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}