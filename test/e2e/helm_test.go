package e2e

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

func TestHelm(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kapp := Kapp{t, env.Namespace, logger}

	yaml1 := `
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-helm
spec:
  fetch:
  - helmChart:
      name: stable/redis
  template:
  - helmTemplate:
      valuesFrom:
      - secretRef:
          name: test-helm-values
  deploy:
  - kapp:
      delete:
        rawOptions: ["--apply-ignored"]
---
apiVersion: v1
kind: Secret
metadata:
  name: test-helm-values
stringData:
  data.yml: |
    password: "1234567891234"
`

	name := "test-helm"
	cleanUp := func() {
		kapp.RunWithOpts([]string{"delete", "-a", name}, RunOpts{AllowError: true})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy", func() {
		kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
			RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})

		out := kapp.Run([]string{"inspect", "-a", name, "--raw", "--tty=false", "--filter-kind", "App"})

		var cr v1alpha1.App

		err := yaml.Unmarshal([]byte(out), &cr)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s", err)
		}

		expectedStatus := v1alpha1.AppStatus{
			Conditions: []v1alpha1.AppCondition{{
				Type: v1alpha1.ReconcileSucceeded,
				Status: corev1.ConditionTrue,
			}},
			Deploy: &v1alpha1.AppStatusDeploy{
				ExitCode: 0,
				Finished: true,
			},
			Fetch: &v1alpha1.AppStatusFetch{
				ExitCode: 0,
			},
			Inspect: &v1alpha1.AppStatusInspect{
				ExitCode: 0,
			},
			Template: &v1alpha1.AppStatusTemplate{
				ExitCode: 0,
			},
			ObservedGeneration: 1,
		}

		{
			// deploy
			if !strings.Contains(cr.Status.Deploy.Stdout, "Wait to:") {
				t.Fatalf("Expected non-empty deploy output: '%s'", cr.Status.Deploy.Stdout)	
			}
			cr.Status.Deploy.StartedAt = metav1.Time{}
			cr.Status.Deploy.UpdatedAt = metav1.Time{}
			cr.Status.Deploy.Stdout = ""

			// fetch
			cr.Status.Fetch.StartedAt = metav1.Time{}
			cr.Status.Fetch.UpdatedAt = metav1.Time{}

			// inspect
			if !strings.Contains(cr.Status.Inspect.Stdout, "Resources in app 'test-helm-ctrl'") {
				t.Fatalf("Expected non-empty inspect output: '%s'", cr.Status.Inspect.Stdout)
			}
			cr.Status.Inspect.UpdatedAt = metav1.Time{}
			cr.Status.Inspect.Stdout = ""

			// template
			cr.Status.Template.UpdatedAt = metav1.Time{}
			cr.Status.Template.Stderr = ""
		}

		if !reflect.DeepEqual(expectedStatus, cr.Status) {
			t.Fatalf("Status is not same: %#v vs %#v", expectedStatus, cr.Status)
		}
	})
}
