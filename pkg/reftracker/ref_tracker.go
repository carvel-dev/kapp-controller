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
}

func NewAppRefTracker() *AppRefTracker {
	refsToApps := map[string]map[string]struct{}{}
	lock := sync.Mutex{}
	return &AppRefTracker{refsToApps: refsToApps, lock: lock}
}

func (a AppRefTracker) AddAppToRefMap(resourceKind, resourceName, namespace, appName string) {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))

	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		// If refKey not found, need to initialize map
		apps = map[string]struct{}{}
	}
	appKey := strings.ToLower(appName)
	apps[appKey] = struct{}{}
	a.refsToApps[refKey] = apps
}

func (a AppRefTracker) GetAppsForRef(resourceKind, resourceName, namespace string) (map[string]struct{}, error) {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.refsToApps[refKey] == nil {
		return nil, fmt.Errorf("could not find App ref %s/%s", resourceKind, resourceName)
	}
	return a.refsToApps[refKey], nil
}

func (a AppRefTracker) CheckAppExistsForRef(resourceKind, resourceName, namespace, appName string) bool {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	appKey := strings.ToLower(appName)

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.refsToApps[refKey] == nil {
		return false
	}
	_, keyExists := a.refsToApps[refKey][appKey]
	return keyExists
}

func (a AppRefTracker) RemoveRefFromMap(resourceKind, resourceName, namespace string) {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.refsToApps, refKey)
}

func (a AppRefTracker) RemoveAppFromRefMap(resourceKind, resourceName, namespace, appName string) error {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))

	a.lock.Lock()
	defer a.lock.Unlock()

	apps := a.refsToApps[refKey]
	if apps == nil {
		return fmt.Errorf("could not find App ref %s/%s", resourceKind, resourceName)
	}

	appKey := strings.ToLower(appName)
	delete(apps, appKey)
	a.refsToApps[refKey] = apps
	return nil
}

type AppUpdateStatus struct {
	lock               sync.Mutex
	appsToUpdateStatus map[string]struct{}
}

func NewAppUpdateStatus() *AppUpdateStatus {
	appsUpdateStatus := map[string]struct{}{}
	lock := sync.Mutex{}
	return &AppUpdateStatus{appsToUpdateStatus: appsUpdateStatus, lock: lock}
}

func (a AppUpdateStatus) MarkNeedsUpdate(appName, namespace string) {
	appKey := strings.ToLower(fmt.Sprintf(`%s:%s`, appName, namespace))

	a.lock.Lock()
	defer a.lock.Unlock()

	a.appsToUpdateStatus[appKey] = struct{}{}
}

func (a AppUpdateStatus) IsUpdateNeeded(appName, namespace string) bool {
	appKey := strings.ToLower(fmt.Sprintf(`%s:%s`, appName, namespace))

	a.lock.Lock()
	defer a.lock.Unlock()

	_, keyExists := a.appsToUpdateStatus[appKey]
	return keyExists
}

func (a AppUpdateStatus) MarkUpdated(appName, namespace string) {
	appKey := strings.ToLower(fmt.Sprintf(`%s:%s`, appName, namespace))

	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.appsToUpdateStatus, appKey)
}
