// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package reftracker_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"carvel.dev/kapp-controller/pkg/reftracker"
)

// Test_AppsForRef_ConcurrentIterationAndWrite_DoesNotPanic reproduces the race
// reported in #1812/#1843: AppsForRef returned the internal map by reference, so
// a caller iterating the result (e.g. SecretHandler.enqueueAppsForUpdate) could
// run concurrently with another goroutine mutating the same map via
// ReconcileRefs/RemoveAppFromAllRefs, causing a
// "fatal error: concurrent map iteration and map write" crash.
//
// Before the fix, which makes AppsForRef return a copy, `go test -race` reports
// the race on every run, while a plain run only sometimes crashes.
func Test_AppsForRef_ConcurrentIterationAndWrite_DoesNotPanic(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()
	appKey := reftracker.NewAppKey("app", "default")

	seedRefs := map[reftracker.RefKey]struct{}{}
	for i := 0; i < 50; i++ {
		seedRefs[reftracker.NewSecretKey(fmt.Sprintf("secret-%d", i), "default")] = struct{}{}
	}
	appRefTracker.ReconcileRefs(seedRefs, appKey)

	var wg sync.WaitGroup
	stop := make(chan struct{})
	iterated := 0

	// Writer: add and remove a second App, so each ref's app set keeps being
	// written to but never becomes empty for the reader.
	otherAppKey := reftracker.NewAppKey("other-app", "default")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			appRefTracker.ReconcileRefs(seedRefs, otherAppKey)
			appRefTracker.RemoveAppFromAllRefs(otherAppKey)
		}
	}()

	// Reader: read a ref's app set and iterate it outside the lock, mimicking
	// SecretHandler.enqueueAppsForUpdate.
	wg.Add(1)
	go func() {
		defer wg.Done()
		refKey := reftracker.NewSecretKey("secret-0", "default")
		for {
			select {
			case <-stop:
				return
			default:
			}
			apps, err := appRefTracker.AppsForRef(refKey)
			if err != nil {
				continue
			}
			for range apps {
				iterated++
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()

	if iterated == 0 {
		t.Fatal("the reader never iterated a result, so the race was not exercised")
	}
}

// Test_RefsForApp_ConcurrentIterationAndWrite_DoesNotPanic covers the symmetric
// path: RefsForApp also returned its internal map by reference.
func Test_RefsForApp_ConcurrentIterationAndWrite_DoesNotPanic(t *testing.T) {
	appRefTracker := reftracker.NewAppRefTracker()
	appKey := reftracker.NewAppKey("app", "default")

	seedRefs := map[reftracker.RefKey]struct{}{}
	for i := 0; i < 50; i++ {
		seedRefs[reftracker.NewSecretKey(fmt.Sprintf("secret-%d", i), "default")] = struct{}{}
	}
	appRefTracker.ReconcileRefs(seedRefs, appKey)

	var wg sync.WaitGroup
	stop := make(chan struct{})
	iterated := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			appRefTracker.ReconcileRefs(seedRefs, appKey)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			refs, err := appRefTracker.RefsForApp(appKey)
			if err != nil {
				continue
			}
			for range refs {
				iterated++
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()

	if iterated == 0 {
		t.Fatal("the reader never iterated a result, so the race was not exercised")
	}
}
