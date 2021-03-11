package reftracker

import (
	"fmt"
	"sync"
)

type AppUpdateStatus struct {
	lock               sync.Mutex
	appsToUpdateStatus map[string]struct{}
}

func NewAppUpdateStatus() *AppUpdateStatus {
	return &AppUpdateStatus{appsToUpdateStatus: map[string]struct{}{}}
}

func (a AppUpdateStatus) MarkNeedsUpdate(appName, namespace string) {
	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)

	a.lock.Lock()
	defer a.lock.Unlock()

	a.appsToUpdateStatus[appKey] = struct{}{}
}

func (a AppUpdateStatus) IsUpdateNeeded(appName, namespace string) bool {
	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)

	a.lock.Lock()
	defer a.lock.Unlock()

	_, keyExists := a.appsToUpdateStatus[appKey]
	return keyExists
}

func (a AppUpdateStatus) MarkUpdated(appName, namespace string) {
	appKey := fmt.Sprintf(`%s:%s`, appName, namespace)

	a.lock.Lock()
	defer a.lock.Unlock()

	delete(a.appsToUpdateStatus, appKey)
}
