// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type CRDAppWatcher struct {
	app       v1alpha1.App
	appClient kcclient.Interface
}

func NewCRDAppWatcher(app v1alpha1.App, appClient kcclient.Interface) CRDAppWatcher {
	return CRDAppWatcher{app, appClient}
}

func (w CRDAppWatcher) Watch(callback func(v1alpha1.App), cancelCh chan struct{}) error {
	// canceled before starting
	select {
	case <-cancelCh:
		return nil
	default:
	}

	for {
		retry, err := w.watch(callback, cancelCh)
		if err != nil {
			return err
		}
		if !retry {
			return nil
		}
	}
}

func (w CRDAppWatcher) watch(callback func(v1alpha1.App), cancelCh chan struct{}) (bool, error) {
	listOpts := metav1.ListOptions{
		// Only single resource that has ns+name combination.
		// metadata.uid cannot be used as it's not indexed.
		FieldSelector: fields.OneTermEqualSelector("metadata.name", string(w.app.Name)).String(),
	}

	watcher, err := w.appClient.KappctrlV1alpha1().Apps(w.app.Namespace).Watch(context.Background(), listOpts)
	if err != nil {
		return false, fmt.Errorf("Creating app watcher: %s", err)
	}

	defer func() {
		watcher.Stop()

		// Watcher may be trying to send Event before being channel is closed
		// https://github.com/kubernetes/apimachinery/blob/d8530e6c952f75365336be8ea29cfd758ce49ee8/pkg/watch/streamwatcher.go#L127
		// (Found this problem by observing open stuck goroutines via pprof)
		for range watcher.ResultChan() {
			// Drain any pending events, so that channel is
			// closed and internal goroutine is discarded
		}
	}()

	for {
		select {
		case e, ok := <-watcher.ResultChan():
			if !ok || e.Object == nil {
				// Watcher may expire, hence try to retry
				return true, nil
			}

			app, ok := e.Object.(*v1alpha1.App)
			if !ok {
				continue
			}

			callback(*app)

		case <-cancelCh:
			return false, nil
		}
	}
}
