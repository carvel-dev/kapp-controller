// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"context"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	fakediscovery "k8s.io/client-go/discovery/fake"
	fakekube "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	fakeaggregator "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/fake"
)

// fakeCAProvider implements dynamiccertificates.CAContentProvider
type fakeCAProvider struct {
	bundle []byte
}

func (f *fakeCAProvider) Name() string                               { return "fake-ca-provider" }
func (f *fakeCAProvider) CurrentCABundleContent() []byte             { return f.bundle }
func (f *fakeCAProvider) AddListener(_ dynamiccertificates.Listener) {}
func (f *fakeCAProvider) VerifyOptions() (x509.VerifyOptions, bool) {
	return x509.VerifyOptions{}, false
}

func Test_updateAPIService(t *testing.T) {
	logger := logr.Discard()

	tests := []struct {
		name           string
		existingBundle []byte
		newBundle      []byte
		expectUpdate   bool
	}{
		{
			name:           "updates APIService when CA bundle is different",
			existingBundle: []byte("old-dead-pod-cert"),
			newBundle:      []byte("new-active-pod-cert"),
			expectUpdate:   true,
		},
		{
			name:           "does nothing when CA bundle is identical",
			existingBundle: []byte("active-pod-cert"),
			newBundle:      []byte("active-pod-cert"),
			expectUpdate:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			apiSvc := &apiregv1.APIService{
				ObjectMeta: metav1.ObjectMeta{Name: apiServiceName},
				Spec: apiregv1.APIServiceSpec{
					CABundle: tc.existingBundle,
				},
			}
			fakeClient := fakeaggregator.NewSimpleClientset(apiSvc)

			fakeProvider := &fakeCAProvider{bundle: tc.newBundle}

			err := updateAPIService(context.TODO(), logger, fakeClient, fakeProvider)
			require.NoError(t, err)

			actions := fakeClient.Actions()
			require.GreaterOrEqual(t, len(actions), 1, "expected at least a GET action")
			require.Equal(t, "get", actions[0].GetVerb())

			var updateActionFound bool
			for _, action := range actions {
				if action.GetVerb() == "update" {
					updateActionFound = true
					updateAction, ok := action.(clienttesting.UpdateAction)
					if !ok {
						t.Fatalf("Expected UpdateAction, got %T", action)
					}
					updatedSvc := updateAction.GetObject().(*apiregv1.APIService)
					require.Equal(t, tc.newBundle, updatedSvc.Spec.CABundle)
				}
			}

			if tc.expectUpdate {
				require.True(t, updateActionFound, "expected an UPDATE action to be executed, but none was found")
			} else {
				require.False(t, updateActionFound, "expected NO UPDATE action, but one was executed")
			}
		})
	}
}

func Test_canUseAdmissionPolicyPlugin(t *testing.T) {
	const (
		policyResource  = "mutatingadmissionpolicies"
		bindingResource = "mutatingadmissionpolicybindings"
	)
	fullDiscovery := []*metav1.APIResourceList{{
		GroupVersion: admissionPolicyAPIGroupVersion,
		APIResources: []metav1.APIResource{
			{Name: policyResource, Namespaced: false, Kind: "MutatingAdmissionPolicy"},
			{Name: bindingResource, Namespaced: false, Kind: "MutatingAdmissionPolicyBinding"},
		},
	}}

	// ssarReactor builds a testing reactor that allows the given
	// (resource, verb) pairs and denies everything else. This lets each
	// test model a specific RBAC state without a real cluster.
	ssarReactor := func(allowed map[string]bool) clienttesting.ReactionFunc {
		return func(action clienttesting.Action) (bool, runtime.Object, error) {
			create, ok := action.(clienttesting.CreateAction)
			if !ok {
				return false, nil, nil
			}
			ssar, ok := create.GetObject().(*authorizationv1.SelfSubjectAccessReview)
			if !ok {
				return false, nil, nil
			}
			attrs := ssar.Spec.ResourceAttributes
			key := fmt.Sprintf("%s/%s", attrs.Resource, attrs.Verb)
			ssar.Status.Allowed = allowed[key]
			return true, ssar, nil
		}
	}

	tests := []struct {
		name       string
		resources  []*metav1.APIResourceList
		ssarAllow  map[string]bool
		ssarError  error
		wantResult bool
	}{
		{
			name:       "resource absent from discovery -> skip plugin",
			resources:  nil,
			wantResult: false,
		},
		{
			name: "policy present but binding missing -> skip plugin",
			resources: []*metav1.APIResourceList{{
				GroupVersion: admissionPolicyAPIGroupVersion,
				APIResources: []metav1.APIResource{
					{Name: policyResource, Namespaced: false, Kind: "MutatingAdmissionPolicy"},
				},
			}},
			wantResult: false,
		},
		{
			name:      "both resources present, list+watch allowed on both -> enable plugin",
			resources: fullDiscovery,
			ssarAllow: map[string]bool{
				policyResource + "/list":   true,
				policyResource + "/watch":  true,
				bindingResource + "/list":  true,
				bindingResource + "/watch": true,
			},
			wantResult: true,
		},
		{
			name:      "list allowed but watch denied on policy -> skip plugin",
			resources: fullDiscovery,
			ssarAllow: map[string]bool{
				policyResource + "/list":   true,
				policyResource + "/watch":  false,
				bindingResource + "/list":  true,
				bindingResource + "/watch": true,
			},
			wantResult: false,
		},
		{
			name:      "list+watch denied on policy -> skip plugin",
			resources: fullDiscovery,
			ssarAllow: map[string]bool{
				policyResource + "/list":   false,
				policyResource + "/watch":  false,
				bindingResource + "/list":  true,
				bindingResource + "/watch": true,
			},
			wantResult: false,
		},
		{
			name:      "policy allowed but binding denied -> skip plugin",
			resources: fullDiscovery,
			ssarAllow: map[string]bool{
				policyResource + "/list":   true,
				policyResource + "/watch":  true,
				bindingResource + "/list":  false,
				bindingResource + "/watch": false,
			},
			wantResult: false,
		},
		{
			name:       "SubjectAccessReview API error -> skip plugin conservatively",
			resources:  fullDiscovery,
			ssarError:  fmt.Errorf("simulated api error"),
			wantResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kubeClient := fakekube.NewSimpleClientset()
			kubeClient.Discovery().(*fakediscovery.FakeDiscovery).Resources = tc.resources

			if tc.ssarError != nil {
				kubeClient.PrependReactor("create", "selfsubjectaccessreviews", func(action clienttesting.Action) (bool, runtime.Object, error) {
					return true, nil, tc.ssarError
				})
			} else if tc.ssarAllow != nil {
				kubeClient.PrependReactor("create", "selfsubjectaccessreviews", ssarReactor(tc.ssarAllow))
			}

			got := canUseAdmissionPolicyPlugin(kubeClient.Discovery(), kubeClient, policyResource, bindingResource, logr.Discard())
			require.Equal(t, tc.wantResult, got)
		})
	}
}
