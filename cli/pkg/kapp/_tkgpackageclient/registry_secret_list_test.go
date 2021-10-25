// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	secretgenctrl "github.com/vmware-tanzu/carvel-secretgen-controller/pkg/apis/secretgen2/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

const testSecretExportName = "test-secret"

var testDockerConfig = DockerConfigJSON{Auths: map[string]dockerConfigEntry{"us-east4-docker.pkg.dev": {Username: "test_user", Password: "test_password"}}}

var testSecret = &corev1.Secret{
	TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindSecret, APIVersion: corev1.SchemeGroupVersion.String()},
	ObjectMeta: metav1.ObjectMeta{Name: testSecretName, Namespace: testNamespaceName},
	Type:       corev1.SecretTypeDockerConfigJson,
	Data:       map[string][]byte{},
}

var testSecretExport = &secretgenctrl.SecretExport{
	TypeMeta:   metav1.TypeMeta{Kind: tkgpackagedatamodel.KindSecretExport, APIVersion: secretgenctrl.SchemeGroupVersion.String()},
	ObjectMeta: metav1.ObjectMeta{Name: testSecretExportName, Namespace: testNamespaceName},
	Spec:       secretgenctrl.SecretExportSpec{ToNamespaces: []string{"*"}},
}

var _ = Describe("List Secrets", func() {
	var (
		ctl     *pkgClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.RegistrySecretOptions{
			Namespace:     testNamespaceName,
			AllNamespaces: false,
		}
		options    = opts
		secrets    *corev1.SecretList
		secretList = &corev1.SecretList{
			TypeMeta: metav1.TypeMeta{Kind: "SecretList"},
			Items:    []corev1.Secret{*testSecret},
		}
	)

	JustBeforeEach(func() {
		ctl = &pkgClient{kappClient: kappCtl}
		dockerCfgContent, _ := json.Marshal(testDockerConfig)
		testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
		secrets, err = ctl.ListRegistrySecrets(&options)
	})

	Context("failure in listing secrets due to ListRegistrySecrets API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.ListRegistrySecretsReturns(nil, errors.New("failure in ListRegistrySecrets"))
			ctl = &pkgClient{kappClient: kappCtl}
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in ListRegistrySecrets"))
			Expect(secrets).To(BeNil())
		})
	})

	Context("success in listing secrets", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.ListRegistrySecretsReturns(secretList, nil)
			ctl = &pkgClient{kappClient: kappCtl}
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(secrets).NotTo(BeNil())
			Expect(secrets).To(Equal(secretList))
		})
	})
})

var _ = Describe("List Secret Exports", func() {
	var (
		ctl     *pkgClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.RegistrySecretOptions{
			Namespace:     testNamespaceName,
			AllNamespaces: false,
		}
		options          = opts
		secretExports    *secretgenctrl.SecretExportList
		secretExportList = &secretgenctrl.SecretExportList{
			TypeMeta: metav1.TypeMeta{Kind: "SecretList"},
			Items:    []secretgenctrl.SecretExport{*testSecretExport},
		}
	)

	JustBeforeEach(func() {
		ctl = &pkgClient{kappClient: kappCtl}
		secretExports, err = ctl.ListSecretExports(&options)
	})

	Context("failure in listing secret exports due to ListSecretExports API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.ListSecretExportsReturns(nil, errors.New("failure in ListSecretExports"))
			ctl = &pkgClient{kappClient: kappCtl}
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in ListSecretExports"))
			Expect(secretExports).To(BeNil())
		})
	})

	Context("success in listing secret exports", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.ListSecretExportsReturns(secretExportList, nil)
			ctl = &pkgClient{kappClient: kappCtl}
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(secretExports).NotTo(BeNil())
			Expect(secretExports).To(Equal(secretExportList))
		})
	})
})
