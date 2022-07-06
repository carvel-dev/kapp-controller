// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func TestSopsAge(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	yaml1 := `
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-sops-age
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        cm2.sops.yml: |
          apiVersion: ENC[AES256_GCM,data:vz4=,iv:M40G8gC5zZVYTZON3ECUYEIqesc3ixmJv7mLDDMKUpk=,tag:gQ+8lLPn6tE+zZQpLsug6A==,type:str]
          kind: ENC[AES256_GCM,data:fApRFHY+r/rw,iv:ONcPw4wybz5pM7FxEGCrtrCAfnLZ2O2W044fzSPPXEg=,tag:c2DX9UViBsQdbUBoeRJDhw==,type:str]
          metadata:
            name: ENC[AES256_GCM,data:yvYQ,iv:c5C/q9lb7vIWaGtFBzNNJsgiMgcd2MPjdBVwB1IrfDs=,tag:Qe6yMEPWg7dWXoikfL9sZA==,type:str]
          data:
            key: ENC[AES256_GCM,data:R/29Esw1+0jfUpjYWQ==,iv:bSvLOfZ+jxxr6kkU9RY++JWSF4wwnW6igpC8OsXFN3c=,tag:kp9O76jCLJoVzKUBSx6OsA==,type:str]
          sops:
            kms: []
            gcp_kms: []
            azure_kv: []
            hc_vault: []
            age:
            - recipient: age1s6cmqg446e5atg7akl4ryg29dlrt8j5xgw6lj2n24jgk0tz9x39ql4tu0f
              enc: |
                  -----BEGIN AGE ENCRYPTED FILE-----
                  YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBnZk5SMzh2dDVCMVcxdXo0
                  WGFHbEVUZEtseWNxY1NSelNWQllIV09yZTJVCkhmbWlicmtETGgyVDR2NXp0OURW
                  RndYR3Arb0V3bnhmeXVOUDkzVVYwRVUKLS0tIHhRaXNtT2ZESldsQVNaZHhpaWpO
                  aTNGQTlVY2FHY0FBUjYxQi9GMzZ0WlkKPZyI49+KTI+3RhqQlDqQo8WN+M6oCbrZ
                  ak/Ih323P6gBv+b//DFp0/SiQIejOMLGZsWg4vpPnNH1fJgaT0JpEA==
                  -----END AGE ENCRYPTED FILE-----
            lastmodified: "2021-09-24T17:22:16Z"
            mac: ENC[AES256_GCM,data:/GDo+52OIJKt8flWt+YUbsHDA6qZqiQ9Hc2Svj3GhRNMka1eWSyjSXbYn3JRaySiqobkn5AfpWdU0hXl/3tT/dBAtLF8Tt/BKWPObjzeQWpRutZsRyXP1Riyxm50ZqqWxi+Ut9FtmhF06QEfURaVwbKDmoNy9Ut5rNMJ5rKTxHY=,iv:EzNH+2WJUpuuUIr1hfYSlLjTsQ/+eTNce83E1DIu4fo=,tag:LNyZt5BbUhiXgZ11qlhu8w==,type:str]
            pgp: []
            unencrypted_suffix: _unencrypted
            version: 3.7.1
  template:
  - sops:
      age:
        privateKeysSecretRef:
          name: age-key
  - ytt: {}
  deploy:
  - kapp: {}
---
apiVersion: v1
kind: Secret
metadata:
  name: age-key
stringData:
  key.txt: |
    # created: 2021-09-23T17:12:23-07:00
    # public key: age1s6cmqg446e5atg7akl4ryg29dlrt8j5xgw6lj2n24jgk0tz9x39ql4tu0f
    AGE-SECRET-KEY-16VAQN88H4DAJQXXEAH4HL8UK44409XEYUQMD9KUZ47K8H7TERYVQGNYHL6

` + sas.ForNamespaceYAML()

	name := "test-sops-age"
	cleanUp := func() {
		kapp.Run([]string{"delete", "-a", name})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			e2e.RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})
	})

	logger.Section("verify fully encrypted configmap", func() {
		out := kapp.Run([]string{"inspect", "-a", name + ".app", "--raw", "--tty=false", "--filter-kind-name", "ConfigMap/cm2"})

		var cm corev1.ConfigMap

		err := yaml.Unmarshal([]byte(out), &cm)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}
		if cm.Data["key"] != "cm2-encrypted" {
			t.Fatalf(`Expected data.key to be "cm2-encrypted" got %#v`, cm.Data["key"])
		}
	})
}

/*

cm2.yml

apiVersion: v1
kind: ConfigMap
metadata:
  name: cm2
data:
  key: cm2-encrypted

*/

/*

AGE/SOPS usage:

# this is the file in age-key StringData above
$ age-keygen -o key.txt

# this is the command to encrypt a yaml file with the key file: cmd line takes public key and the keyfile (public+private keys)
$ SOPS_AGE_KEY_FILE=./key.txt  sops --encrypt --age age1s6cmqg446e5atg7akl4ryg29dlrt8j5xgw6lj2n24jgk0tz9x39ql4tu0f examples/simple-app-http.yml

*/
