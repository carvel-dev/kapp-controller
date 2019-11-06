package app

import (
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/template"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	cfgMapAppAnn    = "kapp-ctrl.k14s.io"
	cfgMapStatusAnn = "kapp-ctrl.k14s.io/status"

	appCfgMapDataKey    = "spec"
	statusCfgMapDataKey = "status"
)

type ConfigMapApp struct {
	app       *App
	appCfgMap *corev1.ConfigMap
	nsName    types.NamespacedName // TODO fill in?

	log        logr.Logger
	coreClient kubernetes.Interface
}

func NewConfigMapApp(appCfgMap *corev1.ConfigMap, log logr.Logger,
	coreClient kubernetes.Interface, fetchFactory fetch.Factory,
	templateFactory template.Factory) (*ConfigMapApp, error) {

	// Not an app config map
	if _, found := appCfgMap.Annotations[cfgMapAppAnn]; !found {
		return nil, nil
	}
	// Status config map, so ignore them
	if _, found := appCfgMap.Annotations[cfgMapStatusAnn]; found {
		return nil, nil
	}

	if appCfgMap.Data == nil {
		appCfgMap.Data = map[string]string{}
	}

	var spec v1alpha1.AppSpec

	err := yaml.Unmarshal([]byte(appCfgMap.Data[appCfgMapDataKey]), &spec)
	if err != nil {
		return nil, err
	}

	appCR := v1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:              appCfgMap.Name,
			Namespace:         appCfgMap.Namespace,
			DeletionTimestamp: appCfgMap.DeletionTimestamp,
		},
		Spec: spec,
	}

	cmApp := &ConfigMapApp{appCfgMap: appCfgMap, log: log, coreClient: coreClient}

	cmApp.app = NewApp(appCR, AppHooks{
		BlockDeletion:   cmApp.blockDeletion,
		UnblockDeletion: cmApp.unblockDeletion,
		UpdateStatus:    cmApp.updateStatus,
	}, fetchFactory, templateFactory)

	return cmApp, nil
}

func NewConfigMapAppFromName(nsName types.NamespacedName, log logr.Logger,
	coreClient kubernetes.Interface) *ConfigMapApp {

	return &ConfigMapApp{nil, nil, nsName, log, coreClient}
}

func (a *ConfigMapApp) StatusConfigMap() (*corev1.ConfigMap, error) {
	appStatusBs, err := a.app.StatusAsYAMLBytes()
	if err != nil {
		return nil, err
	}

	statusCfgMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        a.app.Name() + "-status",
			Namespace:   a.app.Namespace(),
			Annotations: a.appCfgMap.Annotations,
			Labels:      a.appCfgMap.Labels,
			// Set owner references for cascading deletion
			OwnerReferences: []metav1.OwnerReference{
				// TODO calculate schema automatically
				*metav1.NewControllerRef(a.appCfgMap.GetObjectMeta(), schema.GroupVersionKind{
					Group:   "",
					Version: "v1",
					Kind:    "ConfigMap",
				}),
			},
		},
		Data: map[string]string{
			statusCfgMapDataKey: string(appStatusBs),
		},
	}

	// Mark as status config map
	delete(statusCfgMap.Annotations, cfgMapAppAnn)
	statusCfgMap.Annotations[cfgMapStatusAnn] = ""

	return statusCfgMap, nil
}

func (a *ConfigMapApp) blockDeletion() error {
	a.log.Info("Blocking deletion")
	return a.updateAppCfgMap(func(appCfgMap *corev1.ConfigMap) {
		if !containsString(appCfgMap.ObjectMeta.Finalizers, deleteFinalizerName) {
			appCfgMap.ObjectMeta.Finalizers = append(appCfgMap.ObjectMeta.Finalizers, deleteFinalizerName)
		}
	})
}

func (a *ConfigMapApp) unblockDeletion() error {
	a.log.Info("Unblocking deletion")
	return a.updateAppCfgMap(func(appCfgMap *corev1.ConfigMap) {
		appCfgMap.ObjectMeta.Finalizers = removeString(appCfgMap.ObjectMeta.Finalizers, deleteFinalizerName)
	})
}

func (a *ConfigMapApp) updateStatus() error {
	_, err := a.createOrUpdateStatusCfgMap()
	return err
}

func (a *ConfigMapApp) updateAppCfgMap(updateFunc func(*corev1.ConfigMap)) error {
	existingAppCfgMap, err := a.coreClient.CoreV1().ConfigMaps(a.appCfgMap.Namespace).Get(a.appCfgMap.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	updateFunc(existingAppCfgMap)

	_, err = a.coreClient.CoreV1().ConfigMaps(existingAppCfgMap.Namespace).Update(existingAppCfgMap)
	return err
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
