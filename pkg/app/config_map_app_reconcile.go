package app

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (a *ConfigMapApp) Reconcile() (reconcile.Result, error) {
	res, err := a.createOrUpdateStatusCfgMap()
	if err != nil {
		return res, err
	}

	err = a.app.Reconcile()
	if err != nil {
		return reconcile.Result{Requeue: true}, err // TODO move up
	}

	if a.appCfgMap.DeletionTimestamp != nil {
		return a.deleteCfgMaps()
	}

	return reconcile.Result{}, nil
}

func (a *ConfigMapApp) Delete() (reconcile.Result, error) {
	// TODO how to implement?
	return reconcile.Result{}, nil
}

func (r *ConfigMapApp) createOrUpdateStatusCfgMap() (reconcile.Result, error) {
	statusCfgMap, err := r.StatusConfigMap()
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	existingStatusCfgMap, err := r.coreClient.CoreV1().ConfigMaps(statusCfgMap.Namespace).Get(statusCfgMap.Name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			r.log.Error(err, "Could not fetch status ConfigMap")
			return reconcile.Result{Requeue: true}, err
		}
	}

	if existingStatusCfgMap != nil && len(existingStatusCfgMap.Name) > 0 {
		existingStatusCfgMap.ObjectMeta.OwnerReferences = statusCfgMap.ObjectMeta.OwnerReferences
		existingStatusCfgMap.Data = statusCfgMap.Data

		_, err = r.coreClient.CoreV1().ConfigMaps(existingStatusCfgMap.Namespace).Update(existingStatusCfgMap)
		if err != nil {
			r.log.Error(err, "Could not update status ConfigMap")
			return reconcile.Result{Requeue: true}, err
		}

		r.log.Info("Updated status ConfigMap") // TODO
		return reconcile.Result{}, nil
	}

	_, err = r.coreClient.CoreV1().ConfigMaps(statusCfgMap.Namespace).Create(statusCfgMap)
	if err != nil {
		r.log.Error(err, "Could not create status ConfigMap")
		return reconcile.Result{Requeue: true}, err
	}

	r.log.Info("Created status ConfigMap") // TODO
	return reconcile.Result{}, nil
}

func (r *ConfigMapApp) deleteCfgMaps() (reconcile.Result, error) {
	// Expecting that finalizers would have been cleared by now...

	err := r.coreClient.CoreV1().ConfigMaps(r.appCfgMap.Namespace).Delete(r.appCfgMap.Name, &metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			r.log.Error(err, "Could not delete app ConfigMap")
			return reconcile.Result{Requeue: true}, err
		}
	}

	// Status config map is deleted via owner reference

	return reconcile.Result{}, nil
}
