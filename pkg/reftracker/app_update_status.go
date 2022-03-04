// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package reftracker

import (
	"sync"
)

type AppUpdateStatus struct {
	lock               sync.Mutex
	appsToUpdateStatus map[RefKey]struct{}
}

func NewAppUpdateStatus() *AppUpdateStatus {
	return &AppUpdateStatus{appsToUpdateStatus: map[RefKey]struct{}{}}
}

func (a *AppUpdateStatus) MarkNeedsUpdate(appKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.appsToUpdateStatus[appKey] = struct{}{}
}

func (a *AppUpdateStatus) IsUpdateNeeded(appKey RefKey) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	_, keyExists := a.appsToUpdateStatus[appKey]
	return keyExists
}

func (a *AppUpdateStatus) MarkUpdated(appKey RefKey) {
	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.appsToUpdateStatus, appKey)
}
