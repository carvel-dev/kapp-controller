// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"fmt"
	"strings"
)

var appExists struct{}

type AppRefTracker struct {
	refsToApps         map[string]map[string]struct{}
	appsToUpdateStatus map[string]struct{}
}

func NewAppRefTracker() AppRefTracker {
	refsToApps := make(map[string]map[string]struct{})
	appsToUpdateStatus := make(map[string]struct{})
	return AppRefTracker{refsToApps: refsToApps, appsToUpdateStatus: appsToUpdateStatus}
}

func (a AppRefTracker) AddAppToRefMap(resourceKind, resourceName, namespace, appName string) {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	apps, err := a.GetAppsForRef(resourceKind, resourceName, namespace)
	if err != nil {
		// If refKey not found, need to initialize map
		apps = make(map[string]struct{})
	}
	appKey := strings.ToLower(appName)
	apps[appKey] = appExists
	a.refsToApps[refKey] = apps
}

func (a AppRefTracker) GetAppsForRef(resourceKind, resourceName, namespace string) (map[string]struct{}, error) {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	if a.refsToApps[refKey] == nil {
		return nil, fmt.Errorf("could not find App ref %s/%s", resourceKind, resourceName)
	}
	return a.refsToApps[refKey], nil
}

func (a AppRefTracker) CheckAppExistsForRef(resourceKind, resourceName, namespace, appName string) bool {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	appKey := strings.ToLower(appName)
	if a.refsToApps[refKey] == nil {
		return false
	}
	_, keyExists := a.refsToApps[refKey][appKey]
	return keyExists
}

func (a AppRefTracker) RemoveRefFromMap(resourceKind, resourceName, namespace string) {
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	delete(a.refsToApps, refKey)
}

func (a AppRefTracker) RemoveAppFromRefMap(resourceKind, resourceName, namespace, appName string) error {
	apps, err := a.GetAppsForRef(resourceKind, resourceName, namespace)
	if err != nil {
		return err
	}

	appKey := strings.ToLower(appName)
	delete(apps, appKey)
	refKey := strings.ToLower(fmt.Sprintf(`%s:%s:%s`, resourceKind, resourceName, namespace))
	a.refsToApps[refKey] = apps

	return nil
}

func (a AppRefTracker) MarkAppForUpdate(appName, namespace string) {
	appKey := strings.ToLower(fmt.Sprintf(`%s:%s`, appName, namespace))
	a.appsToUpdateStatus[appKey] = appExists
}

func (a AppRefTracker) MarkAppUpdated(appName, namespace string) {
	a.RemoveAppFromUpdateMap(appName, namespace)
}

func (a AppRefTracker) GetAppUpdateStatus(appName, namespace string) bool {
	appKey := strings.ToLower(fmt.Sprintf(`%s:%s`, appName, namespace))
	_, keyExists := a.appsToUpdateStatus[appKey]
	return keyExists
}

func (a AppRefTracker) RemoveAppFromUpdateMap(appName, namespace string) {
	appKey := strings.ToLower(fmt.Sprintf(`%s:%s`, appName, namespace))
	delete(a.appsToUpdateStatus, appKey)
}
