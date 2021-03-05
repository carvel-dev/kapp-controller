// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package resourcetracker

import (
	"fmt"
)

type AppSecretsEntry struct {
	appName         string
	resourceVersion string
}

type AppSecrets struct {
	secretToApps map[string][]AppSecretsEntry
}

func NewAppsSecretsEntry(name string, resourveVersion string) AppSecretsEntry {
	return AppSecretsEntry{name, resourveVersion}
}

func (a AppSecretsEntry) GetResourceVersion() string {
	return a.resourceVersion
}

func (a AppSecretsEntry) GetAppName() string {
	return a.appName
}

func NewAppSecrets() AppSecrets {
	secretToApps := make(map[string][]AppSecretsEntry)
	return AppSecrets{secretToApps: secretToApps}
}

func (a AppSecrets) AddAppToMap(secretName, namespace, appName, resourceVersion string) {
	secretKey := fmt.Sprintf(`%s:%s`, secretName, namespace)
	appEntry := a.GetAppsForSecret(secretName, namespace)
	a.secretToApps[secretKey] = append(appEntry, NewAppsSecretsEntry(appName, resourceVersion))
}

func (a AppSecrets) GetAppsForSecret(secretName, namespace string) []AppSecretsEntry {
	secretKey := fmt.Sprintf(`%s:%s`, secretName, namespace)
	return a.secretToApps[secretKey]
}

func (a AppSecrets) GetSpecificAppForSecret(secretName, namespace, appName string) (AppSecretsEntry, error) {
	secretKey := fmt.Sprintf(`%s:%s`, secretName, namespace)
	for _, appEntry := range a.secretToApps[secretKey] {
		if appEntry.GetAppName() == appName {
			return appEntry, nil
		}
	}
	return AppSecretsEntry{}, fmt.Errorf("could not find App %s", appName)
}

func (a AppSecrets) RemoveSecretFromMap(secretName, namespace string) {
	secretKey := fmt.Sprintf(`%s:%s`, secretName, namespace)
	delete(a.secretToApps, secretKey)
}

// TODO: Use when App is deleted?
func (a AppSecrets) RemoveAppFromMap(secretName, namespace, appName string) {
	secretKey := fmt.Sprintf(`%s:%s`, secretName, namespace)
	appEntries := a.GetAppsForSecret(secretKey, namespace)
	for i, appEntry := range appEntries {
		if appEntry.GetAppName() == appName {
			appEntries[i] = appEntries[len(appEntries)-1]
			appEntries[len(appEntries)-1] = AppSecretsEntry{}
			appEntries = appEntries[:len(appEntries)-1]
		}
	}
}

func (a AppSecrets) UpdateAppInMap(secretName, namespace, appName, resourceVersion string) {
	secretKey := fmt.Sprintf(`%s:%s`, secretName, namespace)
	appEntries := a.GetAppsForSecret(secretKey, namespace)
	for i, appEntry := range appEntries {
		if appEntry.GetAppName() == appName {
			appEntries[i] = NewAppsSecretsEntry(appName, resourceVersion)
		}
	}
	a.secretToApps[secretKey] = appEntries
}
