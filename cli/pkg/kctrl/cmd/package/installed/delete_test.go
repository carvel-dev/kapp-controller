// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package installed

import (
	"context"
	"testing"

	cmdcore "carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/core"
	kcpkgv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const (
	testPkgiName = "test-pkgi"
	testPkgiNs   = "test-ns"
)

func Test_DeleteInstallCreatedResources(t *testing.T) {
	createdAnnotations := NewCreatedResourceAnnotations(testPkgiName, testPkgiNs)
	ownedAnnotations := map[string]string{KctrlPkgAnnotation: createdAnnotations.PackageAnnValue()}
	tanzuOwnedAnnotations := map[string]string{TanzuPkgAnnotation: createdAnnotations.PackageAnnValue()}

	testCases := []struct {
		desc              string
		pkgiAnnotations   map[string]string
		existing          []fakeDynamicObject
		expectedDeleted   []fakeDynamicKey
		expectedRemaining []fakeDynamicKey
	}{
		{
			desc: "deletes resources created by kctrl for the package install",
			pkgiAnnotations: map[string]string{
				KctrlPkgAnnotation + "-" + KindServiceAccount.AsString():     createdAnnotations.ServiceAccountAnnValue(),
				KctrlPkgAnnotation + "-" + KindClusterRole.AsString():        createdAnnotations.ClusterRoleAnnValue(),
				KctrlPkgAnnotation + "-" + KindClusterRoleBinding.AsString(): createdAnnotations.ClusterRoleBindingAnnValue(),
				KctrlPkgAnnotation + "-" + KindSecret.AsString():             createdAnnotations.SecretAnnValue(),
				TanzuPkgAnnotation + "-" + KindServiceAccount.AsString():     createdAnnotations.ServiceAccountAnnValue(),
				yttOverlayAnnotation: "test-pkgi-test-ns-overlays",
			},
			existing: []fakeDynamicObject{
				{key: serviceAccountKey(createdAnnotations.ServiceAccountAnnValue()), annotations: ownedAnnotations},
				{key: clusterRoleKey(createdAnnotations.ClusterRoleAnnValue()), annotations: ownedAnnotations},
				{key: clusterRoleBindingKey(createdAnnotations.ClusterRoleBindingAnnValue()), annotations: ownedAnnotations},
				{key: secretKey(createdAnnotations.SecretAnnValue()), annotations: ownedAnnotations},
				{key: secretKey("test-pkgi-test-ns-overlays"), annotations: ownedAnnotations},
			},
			expectedDeleted: []fakeDynamicKey{
				serviceAccountKey(createdAnnotations.ServiceAccountAnnValue()),
				clusterRoleKey(createdAnnotations.ClusterRoleAnnValue()),
				clusterRoleBindingKey(createdAnnotations.ClusterRoleBindingAnnValue()),
				secretKey(createdAnnotations.SecretAnnValue()),
				secretKey("test-pkgi-test-ns-overlays"),
			},
		},
		{
			desc: "deletes resources created by older versions of Tanzu CLI",
			pkgiAnnotations: map[string]string{
				TanzuPkgAnnotation + "-" + KindServiceAccount.AsString(): createdAnnotations.ServiceAccountAnnValue(),
				TanzuPkgAnnotation + "-" + KindClusterRole.AsString():    createdAnnotations.ClusterRoleAnnValue(),
			},
			existing: []fakeDynamicObject{
				{key: serviceAccountKey(createdAnnotations.ServiceAccountAnnValue()), annotations: tanzuOwnedAnnotations},
				{key: clusterRoleKey(createdAnnotations.ClusterRoleAnnValue()), annotations: tanzuOwnedAnnotations},
			},
			expectedDeleted: []fakeDynamicKey{
				serviceAccountKey(createdAnnotations.ServiceAccountAnnValue()),
				clusterRoleKey(createdAnnotations.ClusterRoleAnnValue()),
			},
		},
		{
			desc: "does not delete resources not created for the package install",
			pkgiAnnotations: map[string]string{
				KctrlPkgAnnotation + "-" + KindClusterRoleBinding.AsString(): "unrelated-crb",
				KctrlPkgAnnotation + "-" + KindSecret.AsString():             "unrelated-secret",
				yttOverlayAnnotation: "unrelated-overlays-secret",
			},
			existing: []fakeDynamicObject{
				{key: clusterRoleBindingKey("unrelated-crb")},
				{key: secretKey("unrelated-secret"), annotations: map[string]string{KctrlPkgAnnotation: "other-pkgi-test-ns"}},
				{key: secretKey("unrelated-overlays-secret")},
			},
			expectedRemaining: []fakeDynamicKey{
				clusterRoleBindingKey("unrelated-crb"),
				secretKey("unrelated-secret"),
				secretKey("unrelated-overlays-secret"),
			},
		},
		{
			desc: "ignores annotation keys not written by kctrl",
			pkgiAnnotations: map[string]string{
				"example.com/package-ClusterRoleBinding": createdAnnotations.ClusterRoleBindingAnnValue(),
			},
			existing: []fakeDynamicObject{
				{key: clusterRoleBindingKey(createdAnnotations.ClusterRoleBindingAnnValue()), annotations: ownedAnnotations},
			},
			expectedRemaining: []fakeDynamicKey{
				clusterRoleBindingKey(createdAnnotations.ClusterRoleBindingAnnValue()),
			},
		},
		{
			desc: "ignores resources that do not exist",
			pkgiAnnotations: map[string]string{
				KctrlPkgAnnotation + "-" + KindServiceAccount.AsString(): createdAnnotations.ServiceAccountAnnValue(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			dynamicClient := newFakeDynamicClient(tc.existing)
			noopUI := ui.NewNoopUI()
			o := &DeleteOptions{ui: noopUI, statusUI: cmdcore.NewStatusLoggingUI(noopUI),
				Name: testPkgiName, NamespaceFlags: cmdcore.NamespaceFlags{Name: testPkgiNs}}

			pkgi := &kcpkgv1alpha1.PackageInstall{ObjectMeta: metav1.ObjectMeta{
				Name: testPkgiName, Namespace: testPkgiNs, Annotations: tc.pkgiAnnotations}}

			err := o.deleteInstallCreatedResources(pkgi, dynamicClient)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedDeleted, dynamicClient.deleted)
			for _, key := range tc.expectedRemaining {
				assert.Contains(t, dynamicClient.objects, key, "expected %v not to be deleted", key)
			}
		})
	}
}

type fakeDynamicKey struct {
	resource  string
	namespace string
	name      string
}

type fakeDynamicObject struct {
	key         fakeDynamicKey
	annotations map[string]string
}

func serviceAccountKey(name string) fakeDynamicKey {
	return fakeDynamicKey{KindServiceAccount.Resource(), testPkgiNs, name}
}

func secretKey(name string) fakeDynamicKey {
	return fakeDynamicKey{KindSecret.Resource(), testPkgiNs, name}
}

func clusterRoleKey(name string) fakeDynamicKey {
	return fakeDynamicKey{KindClusterRole.Resource(), "", name}
}

func clusterRoleBindingKey(name string) fakeDynamicKey {
	return fakeDynamicKey{KindClusterRoleBinding.Resource(), "", name}
}

// fakeDynamicClient is a minimal in-memory dynamic client supporting Get and Delete
type fakeDynamicClient struct {
	objects map[fakeDynamicKey]*unstructured.Unstructured
	deleted []fakeDynamicKey
}

var _ dynamic.Interface = &fakeDynamicClient{}

func newFakeDynamicClient(objs []fakeDynamicObject) *fakeDynamicClient {
	c := &fakeDynamicClient{objects: map[fakeDynamicKey]*unstructured.Unstructured{}}
	for _, obj := range objs {
		u := &unstructured.Unstructured{Object: map[string]interface{}{}}
		u.SetName(obj.key.name)
		u.SetNamespace(obj.key.namespace)
		u.SetAnnotations(obj.annotations)
		c.objects[obj.key] = u
	}
	return c
}

func (c *fakeDynamicClient) Resource(gvr schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	return &fakeDynamicResource{client: c, gvr: gvr}
}

type fakeDynamicResource struct {
	dynamic.NamespaceableResourceInterface
	client    *fakeDynamicClient
	gvr       schema.GroupVersionResource
	namespace string
}

func (r *fakeDynamicResource) Namespace(ns string) dynamic.ResourceInterface {
	return &fakeDynamicResource{client: r.client, gvr: r.gvr, namespace: ns}
}

func (r *fakeDynamicResource) key(name string) fakeDynamicKey {
	return fakeDynamicKey{resource: r.gvr.Resource, namespace: r.namespace, name: name}
}

func (r *fakeDynamicResource) Get(_ context.Context, name string, _ metav1.GetOptions, _ ...string) (*unstructured.Unstructured, error) {
	obj, found := r.client.objects[r.key(name)]
	if !found {
		return nil, errors.NewNotFound(r.gvr.GroupResource(), name)
	}
	return obj, nil
}

func (r *fakeDynamicResource) Delete(_ context.Context, name string, _ metav1.DeleteOptions, _ ...string) error {
	if _, found := r.client.objects[r.key(name)]; !found {
		return errors.NewNotFound(r.gvr.GroupResource(), name)
	}
	delete(r.client.objects, r.key(name))
	r.client.deleted = append(r.client.deleted, r.key(name))
	return nil
}
