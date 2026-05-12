// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package reftracker_test

import (
	"sync"
	"testing"

	"carvel.dev/kapp-controller/pkg/reftracker"
)

func Test_AddAppForRef_AddsApp_WhenRefNotInMap(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewSecretKey("secretName", "default")
	refKeyMap := map[reftracker.RefKey]struct{}{
		refKey: {},
	}
	appKey := reftracker.NewAppKey("app", "default")
	appRefTracker.ReconcileRefs(refKeyMap, appKey)

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := apps[appKey]; !ok {
		t.Fatalf("app was not added to appRefTracker when ref key did not exist")
	}

	refs, err := appRefTracker.RefsForApp(reftracker.NewAppKey("app", "default"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := refs[refKey]; !ok {
		t.Fatalf("ref was not added to appRefTracker when App key did not exist")
	}
}

func Test_RemoveAppFromAllRefs_RemovesApp(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewSecretKey("secretName", "default")
	refKeyMap := map[reftracker.RefKey]struct{}{
		refKey: {},
	}
	appKey := reftracker.NewAppKey("app", "default")
	appRefTracker.ReconcileRefs(refKeyMap, appKey)

	appRefTracker.RemoveAppFromAllRefs(appKey)

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := apps[appKey]; ok {
		t.Fatalf("expected app to be removed from appRefTracker after deletion")
	}
}

func Test_AppsForRef_SafeForConcurrentIteration(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewSecretKey("secretName", "default")
	refKeyMap := map[reftracker.RefKey]struct{}{refKey: {}}

	for i := 0; i < 20; i++ {
		appKey := reftracker.NewAppKey("seed-app", "default")
		appKey = reftracker.NewAppKey(appKey.RefName()+string(rune('a'+i)), "default")
		appRefTracker.ReconcileRefs(refKeyMap, appKey)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				appKey := reftracker.NewAppKey("mutator-app", "default")
				appKey = reftracker.NewAppKey(appKey.RefName()+string(rune('a'+(i%26))), "default")
				appRefTracker.ReconcileRefs(refKeyMap, appKey)
				appRefTracker.RemoveAppFromAllRefs(appKey)
				i++
			}
		}
	}()

	iterDone := make(chan struct{})
	go func() {
		defer close(iterDone)
		for i := 0; i < 200; i++ {
			apps, err := appRefTracker.AppsForRef(refKey)
			if err != nil {
				continue
			}
			for range apps {
			}
		}
	}()

	<-iterDone
	close(stop)
	wg.Wait()
}
