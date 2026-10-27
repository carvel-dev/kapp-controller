// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package reftracker_test

import (
	"runtime"
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

func Test_AppsForRef_RefsForApp_ReturnSnapshot(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewSecretKey("secretName", "default")
	appKey := reftracker.NewAppKey("app", "default")
	appRefTracker.ReconcileRefs(map[reftracker.RefKey]struct{}{refKey: {}}, appKey)

	apps, err := appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	refs, err := appRefTracker.RefsForApp(appKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	appRefTracker.RemoveAppFromAllRefs(appKey)

	if _, ok := apps[appKey]; !ok {
		t.Fatalf("expected returned apps to be unaffected by later tracker changes")
	}
	if _, ok := refs[refKey]; !ok {
		t.Fatalf("expected returned refs to be unaffected by later tracker changes")
	}

	delete(apps, appKey)
	appRefTracker.ReconcileRefs(map[reftracker.RefKey]struct{}{refKey: {}}, appKey)
	apps, err = appRefTracker.AppsForRef(refKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := apps[appKey]; !ok {
		t.Fatalf("expected tracker to be unaffected by changes to a returned map")
	}
}

// Run with -race to detect unsynchronized access to the tracker's internal maps.
func Test_AppsForRef_ConcurrentWithReconcileRefs(_ *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()

	refKey := reftracker.NewSecretKey("secretName", "default")
	appKeys := []reftracker.RefKey{
		reftracker.NewAppKey("app1", "default"),
		reftracker.NewAppKey("app2", "default"),
	}
	for _, appKey := range appKeys {
		appRefTracker.ReconcileRefs(map[reftracker.RefKey]struct{}{refKey: {}}, appKey)
	}

	start := make(chan struct{})
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(done)
		<-start
		for i := 0; i < 1000; i++ {
			appKey := appKeys[i%len(appKeys)]
			appRefTracker.RemoveAppFromAllRefs(appKey)
			appRefTracker.ReconcileRefs(map[reftracker.RefKey]struct{}{refKey: {}}, appKey)
			runtime.Gosched()
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		for {
			select {
			case <-done:
				return
			default:
			}
			apps, _ := appRefTracker.AppsForRef(refKey)
			for range apps {
				runtime.Gosched()
			}
			refs, _ := appRefTracker.RefsForApp(appKeys[0])
			for range refs {
				runtime.Gosched()
			}
		}
	}()

	close(start)
	wg.Wait()
}
