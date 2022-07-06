// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kappcontroller

import (
	"strings"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func TestSops(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	sas := e2e.ServiceAccounts{env.Namespace}

	yaml1 := `
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-sops
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        cm2.sops.yml: |
          apiVersion: ENC[AES256_GCM,data:TCc=,iv:YyeMABS+DjwlLBcg7RgSYD2lvUC9arab3g5XH02l4Ug=,tag:m1jlYfHpZwS5lF8QGmLfQg==,type:str]
          kind: ENC[AES256_GCM,data:0IR0Z3OPGpMJ,iv:dhPkR6HWq9Baq3qBEgh6Rtap9fYy4/LPwB+HJw3e9Oc=,tag:Z2d/P53etuarM73vqizOAg==,type:str]
          metadata:
              name: ENC[AES256_GCM,data:Nbhb,iv:IdmS1Q1dPjcYYwwKQHqMICcG2AVvsesZ2YWZMRRjBAw=,tag:j90x3EDjDBNZuI7CYjqcXg==,type:str]
          data:
              key: ENC[AES256_GCM,data:W5fepVbN5Nlt+lLhkQ==,iv:7HtxrTusZ6myNkjUBH6uudqUc9u/81NSmvkMP1rYQqs=,tag:obOHC3u9uNPCyxVlZdvecQ==,type:str]
          sops:
              kms: []
              gcp_kms: []
              azure_kv: []
              hc_vault: []
              lastmodified: '2020-10-05T15:55:37Z'
              mac: ENC[AES256_GCM,data:a/NUWKO1umLeGJaVoV7HKMxtTGEVB+C1MOXJ+RwWhsplfQn1S09//C4WHFkkLUaM6+chBhc9HDCLgp5zoZ0EyW7JnQHCq18CE/xGA/qBljeKMQVtaR0r82EwnNMCBC8Lx7jlRDky0ORLkK0JD0Pcu7ARSwTCOgf8CiuAJlgT0ZE=,iv:i4L5PEV0tIX8QI7niAOT+WyvlH6/TBBjmkFwflJ7PKE=,tag:emFBV2lk9/3nL8qCHQmqNw==,type:str]
              pgp:
              -   created_at: '2020-10-05T15:55:37Z'
                  enc: |
                      -----BEGIN PGP MESSAGE-----
                      Version: GnuPG v1

                      hQIMA/7je44gmO38ARAA0L44XpVaU4XSn+NQUJu07C4evREOsCYvdXaDOQGaZgW7
                      FoTfwDwtKIGdPfDo0pMBs5jnK5D0ukiuQEgzOGiWYkZEpAHtaT3lzbGo3x8sO5rE
                      qsz+dY0h/QMeBfgmdeDAzUlh6RIHmWf7xQuFKBj7s7q4Q3CtjSucUBa2SIlpaQ1/
                      7MYkf9CGnINXQxHu0ZJXeniV38Tsw/0UiWRGRu2grsC/nziTONdlDW5Jcq6CUahT
                      hH6Bqb41C6FE43numNtcuVr7kR+DIV1laFn0J9XqSO1+SAnaI36CSHvLOBWAPrJE
                      cOOGWDYcSpYakDOLDdDI0O9jeXVCgQjSBsMlvAHEap2946B1ORK9o7rtDY+hbglk
                      hHEDdJ0tPIf9O5Hk0u6fMnoZxDPQOx9kMZT5bsbdlBDg27PHYbfAEVAZuOywf9u5
                      Gq/CUbE7AVHkBf0q0iDIeGV3jbfeRT7jHG+LR6xrvdhDlciwsiKDbXlS/ges7Yda
                      IILY4EEBZmt1tZMzCcE5gwkaXRRGKOm5NLgTfRnneVDDDZmv5NJYfyIvRumalOAZ
                      pK7Md3EMGKWFxLMtIXlPSAN7Cd2Ce+dp0MdpOv3ioveR5wdIJVMvJc5FpkEoZvAW
                      EQpydWZtMUMeEUtuEC6QtVIyNgDsWnVvVE4bM4KlupeGn6lIVjEjMjEoHzj+L7/S
                      XAFv51Vet4dHpw0FTkYnTwpx5QrieF2zDlZ6OBUsAa9h7gE21K7HqD7+zl0rHF/A
                      cBj6pvgBGW8qFa67LLiWPT701QNA4JDkXIAd9d/f/0zwEnJV2k+9c7R5Ec9G
                      =53Ry
                      -----END PGP MESSAGE-----
                  fp: B464DFD255C6B9F8
              unencrypted_suffix: _unencrypted
              version: 3.6.1
  template:
  - sops:
      pgp:
        privateKeysSecretRef:
          name: pgp-key
  - ytt: {}
  deploy:
  - kapp: {}
---
apiVersion: v1
kind: Secret
metadata:
  name: pgp-key
stringData:
  all.pk: |
    -----BEGIN PGP PRIVATE KEY BLOCK-----
    Version: GnuPG v1

    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
    =2uSU
    -----END PGP PRIVATE KEY BLOCK-----
` + sas.ForNamespaceYAML()

	name := "test-sops"
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

GPG/SOPS usage:

$ gpg --gen-key

$ gpg --list-secret-keys --keyid-format LONG
/root/.gnupg/secring.gpg
------------------------
sec   4096R/B464DFD255C6B9F8 2020-10-03
uid                          test test (test) <test@test.com>
ssb   4096R/FEE37B8E2098EDFC 2020-10-03

$ cat foo.yml

data: |
  #@data/values
  ---
  value: encrypted-valued

$ sops --encrypt --pgp B464DFD255C6B9F8 foo.yml

$ gpg --armor --export-secret-keys B464DFD255C6B9F8

*/
