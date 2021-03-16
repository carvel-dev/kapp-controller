// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import "testing"

func Test_AddAppForRef_AddsApp_WhenRefNotInMap(t *testing.T) {
	appRefTracker := NewAppRefTracker()

	refKey := RefKey{"secret", "secretName", "default"}
	appRefTracker.AddAppForRef(refKey, "app")

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := apps["app"]; !ok {
		t.Fatalf("app was not added to appRefTracker when ref key did not exist")
	}

	refs, err := appRefTracker.RefsForApp("app", "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := refs[refKey]; !ok {
		t.Fatalf("ref was not added to appRefTracker when App key did not exist")
	}
}

func Test_RemoveAppFromAllRefs_RemovesApp(t *testing.T) {
	appRefTracker := NewAppRefTracker()

	refKey := RefKey{"secret", "secretName", "default"}
	appRefTracker.AddAppForRef(refKey, "app")

	appRefTracker.RemoveAppFromAllRefs("app", "default")

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := apps["app"]; ok {
		t.Fatalf("expected app to be removed from appRefTracker after deletion")
	}
}
