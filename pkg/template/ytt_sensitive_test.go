// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"errors"
	"fmt"
	goexec "os/exec"
	"strings"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"carvel.dev/kapp-controller/pkg/fetch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type secretDiagnosticRunner struct{ fail bool }

func (r secretDiagnosticRunner) Run(cmd *goexec.Cmd) error {
	_, _ = fmt.Fprint(cmd.Stdout, "rendered-canary")
	_, _ = fmt.Fprint(cmd.Stderr, "schema rejected password-canary")
	if r.fail {
		return errors.New("password-canary")
	}
	return nil
}
func (r secretDiagnosticRunner) RunWithCancel(cmd *goexec.Cmd, _ chan struct{}) error {
	return r.Run(cmd)
}
func TestSecretBackedYttDoesNotExposeDiagnosticValues(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "values", Namespace: "default"}, Data: map[string][]byte{"values.yml": []byte("password: password-canary")}})
	for _, fail := range []bool{false, true} {
		opts := v1alpha1.AppTemplateYtt{ValuesFrom: []v1alpha1.AppTemplateValuesSource{{SecretRef: &v1alpha1.AppTemplateValuesSourceRef{Name: "values"}}}}
		subject := NewYtt(opts, AppContext{Namespace: "default"}, client, fetch.Factory{}, secretDiagnosticRunner{fail: fail}, AdditionalDownwardAPIValues{})
		result, _ := subject.TemplateDir(t.TempDir())
		if strings.Contains(result.Stderr+result.ErrorStr(), "canary") {
			t.Fatal("Secret value in diagnostics")
		}
		if result.Stdout != "rendered-canary" {
			t.Fatal("rendered output changed")
		}
	}
}
