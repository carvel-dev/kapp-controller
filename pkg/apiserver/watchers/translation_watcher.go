// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package watchers

import (
	"k8s.io/apimachinery/pkg/watch"
)

// Translates the object of the original event to a new object
// before sending
type translationFunc func(evt watch.Event) watch.Event

// Determines if an event should be sent
type filterFunc func(evt watch.Event) bool

type TranslationWatcher struct {
	translate      translationFunc
	filter         filterFunc
	proxiedWatcher watch.Interface

	resultChan chan watch.Event
	stopCh     chan struct{}
}

var _ watch.Interface = &TranslationWatcher{}

func NewTranslationWatcher(translateFunc translationFunc, filterFunc filterFunc, proxiedWatcher watch.Interface) *TranslationWatcher {
	tw := &TranslationWatcher{
		proxiedWatcher: proxiedWatcher,
		translate:      translateFunc,
		filter:         filterFunc,
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
	close(tw.stopCh)
	tw.proxiedWatcher.Stop()
}

// Not sure why we get empty events from the CRD watch instance,
// but since they are empty, filter them out here, as including
// them leads to some error messages
func (tw *TranslationWatcher) proxyEvents() {
	proxiedChan := tw.proxiedWatcher.ResultChan()
	for {
		select {
		case evt, ok := <-proxiedChan:
			if !ok {
				close(tw.resultChan)
				return
			}

			trasnlatedEvt := tw.translate(evt)
			if tw.filter(trasnlatedEvt) {
				tw.resultChan <- tw.translate(evt)
			}
		case <-tw.stopCh:
			close(tw.resultChan)
			return
		}
	}
}
