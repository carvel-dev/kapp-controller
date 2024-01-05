// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package fetch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcconfig "github.com/vmware-tanzu/carvel-kapp-controller/pkg/config"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/fetch"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func Test_AddDir_skipsTLS(t *testing.T) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"dangerousSkipTLSVerify": "always.trustworthy.com, selectively.trusted.net:123456",
		},
	}
	k8scs := k8sfake.NewSimpleClientset(configMap)
	config, err := kcconfig.NewConfig(k8scs)
	assert.NoError(t, err)

	vendir := fetch.NewVendir("default", k8scs,
		fetch.VendirOpts{SkipTLSConfig: config}, exec.NewPlainCmdRunner())

	type testCase struct {
		URL           string
		shouldSkipTLS bool
	}
	testCases := []testCase{
		{"always.trustworthy.com/myrepo/myimage:tag", true},
		{"never.trustworthy.com/myrepo/myimage:tag", false},
		{"selectively.trusted.net:123456/myrepo/myimage:tag", true},
		{"selectively.trusted.net:7777/myrepo/myimage:tag", false},
	}
	for i, tc := range testCases {
		err = vendir.AddDir(v1alpha1.AppFetch{
			Image: &v1alpha1.AppFetchImage{
				URL: tc.URL,
			},
		},
			"dirpath/0")
		assert.NoError(t, err)

		vConf := vendir.Config()
		assert.Equal(t, i+1, len(vConf.Directories), "Failed on iteration %d", i)
		assert.Equal(t, tc.shouldSkipTLS, vConf.Directories[i].Contents[0].Image.DangerousSkipTLSVerify, "Failed with URL %s", tc.URL)
	}
}

func TestExtractHost(t *testing.T) {
	tests := []struct {
		name     string
		isGitURL bool
		want     string
	}{
		{
			name:     "ubuntu:latest",
			isGitURL: false,
			want:     "index.docker.io",
		},
		{
			name:     "foo/bar:v1.2.3",
			isGitURL: false,
			want:     "index.docker.io",
		},
		{
			name:     "ghcr.io/foo/bar:foo",
			isGitURL: false,
			want:     "ghcr.io",
		},
		{
			name:     "foo.domain:5426/foo/bar@sha256:blah",
			isGitURL: false,
			want:     "foo.domain:5426",
		},
		{
			name:     "https://github.com/bitnami/charts/",
			isGitURL: true,
			want:     "github.com",
		},
		{
			name:     "http://github.com/bitnami/charts/",
			isGitURL: true,
			want:     "github.com",
		},
		{
			name:     "ssh://username@hostname.com:/path/to/repo.git",
			isGitURL: true,
			want:     "hostname.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fetch.ExtractHost(tt.name, tt.isGitURL); got != tt.want {
				t.Errorf("ExtractHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
