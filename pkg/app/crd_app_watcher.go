package app

import (
	"fmt"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcclient "github.com/k14s/kapp-controller/pkg/client/clientset/versioned"
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

	watcher, err := w.appClient.KappctrlV1alpha1().Apps(w.app.Namespace).Watch(listOpts)
	if err != nil {
		return false, fmt.Errorf("Creating app watcher: %s", err)
	}

	defer watcher.Stop()

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
