// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"fmt"
	"sync"
)

type AppRefTracker struct {
	lock       sync.Mutex
	refsToApps map[string]map[string]struct{}
}

func NewAppRefTracker() *AppRefTracker {
	return &AppRefTracker{refsToApps: map[string]map[string]struct{}{}}
}

func (a AppRefTracker) AddAppForRef(resourceKind, resourceName, namespace, appName string) {
	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)

	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		// If refKey not found, need to initialize map
		apps = map[string]struct{}{}
	}
	appKey := appName
	apps[appKey] = struct{}{}
	a.refsToApps[refKey] = apps
}

func (a AppRefTracker) AppsForRef(resourceKind, resourceName, namespace string) (map[string]struct{}, error) {
	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.refsToApps[refKey] == nil {
		return nil, fmt.Errorf("could not find App ref %s/%s", resourceKind, resourceName)
	}
	return a.refsToApps[refKey], nil
}

func (a AppRefTracker) RemoveRef(resourceKind, resourceName, namespace string) {
	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	a.lock.Lock()
	defer a.lock.Unlock()

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
}
