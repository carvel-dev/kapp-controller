// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"context"

	authenticationv1api "k8s.io/api/authentication/v1"
	corev1api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	aplcorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	discovery "k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	internalv1alpha1 "k8s.io/client-go/kubernetes/typed/apiserverinternal/v1alpha1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	appsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	authenticationv1 "k8s.io/client-go/kubernetes/typed/authentication/v1"
	authenticationv1beta1 "k8s.io/client-go/kubernetes/typed/authentication/v1beta1"
	authorizationv1 "k8s.io/client-go/kubernetes/typed/authorization/v1"
	authorizationv1beta1 "k8s.io/client-go/kubernetes/typed/authorization/v1beta1"
	autoscalingv1 "k8s.io/client-go/kubernetes/typed/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta2"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	batchv1beta1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
	certificatesv1 "k8s.io/client-go/kubernetes/typed/certificates/v1"
	certificatesv1beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	coordinationv1beta1 "k8s.io/client-go/kubernetes/typed/coordination/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	discoveryv1 "k8s.io/client-go/kubernetes/typed/discovery/v1"
	discoveryv1beta1 "k8s.io/client-go/kubernetes/typed/discovery/v1beta1"
	eventsv1 "k8s.io/client-go/kubernetes/typed/events/v1"
	eventsv1beta1 "k8s.io/client-go/kubernetes/typed/events/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	flowcontrolv1alpha1 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1alpha1"
	flowcontrolv1beta1 "k8s.io/client-go/kubernetes/typed/flowcontrol/v1beta1"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	networkingv1beta1 "k8s.io/client-go/kubernetes/typed/networking/v1beta1"
	nodev1 "k8s.io/client-go/kubernetes/typed/node/v1"
	nodev1alpha1 "k8s.io/client-go/kubernetes/typed/node/v1alpha1"
	nodev1beta1 "k8s.io/client-go/kubernetes/typed/node/v1beta1"
	policyv1 "k8s.io/client-go/kubernetes/typed/policy/v1"
	policyv1beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	rbacv1alpha1 "k8s.io/client-go/kubernetes/typed/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	schedulingv1 "k8s.io/client-go/kubernetes/typed/scheduling/v1"
	schedulingv1alpha1 "k8s.io/client-go/kubernetes/typed/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/client-go/kubernetes/typed/scheduling/v1beta1"
	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	storagev1alpha1 "k8s.io/client-go/kubernetes/typed/storage/v1alpha1"
	storagev1beta1 "k8s.io/client-go/kubernetes/typed/storage/v1beta1"
	rest "k8s.io/client-go/rest"
)

type localSecrets struct {
	Secrets []corev1api.Secret
}

func (localSecrets) SecretWithData(sec corev1api.Secret) corev1api.Secret {
	sec = *sec.DeepCopy()
	if len(sec.StringData) > 0 {
		sec.Data = map[string][]byte{}
		for k, v := range sec.StringData {
			sec.Data[k] = []byte(v)
		}
		sec.StringData = nil
	}
	return sec
}

type MinCoreClient struct {
	client          kubernetes.Interface
	localSecrets    *localSecrets
	localConfigMaps []corev1api.ConfigMap
}

var _ kubernetes.Interface = &MinCoreClient{}

func (*MinCoreClient) Discovery() discovery.DiscoveryInterface { panic("Not implemented"); return nil }
func (*MinCoreClient) AdmissionregistrationV1() admissionregistrationv1.AdmissionregistrationV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) InternalV1alpha1() internalv1alpha1.InternalV1alpha1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AppsV1() appsv1.AppsV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) AppsV1beta1() appsv1beta1.AppsV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AppsV1beta2() appsv1beta2.AppsV1beta2Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AuthenticationV1() authenticationv1.AuthenticationV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AuthenticationV1beta1() authenticationv1beta1.AuthenticationV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AuthorizationV1() authorizationv1.AuthorizationV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AuthorizationV1beta1() authorizationv1beta1.AuthorizationV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AutoscalingV1() autoscalingv1.AutoscalingV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AutoscalingV2beta1() autoscalingv2beta1.AutoscalingV2beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) AutoscalingV2beta2() autoscalingv2beta2.AutoscalingV2beta2Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) BatchV1() batchv1.BatchV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) BatchV1beta1() batchv1beta1.BatchV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) CertificatesV1() certificatesv1.CertificatesV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) CertificatesV1beta1() certificatesv1beta1.CertificatesV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) CoordinationV1beta1() coordinationv1beta1.CoordinationV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) CoordinationV1() coordinationv1.CoordinationV1Interface {
	panic("Not implemented")
	return nil
}
func (c *MinCoreClient) CoreV1() corev1.CoreV1Interface {
	return &MinCoreV1Client{c.client.CoreV1(), c.localSecrets, c.localConfigMaps}
}
func (*MinCoreClient) DiscoveryV1() discoveryv1.DiscoveryV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) EventsV1() eventsv1.EventsV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) EventsV1beta1() eventsv1beta1.EventsV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) FlowcontrolV1alpha1() flowcontrolv1alpha1.FlowcontrolV1alpha1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) FlowcontrolV1beta1() flowcontrolv1beta1.FlowcontrolV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) NetworkingV1() networkingv1.NetworkingV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) NetworkingV1beta1() networkingv1beta1.NetworkingV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) NodeV1() nodev1.NodeV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) NodeV1beta1() nodev1beta1.NodeV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) PolicyV1() policyv1.PolicyV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) PolicyV1beta1() policyv1beta1.PolicyV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) RbacV1() rbacv1.RbacV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) RbacV1beta1() rbacv1beta1.RbacV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) RbacV1alpha1() rbacv1alpha1.RbacV1alpha1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) SchedulingV1alpha1() schedulingv1alpha1.SchedulingV1alpha1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) SchedulingV1beta1() schedulingv1beta1.SchedulingV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) SchedulingV1() schedulingv1.SchedulingV1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) StorageV1beta1() storagev1beta1.StorageV1beta1Interface {
	panic("Not implemented")
	return nil
}
func (*MinCoreClient) StorageV1() storagev1.StorageV1Interface { panic("Not implemented"); return nil }
func (*MinCoreClient) StorageV1alpha1() storagev1alpha1.StorageV1alpha1Interface {
	panic("Not implemented")
	return nil
}

type MinCoreV1Client struct {
	client          corev1.CoreV1Interface
	localSecrets    *localSecrets
	localConfigMaps []corev1api.ConfigMap
}

func (*MinCoreV1Client) RESTClient() rest.Interface { panic("Not implemented"); return nil }
func (*MinCoreV1Client) ComponentStatuses() corev1.ComponentStatusInterface {
	panic("Not implemented")
	return nil
}
func (c *MinCoreV1Client) ConfigMaps(namespace string) corev1.ConfigMapInterface {
	return &ConfigMaps{namespace, c.client.ConfigMaps(namespace), c.localConfigMaps}
}
func (*MinCoreV1Client) Endpoints(namespace string) corev1.EndpointsInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) Events(namespace string) corev1.EventInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) LimitRanges(namespace string) corev1.LimitRangeInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) Namespaces() corev1.NamespaceInterface { panic("Not implemented"); return nil }
func (*MinCoreV1Client) Nodes() corev1.NodeInterface           { panic("Not implemented"); return nil }
func (*MinCoreV1Client) PersistentVolumes() corev1.PersistentVolumeInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) PersistentVolumeClaims(namespace string) corev1.PersistentVolumeClaimInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) Pods(namespace string) corev1.PodInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) PodTemplates(namespace string) corev1.PodTemplateInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) ReplicationControllers(namespace string) corev1.ReplicationControllerInterface {
	panic("Not implemented")
	return nil
}
func (*MinCoreV1Client) ResourceQuotas(namespace string) corev1.ResourceQuotaInterface {
	panic("Not implemented")
	return nil
}
func (c *MinCoreV1Client) Secrets(namespace string) corev1.SecretInterface {
	return &Secrets{namespace, c.client.Secrets(namespace), c.localSecrets}
}
func (*MinCoreV1Client) Services(namespace string) corev1.ServiceInterface {
	panic("Not implemented")
	return nil
}
func (c *MinCoreV1Client) ServiceAccounts(namespace string) corev1.ServiceAccountInterface {
	return &ServiceAccounts{namespace, c.client.ServiceAccounts(namespace)}
}

type ConfigMaps struct {
	namespace       string
	client          corev1.ConfigMapInterface
	localConfigMaps []corev1api.ConfigMap
}

var _ corev1.ConfigMapInterface = &ConfigMaps{}

func (*ConfigMaps) Create(ctx context.Context, configMap *corev1api.ConfigMap, opts metav1.CreateOptions) (*corev1api.ConfigMap, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ConfigMaps) Update(ctx context.Context, configMap *corev1api.ConfigMap, opts metav1.UpdateOptions) (*corev1api.ConfigMap, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ConfigMaps) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	panic("Not implemented")
	return nil
}
func (*ConfigMaps) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	panic("Not implemented")
	return nil
}
func (cm *ConfigMaps) Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1api.ConfigMap, error) {
	for _, configMap := range cm.localConfigMaps {
		if configMap.Name == name && configMap.Namespace == cm.namespace {
			configMapCopy := configMap
			return &configMapCopy, nil
		}
	}
	return cm.client.Get(ctx, name, opts)
}
func (*ConfigMaps) List(ctx context.Context, opts metav1.ListOptions) (*corev1api.ConfigMapList, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ConfigMaps) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ConfigMaps) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *corev1api.ConfigMap, err error) {
	panic("Not implemented")
	return nil, nil
}
func (*ConfigMaps) Apply(ctx context.Context, configMap *aplcorev1.ConfigMapApplyConfiguration, opts metav1.ApplyOptions) (result *corev1api.ConfigMap, err error) {
	panic("Not implemented")
	return nil, nil
}

type Secrets struct {
	namespace    string
	client       corev1.SecretInterface
	localSecrets *localSecrets
}

var _ corev1.SecretInterface = &Secrets{}

func (sec *Secrets) Create(ctx context.Context, secret *corev1api.Secret, opts metav1.CreateOptions) (*corev1api.Secret, error) {
	// TODO ignore creation (kapp-controller will create secretgen friendly secret)
	sec.localSecrets.Secrets = append(sec.localSecrets.Secrets, *secret.DeepCopy())
	return secret.DeepCopy(), nil
}
func (*Secrets) Update(ctx context.Context, secret *corev1api.Secret, opts metav1.UpdateOptions) (*corev1api.Secret, error) {
	panic("Not implemented")
	return nil, nil
}
func (*Secrets) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	panic("Not implemented")
	return nil
}
func (*Secrets) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	panic("Not implemented")
	return nil
}
func (sec *Secrets) Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1api.Secret, error) {
	for _, secret := range sec.localSecrets.Secrets {
		if secret.Name == name && secret.Namespace == sec.namespace {
			secCopy := sec.localSecrets.SecretWithData(secret)
			return &secCopy, nil
		}
	}
	return sec.client.Get(ctx, name, opts)
}
func (*Secrets) List(ctx context.Context, opts metav1.ListOptions) (*corev1api.SecretList, error) {
	panic("Not implemented")
	return nil, nil
}
func (*Secrets) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("Not implemented")
	return nil, nil
}
func (*Secrets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *corev1api.Secret, err error) {
	panic("Not implemented")
	return nil, nil
}
func (*Secrets) Apply(ctx context.Context, secret *aplcorev1.SecretApplyConfiguration, opts metav1.ApplyOptions) (result *corev1api.Secret, err error) {
	panic("Not implemented")
	return nil, nil
}

type ServiceAccounts struct {
	namespace string
	client    corev1.ServiceAccountInterface
}

var _ corev1.ServiceAccountInterface = &ServiceAccounts{}

func (*ServiceAccounts) Create(ctx context.Context, serviceAccount *corev1api.ServiceAccount, opts metav1.CreateOptions) (*corev1api.ServiceAccount, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ServiceAccounts) Update(ctx context.Context, serviceAccount *corev1api.ServiceAccount, opts metav1.UpdateOptions) (*corev1api.ServiceAccount, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ServiceAccounts) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	panic("Not implemented")
	return nil
}
func (*ServiceAccounts) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	panic("Not implemented")
	return nil
}
func (sa *ServiceAccounts) Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1api.ServiceAccount, error) {
	return sa.client.Get(ctx, name, opts)
}
func (*ServiceAccounts) List(ctx context.Context, opts metav1.ListOptions) (*corev1api.ServiceAccountList, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ServiceAccounts) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("Not implemented")
	return nil, nil
}
func (*ServiceAccounts) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *corev1api.ServiceAccount, err error) {
	panic("Not implemented")
	return nil, nil
}
func (*ServiceAccounts) Apply(ctx context.Context, serviceAccount *aplcorev1.ServiceAccountApplyConfiguration, opts metav1.ApplyOptions) (result *corev1api.ServiceAccount, err error) {
	panic("Not implemented")
	return nil, nil
}
func (*ServiceAccounts) CreateToken(ctx context.Context, serviceAccountName string, tokenRequest *authenticationv1api.TokenRequest, opts metav1.CreateOptions) (*authenticationv1api.TokenRequest, error) {
	panic("Not implemented")
	return nil, nil
}
