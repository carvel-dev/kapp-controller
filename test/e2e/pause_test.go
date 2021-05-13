// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func Test_InstalledPackage_SetsPauseOnApp(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}
	name := "instl-pkg-pause"

	cleanUp := func() {
		// Need to make sure InstalledPackage App is not paused
		// so need to create as part of deletion.
		installPkgYaml := installedPackageYAML(name, env.Namespace, "original", false)
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	// create originally as not paused otherwise
	// App will not create original resources.
	installPkgYaml := installedPackageYAML(name, env.Namespace, "original", false)
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

	retry(t, 10*time.Second, func() error {
		_, err := kubectl.RunWithOpts([]string{"get", "apps/" + name}, RunOpts{AllowError: true})
		if err != nil {
			return fmt.Errorf("failed to get App for InstalledPackage %s", name)
		}
		return nil
	})

	// update App to be paused
	installPkgYaml = installedPackageYAML(name, env.Namespace, "original", true)
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

	// try to change configmap value
	installPkgYaml = installedPackageYAML(name, env.Namespace, "change", true)
	kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"}, RunOpts{StdinReader: strings.NewReader(installPkgYaml)})

	var cr v1alpha1.App
	retry(t, 10*time.Second, func() error {
		out, err := kubectl.RunWithOpts([]string{"get", "apps/" + name, "-o", "yaml"}, RunOpts{AllowError: true})
		if err != nil {
			return fmt.Errorf("failed to get App for InstalledPackage %s", name)
		}

		err = yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal: %s", err)
		}

		if cr.Status.FriendlyDescription != "Canceled/paused" {
			return fmt.Errorf("expected App for InstalledPackage to have status show Canceled/paused\nGot: %s", cr.Status.FriendlyDescription)
		}

		return nil
	})

	retry(t, 10*time.Second, func() error {
		out := kubectl.Run([]string{"get", "configmap/configmap", "-o", "yaml"})
		var cm corev1.ConfigMap
		err := yaml.Unmarshal([]byte(out), &cm)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %s", err)
		}

		if cm.Data["key"] != "original" {
			return fmt.Errorf("configmap message was updated despite App being paused\nGot: %s", cm.Data["key"])
		}
		return nil
	})
}

func installedPackageYAML(name, namespace, configMapValue string, paused bool) string {
	sas := ServiceAccounts{namespace}
	return fmt.Sprintf(`---
apiVersion: package.carvel.dev/v1alpha1
kind: Package
metadata:
 name: pause.pkg.1.0.0
spec:
 publicName: pause.pkg
 version: 1.0.0
 template:
   spec:
     fetch:
     - inline:
         paths:
           file.yml: |
             apiVersion: v1
             kind: ConfigMap
             metadata:
               name: configmap
             data:
               key: %s
     template:
     - ytt: {}
     deploy:
     - kapp: {}
---
apiVersion: install.package.carvel.dev/v1alpha1
kind: InstalledPackage
metadata:
 name: %s
 namespace: %s
 annotations:
   kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/installedpackages
spec:
 serviceAccountName: kappctrl-e2e-ns-sa
 paused: %t
 packageRef:
   publicName: pause.pkg
   version: 1.0.0
`, configMapValue, name, namespace, paused) + sas.ForNamespaceYAML()
}
