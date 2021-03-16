// Copyright 2021 VMware, Inc.
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

func NewAppRefTracker() *AppRefTracker {
	return &AppRefTracker{refsToApps: map[RefKey]map[RefKey]struct{}{}, appsToRefs: map[RefKey]map[RefKey]struct{}{}}
}

func (a AppRefTracker) AddAppForRef(refKey RefKey, appName string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		apps = map[RefKey]struct{}{}
	}

	appKey := NewAppKey(appName, refKey.Namespace())
	refs := a.appsToRefs[appKey]
	if refs == nil {
		refs = map[RefKey]struct{}{}
	}

	apps[appKey] = struct{}{}
	a.refsToApps[refKey] = apps

	refs[refKey] = struct{}{}
	a.appsToRefs[appKey] = refs
}

func (a AppRefTracker) AppsForRef(refKey RefKey) (map[RefKey]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		return nil, fmt.Errorf("could not find ref %s", refKey.Description())
	}

	return apps, nil
}

func (a AppRefTracker) RefsForApp(appKey RefKey) (map[RefKey]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.appsToRefs[appKey] == nil {
		return nil, fmt.Errorf("could not find refs for App %s", appKey.RefName())
	}

	return a.appsToRefs[appKey], nil
}

func (a AppRefTracker) RemoveRef(refKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.refsToApps, refKey)
}

func (a AppRefTracker) RemoveAppFromAllRefs(appKey RefKey) {
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

func (a AppRefTracker) ReconcileRefs(currentRefs map[RefKey]struct{}, appKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Add all new refs to AppRefTracker
	for refKey := range currentRefs {
		apps := a.refsToApps[refKey]
		if apps == nil {
			apps = map[RefKey]struct{}{}
		}

		appKey := NewAppKey(appKey.RefName(), refKey.Namespace())
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
	var diff []RefKey
	for refKey := range refsInState {
		if _, refExists := currentRefs[refKey]; !refExists {
			diff = append(diff, refKey)
		}
	}

	// Remove any differences between App's
	// current state and previous state
	for _, refKey := range diff {
		apps := a.refsToApps[refKey]
		delete(apps, appKey)
		a.refsToApps[refKey] = apps
	}

	// Make sure appsToRefs uses refs currently
	// on App spec.
	a.appsToRefs[appKey] = currentRefs
}
