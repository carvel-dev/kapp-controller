// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker_test

import (
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/reftracker"
)

func Test_AddAppForRef_AddsApp_WhenRefNotInMap(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewRefKey("secret", "secretName", "default")
	appRefTracker.AddAppForRef(refKey, "app")

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	appKey := reftracker.NewRefKey("app", "app", "default")
	if _, ok := apps[appKey]; !ok {
		t.Fatalf("app was not added to appRefTracker when ref key did not exist")
	}

	refs, err := appRefTracker.RefsForApp(reftracker.NewRefKey("app", "app", "default"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := refs[refKey]; !ok {
		t.Fatalf("ref was not added to appRefTracker when App key did not exist")
	}
}

func Test_RemoveAppFromAllRefs_RemovesApp(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewRefKey("secret", "secretName", "default")
	appRefTracker.AddAppForRef(refKey, "app")

	appKey := reftracker.NewRefKey("app", "app", "default")
	appRefTracker.RemoveAppFromAllRefs(appKey)

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := apps[appKey]; ok {
		t.Fatalf("expected app to be removed from appRefTracker after deletion")
	}
}
