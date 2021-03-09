// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"fmt"
)

type AppRefTracker struct {
	refsToApps         map[string]map[string]string
	appsToUpdateStatus map[string]bool
}

func NewAppRefTracker() AppRefTracker {
	refsToApps := make(map[string]map[string]string)
	appsToUpdateStatus := make(map[string]bool)
	return AppRefTracker{refsToApps: refsToApps, appsToUpdateStatus: appsToUpdateStatus}
}

func (a AppRefTracker) AddAppToRefMap(resourceKind, resourceName, namespace, appName string) {
	// TODO: lowercase all except app name
	key := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	apps := a.GetAppsForRef(resourceKind, resourceName, namespace)
	if apps == nil {
		apps = make(map[string]string)
	}
	apps[appName] = appName
	a.refsToApps[key] = apps
}

func (a AppRefTracker) GetAppsForRef(resourceKind, resourceName, namespace string) map[string]string {
	key := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	return a.refsToApps[key]
}

func (a AppRefTracker) CheckAppExistsForRef(resourceKind, resourceName, namespace, appName string) bool {
	key := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	if a.refsToApps[key][appName] == "" {
		return false
	}
	return true
}

func (a AppRefTracker) RemoveRefFromMap(resourceKind, resourceName, namespace string) {
	key := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	delete(a.refsToApps, key)
}

func (a AppRefTracker) RemoveAppFromRefMap(resourceKind, resourceName, namespace, appName string) error {
	if !a.CheckAppExistsForRef(resourceKind, resourceName, namespace, appName) {
		return fmt.Errorf("could not find App %s for ref %s/%s", appName, resourceKind, resourceName)
	}

	apps := a.GetAppsForRef(resourceKind, resourceName, namespace)
	appKey := fmt.Sprintf(`%s`, appName)
	delete(apps, appKey)
	refKey := fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace)
	a.refsToApps[refKey] = apps

	return nil
}

func (a AppRefTracker) MarkAppForUpdate(appName, namespace string) {
	key := fmt.Sprintf(`%s:%s`, appName, namespace)
	a.appsToUpdateStatus[key] = true
}

func (a AppRefTracker) MarkAppUpdated(appName, namespace string) {
	key := fmt.Sprintf(`%s:%s`, appName, namespace)
	a.appsToUpdateStatus[key] = false
}

func (a AppRefTracker) GetAppUpdateStatus(appName, namespace string) bool {
	key := fmt.Sprintf(`%s:%s`, appName, namespace)
	return a.appsToUpdateStatus[key]
}

func (a AppRefTracker) RemoveAppFromUpdateMap(appName, namespace string) {
	key := fmt.Sprintf(`%s:%s`, appName, namespace)
	delete(a.appsToUpdateStatus, key)
}
