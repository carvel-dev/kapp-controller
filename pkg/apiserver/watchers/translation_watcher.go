// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package watchers

import "k8s.io/apimachinery/pkg/watch"

type translationFunc func(evt watch.Event) watch.Event

type TranslationWatcher struct {
	translate      translationFunc
	proxiedWatcher watch.Interface

	resultChan chan watch.Event
	stopCh     chan struct{}
}

var _ watch.Interface = &TranslationWatcher{}

func NewTranslationWatcher(translateFunc translationFunc, proxiedWatcher watch.Interface) *TranslationWatcher {
	tw := &TranslationWatcher{
		proxiedWatcher: proxiedWatcher,
		translate:      translateFunc,
		resultChan:     make(chan watch.Event, watch.DefaultChanSize),
		stopCh:         make(chan struct{}),
	}
	go tw.proxyEvents()
	return tw

}

func (tw *TranslationWatcher) ResultChan() <-chan watch.Event {
	return tw.resultChan
}

func (tw *TranslationWatcher) Stop() {
	tw.stopCh <- struct{}{}
	tw.proxiedWatcher.Stop()
}

func (tw *TranslationWatcher) proxyEvents() {
	proxiedChan := tw.proxiedWatcher.ResultChan()
	for {
		select {
		case evt := <-proxiedChan:
			tw.resultChan <- tw.translate(evt)
		case <-tw.stopCh:
			// maybe drain evt chan?
			for evt := range proxiedChan {
				tw.resultChan <- tw.translate(evt)
			}
			return
		}
	}
}
