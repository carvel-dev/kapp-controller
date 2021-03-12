// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"fmt"
	"strings"
	"sync"
)

type AppRefTracker struct {
	lock       sync.Mutex
	refsToApps map[string]map[string]struct{}
	appsToRefs map[string]map[string]struct{}
}

func NewAppRefTracker() *AppRefTracker {
	return &AppRefTracker{refsToApps: map[string]map[string]struct{}{}, appsToRefs: map[string]map[string]struct{}{}}
}

func (a AppRefTracker) AddAppForRef(resourceKind, refName, namespace, appName string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, refName, namespace)
	apps := a.refsToApps[refKey]
	if apps == nil {
		apps = map[string]struct{}{}
	}

	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)
	refs := a.appsToRefs[appKey]
	if refs == nil {
		refs = map[string]struct{}{}
	}

	apps[appName] = struct{}{}
	a.refsToApps[refKey] = apps

	appsRefKey := fmt.Sprintf(`%s:%s`, refName, resourceKind)
	refs[appsRefKey] = struct{}{}
	a.appsToRefs[appKey] = refs
}

func (a AppRefTracker) AppsForRef(resourceKind, resourceName, namespace string) (map[string]struct{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	if a.refsToApps[refKey] == nil {
		return nil, fmt.Errorf("could not find App ref %s/%s", resourceKind, resourceName)
	}

	return a.refsToApps[refKey], nil
}

func (a AppRefTracker) RemoveRef(resourceKind, resourceName, namespace string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	delete(a.refsToApps, refKey)
}

func (a AppRefTracker) RemoveAppFromAllRefs(refs map[string]struct{}, resourceKind, namespace, appName string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for ref := range refs {
		refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, ref, namespace)
		apps := a.refsToApps[refKey]
		if apps == nil {
			continue
		}

		delete(apps, appName)
		a.refsToApps[refKey] = apps
	}

	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)
	delete(a.appsToRefs, appKey)
}

func (a AppRefTracker) PruneAppRefs(currentRefs map[string]struct{}, resourceKind, namespace, appName string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)
	refsInState := a.appsToRefs[appKey]

	// Compare current state against App's
	// previous refs.
	var diff []string
	for ref := range refsInState {
		// The format of ref is refName:kind
		// so we need to  get the refName portion
		// to build the key for the refsToApps map.
		index := strings.Index(ref, ":")
		if index == -1 {
			continue
		}
		refName := ref[:index]
		refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, refName, namespace)

		if _, appExists := currentRefs[refName]; !appExists {
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

	// Update App state to reflect new refs
	refs := map[string]struct{}{}
	for refName := range currentRefs {
		appsRefKey := fmt.Sprintf(`%s:%s`, refName, resourceKind)
		refs[appsRefKey] = struct{}{}
	}
	a.appsToRefs[appKey] = refs
}
