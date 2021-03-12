// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import "testing"

func Test_AddAppForRef_AddsApp_WhenRefNotInMap(t *testing.T) {
	appRefTracker := NewAppRefTracker()
	appRefTracker.AddAppForRef("secret", "secretName", "default", "app")

	if _, ok := appRefTracker.refsToApps["secret:secretName:default"]["app"]; !ok {
		t.Fatalf("app was not added to AppRefTracker when ref key did not exist")
	}

	if _, ok := appRefTracker.appsToRefs["app:default"]["secretName:secret"]; !ok {
		t.Fatalf("ref was not added to AppRefTracker when App key did not exist")
	}
}

func Test_RemoveAppFromAllRefs_RemovesApp(t *testing.T) {
	appRefTracker := NewAppRefTracker()
	appRefTracker.AddAppForRef("secret", "secretName", "default", "app")

	refs := map[string]struct{}{
		"secretName": struct{}{},
	}
	appRefTracker.RemoveAppFromAllRefs(refs, "secret", "default", "app")

	if _, ok := appRefTracker.refsToApps["secret:secretName:default"]["app"]; ok {
		t.Fatalf("expected app to be removed from AppRefTracker after deletion")
	}
}
