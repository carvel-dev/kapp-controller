package controller_test

import (
	"reflect"
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/cmd/controller"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_AppRefTracker_HasAppRemovedForSecrets_ThatAreNoLongerUsedByApp(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()
	// Add secrets to AppRefTracker and have all of them
	// be aware of App named "app"
	appRefTracker.AddAppForRef("secret", "secretName", "default", "app")
	appRefTracker.AddAppForRef("secret", "secretName2", "default", "app")
	appRefTracker.AddAppForRef("secret", "secretName3", "default", "app")

	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
	}

	ar := controller.AppsReconciler{
		AppRefTracker: appRefTracker,
	}

	// This map represents the secrets the App has on its spec
	refMap := map[string]struct{}{
		"secretName": {},
	}

	// We expect this method will clean up the AppRefTracker
	// if the App above is no longer using a secret that it
	// once did.
	ar.UpdateAppRefs(refMap, "secret", &app)

	expected := map[string]struct{}{}
	out, _ := ar.AppRefTracker.AppsForRef("secret", "secretName2", "default")
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", expected, out)
	}

	out, _ = ar.AppRefTracker.AppsForRef("secret", "secretName3", "default")
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", expected, out)
	}

	expected = map[string]struct{}{
		"app": {},
	}
	out, _ = ar.AppRefTracker.AppsForRef("secret", "secretName", "default")
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", expected, out)
	}
}

func Test_AppRefTracker_HasNoAppsRemoved_WhenRefsRemainSame(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()
	// Add secrets to AppRefTracker and have all of them
	// be aware of App named "app"
	appRefTracker.AddAppForRef("secret", "secretName", "default", "app")
	appRefTracker.AddAppForRef("secret", "secretName2", "default", "app")
	appRefTracker.AddAppForRef("secret", "secretName3", "default", "app")

	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
	}

	ar := controller.AppsReconciler{
		AppRefTracker: appRefTracker,
	}

	// This map represents the secrets the App has
	// on its spec
	refMap := map[string]struct{}{
		"secretName":  {},
		"secretName2": {},
		"secretName3": {},
	}

	ar.UpdateAppRefs(refMap, "secret", &app)

	// Expect all refs to be associated with app
	expected := map[string]struct{}{
		"app": {},
	}
	out, _ := ar.AppRefTracker.AppsForRef("secret", "secretName2", "default")
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", expected, out)
	}

	out, _ = ar.AppRefTracker.AppsForRef("secret", "secretName3", "default")
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", expected, out)
	}

	out, _ = ar.AppRefTracker.AppsForRef("secret", "secretName", "default")
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", refMap, out)
	}
}
