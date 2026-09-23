// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"context"
	"crypto/x509"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericfeatures "k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	clienttesting "k8s.io/client-go/testing"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	fakeaggregator "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/fake"
)

// The packaging apiserver must run with the WatchList feature disabled so that
// watch&resourceVersion=0 requests are not defaulted with sendInitialEvents,
// which host clusters that ship WatchList off (e.g. K8s 1.33) reject with 422.
func Test_disableWatchListFeatureGate(t *testing.T) {
	err := disableWatchListFeatureGate()
	require.NoError(t, err)

	require.False(t, utilfeature.DefaultFeatureGate.Enabled(genericfeatures.WatchList),
		"WatchList feature gate must be disabled for the packaging apiserver")
}

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
