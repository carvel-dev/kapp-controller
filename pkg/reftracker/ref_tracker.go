// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"fmt"
	"sync"
)

type AppRefTracker struct {
	lock       sync.Mutex
	refsToApps map[RefKey]map[string]struct{}
	appsToRefs map[string]map[RefKey]struct{}
}

func NewAppRefTracker() *AppRefTracker {
	return &AppRefTracker{refsToApps: map[RefKey]map[string]struct{}{}, appsToRefs: map[string]map[RefKey]struct{}{}}
}

func (a AppRefTracker) AddAppForRef(refKey RefKey, appName string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		apps = map[string]struct{}{}
	}

	appKey := fmt.Sprintf(`%s:%s`, appName, refKey.Namespace())
	refs := a.appsToRefs[appKey]
	if refs == nil {
		refs = map[RefKey]struct{}{}
	}

	apps[appName] = struct{}{}
	a.refsToApps[refKey] = apps

	refs[refKey] = struct{}{}
	a.appsToRefs[appKey] = refs
}

func (a AppRefTracker) AppsForRef(refKey RefKey) (map[string]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		return nil, fmt.Errorf("could not find ref %s/%s", refKey.Kind(), refKey.RefName())
	}

	return apps, nil
}

func (a AppRefTracker) RefsForApp(appName, namespace string) (map[RefKey]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)
	if a.appsToRefs[appKey] == nil {
		return nil, fmt.Errorf("could not find refs for App %s", appName)
	}

	return a.appsToRefs[appKey], nil
}

func (a AppRefTracker) RemoveRef(refKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.refsToApps, refKey)
}

func (a AppRefTracker) RemoveAppFromAllRefs(appName, namespace string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)
	refKeys := a.appsToRefs[appKey]
	for refKey := range refKeys {
		apps := a.refsToApps[refKey]
		if apps == nil {
			continue
		}

		delete(apps, appName)
		a.refsToApps[refKey] = apps
	}

	delete(a.appsToRefs, appKey)
}

func (a AppRefTracker) PruneAppFromRefs(currentRefs map[RefKey]struct{}, appName, namespace string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)
	refsInState := a.appsToRefs[appKey]

	// Compare current state against App's
	// previous refs.
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
		delete(apps, appName)
		a.refsToApps[refKey] = apps
	}

	// Make sure appsToRefs uses refs currently
	// on App spec.
	a.appsToRefs[appKey] = currentRefs
}
