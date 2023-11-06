// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app_test

import (
	"reflect"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	apppkg "carvel.dev/kapp-controller/pkg/app"
	"carvel.dev/kapp-controller/pkg/reftracker"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_AppRefTracker_HasAppRemovedForSecrets_ThatAreNoLongerUsedByApp(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()
	// Add secrets to appRefTracker and have all of them
	// be aware of App named "app"
	keySecretName := reftracker.NewSecretKey("secretName", "default")
	keySecretName2 := reftracker.NewSecretKey("secretName2", "default")
	keySecretName3 := reftracker.NewSecretKey("secretName3", "default")
	refKeyMap := map[reftracker.RefKey]struct{}{
		keySecretName:  {},
		keySecretName2: {},
		keySecretName3: {},
	}

	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
	}

	appKey := reftracker.NewAppKey(app.Name, app.Namespace)
	appRefTracker.ReconcileRefs(refKeyMap, appKey)

	ar := apppkg.NewReconciler(nil, testr.New(t), apppkg.CRDAppFactory{}, appRefTracker, nil, FakeComponentInfo{})

	// This map represents the secrets the App has on its spec
	refMap := map[reftracker.RefKey]struct{}{
		keySecretName: {},
	}

	// We expect this method will clean up the appRefTracker
	// if the App above is no longer using a secret that it
	// once did.
	ar.UpdateAppRefs(refMap, &app)

	appRefTracker = ar.AppRefTracker()
	expected := map[reftracker.RefKey]struct{}{}
	out, _ := appRefTracker.AppsForRef(keySecretName2)
	if !reflect.DeepEqual(out, expected) {
		t.Fatalf("\nExpected: %s\nGot: %s", expected, out)
	}

	out, _ = appRefTracker.AppsForRef(keySecretName3)
	assert.Truef(t, reflect.DeepEqual(out, expected), "Expected: %s\nGot: %s", expected, out)

	expected = map[reftracker.RefKey]struct{}{
		appKey: {},
	}
	out, _ = appRefTracker.AppsForRef(keySecretName)
	assert.Truef(t, reflect.DeepEqual(out, expected), "Expected: %s\nGot: %s", expected, out)
}

func Test_AppRefTracker_HasNoAppsRemoved_WhenRefsRemainSame(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()
	// Add secrets to appRefTracker and have all of them
	// be aware of App named "app"
	keySecretName := reftracker.NewSecretKey("secretName", "default")
	keySecretName2 := reftracker.NewSecretKey("secretName2", "default")
	keySecretName3 := reftracker.NewSecretKey("secretName3", "default")
	refKeyMap := map[reftracker.RefKey]struct{}{
		keySecretName:  {},
		keySecretName2: {},
		keySecretName3: {},
	}

	app := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
		},
	}

	appKey := reftracker.NewAppKey(app.Name, app.Namespace)
	appRefTracker.ReconcileRefs(refKeyMap, appKey)

	ar := apppkg.NewReconciler(nil, testr.New(t), apppkg.CRDAppFactory{}, appRefTracker, nil, FakeComponentInfo{})

	// This map represents the secrets the App has
	// on its spec
	refMap := map[reftracker.RefKey]struct{}{
		keySecretName:  {},
		keySecretName2: {},
		keySecretName3: {},
	}

	ar.UpdateAppRefs(refMap, &app)

	// Expect all refs to be associated with app
	expected := map[reftracker.RefKey]struct{}{
		appKey: {},
	}

	for refKey := range refMap {
		out, _ := ar.AppRefTracker().AppsForRef(refKey)
		assert.Truef(t, reflect.DeepEqual(out, expected), "Expected: %s\nGot: %s", expected, out)
	}
}
