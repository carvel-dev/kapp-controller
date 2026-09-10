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
// Before the fix this test crashes the process; after it (AppsForRef returns a
// copy) it passes. Run with -race for extra signal.
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

	// Writer: continuously mutate the internal maps under the lock.
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
			appRefTracker.RemoveAppFromAllRefs(appKey)
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
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
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
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}
