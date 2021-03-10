// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import "testing"

func Test_AddAppToRefMap_AddsApp_WhenRefNotInMap(t *testing.T) {
	appRefTracker := NewAppRefTracker()
	appRefTracker.AddAppToRefMap("secret", "secretName", "default", "app")

	if !appRefTracker.CheckAppExistsForRef("secret", "secretName", "default", "app") {
		t.Fatalf("app was not added to AppRefTracker when Ref key did not exist")
	}
}

func Test_RemoveAppFromRefMap_RemovesApp(t *testing.T) {
	appRefTracker := NewAppRefTracker()
	appRefTracker.AddAppToRefMap("secret", "secretName", "default", "app")
	appRefTracker.RemoveAppFromRefMap("secret", "secretName", "default", "app")

	if appRefTracker.CheckAppExistsForRef("secret", "secretName", "default", "app") {
		t.Fatalf("expected app to be removed from AppRefTracker after deletion")
	}
}
