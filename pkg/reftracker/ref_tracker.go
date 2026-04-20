// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"fmt"
	"sync"
)

type AppRefTracker struct {
	lock       sync.Mutex
	refsToApps map[RefKey]map[RefKey]struct{}
	appsToRefs map[RefKey]map[RefKey]struct{}
}

// TODO: Rename since this doesn't only work with Apps
func NewAppRefTracker() *AppRefTracker {
	return &AppRefTracker{refsToApps: map[RefKey]map[RefKey]struct{}{}, appsToRefs: map[RefKey]map[RefKey]struct{}{}}
}

func (a *AppRefTracker) AppsForRef(refKey RefKey) (map[RefKey]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		return nil, fmt.Errorf("could not find ref %s", refKey.Description())
	}

	// Return a copy. Callers (e.g. SecretHandler.enqueueAppsForUpdate)
	// iterate this map without holding the tracker lock, and concurrent
	// reconciles will call back into ReconcileRefs / RemoveAppFromAllRefs
	// and mutate the same underlying map while the caller is iterating.
	// Without the copy the Go runtime aborts the process with
	// "concurrent map iteration and map write" under load (#1812).
	return cloneRefKeySet(apps), nil
}

func (a *AppRefTracker) RefsForApp(appKey RefKey) (map[RefKey]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	refs := a.appsToRefs[appKey]
	if refs == nil {
		return nil, fmt.Errorf("could not find refs for App %s", appKey.RefName())
	}

	// Same reasoning as AppsForRef: hand back a snapshot so callers can
	// iterate safely while concurrent writers keep mutating the tracker.
	return cloneRefKeySet(refs), nil
}

func cloneRefKeySet(s map[RefKey]struct{}) map[RefKey]struct{} {
	out := make(map[RefKey]struct{}, len(s))
	for k := range s {
		out[k] = struct{}{}
	}
	return out
}

func (a *AppRefTracker) RemoveRef(refKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.refsToApps, refKey)
}

func (a *AppRefTracker) RemoveAppFromAllRefs(appKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	refKeys := a.appsToRefs[appKey]
	for refKey := range refKeys {
		apps := a.refsToApps[refKey]
		if apps == nil {
			continue
		}

		delete(apps, appKey)
		a.refsToApps[refKey] = apps
	}

	delete(a.appsToRefs, appKey)
}

func (a *AppRefTracker) ReconcileRefs(currentRefs map[RefKey]struct{}, appKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Add all new refs to AppRefTracker
	for refKey := range currentRefs {
		apps := a.refsToApps[refKey]
		if apps == nil {
			apps = map[RefKey]struct{}{}
		}

		refs := a.appsToRefs[appKey]
		if refs == nil {
			refs = map[RefKey]struct{}{}
		}

		apps[appKey] = struct{}{}
		a.refsToApps[refKey] = apps

		refs[refKey] = struct{}{}
		a.appsToRefs[appKey] = refs
	}

	// Compare current state against App's
	// previous refs.
	refsInState := a.appsToRefs[appKey]
	var refsToRemove []RefKey
	for refKey := range refsInState {
		if _, refExists := currentRefs[refKey]; !refExists {
			refsToRemove = append(refsToRemove, refKey)
		}
	}

	// Remove any differences between App's
	// current state and previous state
	for _, refKey := range refsToRemove {
		apps := a.refsToApps[refKey]
		delete(apps, appKey)
		a.refsToApps[refKey] = apps
	}

	// Make sure appsToRefs uses refs currently
	// on App spec.
	a.appsToRefs[appKey] = currentRefs
}
