// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package watchers

import (
	"testing"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	"k8s.io/apimachinery/pkg/watch"
)

func TestTranslationWatcher(t *testing.T) {
	proxiedWatcher := watch.NewFake()

	w := NewTranslationWatcher(func(e watch.Event) watch.Event {
		pkg := e.Object.(*datapackaging.Package)
		pkg.Name = "transformed"
		return e
	}, func(e watch.Event) bool {
		if e.Type == watch.Error {
			return false
		}
		return true
	}, proxiedWatcher)

	proxiedWatcher.Add(&datapackaging.Package{})
	event := <-w.ResultChan()
	pkg := event.Object.(*datapackaging.Package)
	if pkg.Name != "transformed" {
		t.Fatal("expected object to be transformed")
	}

	proxiedWatcher.Error(&datapackaging.Package{})
	proxiedWatcher.Delete(&datapackaging.Package{})
	event = <-w.ResultChan()
	if event.Type == watch.Error {
		t.Fatal("expected error event to be filtered out")
	}

	proxiedWatcher.Stop()
	w.Stop()
}
