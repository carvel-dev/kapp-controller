// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("Update Secret", func() {
	type invalidDockerCfgJSON struct {
		Auths map[string]dockerConfigEntry `json:"invalid-key" datapolicy:"token"`
	}
	var (
		ctl              *pkgClient
		crtCtl           *fakes.CRTClusterClient
		kappCtl          *fakes.KappClient
		err              error
		dockerCfgContent []byte
		f                = false
		t                = true
		opts             = tkgpackagedatamodel.RegistrySecretOptions{
			Export:     tkgpackagedatamodel.TypeBoolPtr{ExportToAllNamespaces: &t},
			Namespace:  testNamespaceName,
			Password:   testPassword,
			SecretName: testSecretName,
			Username:   testUsername,
		}
		options                         = opts
		testDockerCfgInvalid            = invalidDockerCfgJSON{Auths: map[string]dockerConfigEntry{"us-east4-docker.pkg.dev": {Username: "test_user", Password: "test_password"}}}
		testDockerCfgMultipleRegistries = DockerConfigJSON{Auths: map[string]dockerConfigEntry{"us-east4-docker.pkg.dev": {Username: "test_user", Password: "test_password"}, "us-west-docker.pkg.dev": {Username: "test_user", Password: "test_password"}}}
		testDockerCfgNoUsername         = DockerConfigJSON{Auths: map[string]dockerConfigEntry{"us-east4-docker.pkg.dev": {Password: "test_password"}}}
		testDockerCfgNoPassword         = DockerConfigJSON{Auths: map[string]dockerConfigEntry{"us-east4-docker.pkg.dev": {Username: "test_user"}}}
	)

	JustBeforeEach(func() {
		ctl = &pkgClient{kappClient: kappCtl}
		err = ctl.UpdateRegistrySecret(&options)
	})

	Context("failure in updating Secret due to Secret Get NotFound error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.GetReturnsOnCall(0, apierrors.NewNotFound(schema.GroupResource{Resource: tkgpackagedatamodel.KindSecret}, testSecretName))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("secret '%s' does not exist in namespace '%s'", options.SecretName, options.Namespace)))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to Secret Get error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.GetReturnsOnCall(0, errors.New("failure in Secret Get"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in Secret Get"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to the non-existence of the 'auths' field in the secret", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerCfgInvalid)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("no 'auths' entry exists in secret '%s'", testSecretName)))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to the existence of multiple registries in the secret", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerCfgMultipleRegistries)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("updating secret '%s' is not allowed as multiple registry entries exists", testSecretName)))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to the non-existence of the 'username' field in the secret", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerCfgNoUsername)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("no 'username' entry exists in secret '%s'", testSecretName)))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to the non-existence of the 'password' field in the secret", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerCfgNoPassword)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("no 'password' entry exists in secret '%s'", testSecretName)))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to Secret Update error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, errors.New("failure in Secret Update"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in Secret Update"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to SecretExport Get error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.GetReturnsOnCall(1, errors.New("failure in SecretExport Get"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in SecretExport Get"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to SecretExport Create error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.GetReturnsOnCall(1, apierrors.NewNotFound(schema.GroupResource{Resource: tkgpackagedatamodel.KindSecretExport}, testSecretName))
			crtCtl.CreateReturnsOnCall(0, errors.New("failure in SecretExport Create"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in SecretExport Create"))
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating the Secret and Creating SecretExport", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.GetReturnsOnCall(1, apierrors.NewNotFound(schema.GroupResource{Resource: tkgpackagedatamodel.KindSecretExport}, testSecretName))
			crtCtl.CreateReturnsOnCall(0, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret due to SecretExport Update error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.GetReturnsOnCall(1, nil)
			crtCtl.UpdateReturnsOnCall(1, errors.New("failure in SecretExport Update"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in SecretExport Update"))
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating the Secret and updating SecretExport", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.GetReturnsOnCall(1, nil)
			crtCtl.UpdateReturnsOnCall(0, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the Secret with --export-to-all-namespaces=false due to SecretExport Delete error", func() {
		BeforeEach(func() {
			options.Export = tkgpackagedatamodel.TypeBoolPtr{ExportToAllNamespaces: &f}
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.DeleteReturnsOnCall(0, errors.New("failure in SecretExport Delete"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in SecretExport Delete"))
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating the Secret with --export-to-all-namespaces=false (SecretExport already not found)", func() {
		BeforeEach(func() {
			options.Export = tkgpackagedatamodel.TypeBoolPtr{ExportToAllNamespaces: &f}
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.DeleteReturnsOnCall(0, apierrors.NewNotFound(schema.GroupResource{Resource: tkgpackagedatamodel.KindSecretExport}, testSecretName))
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating the Secret and deleting SecretExport (--export-to-all-namespaces=false)", func() {
		BeforeEach(func() {
			options.Export = tkgpackagedatamodel.TypeBoolPtr{ExportToAllNamespaces: &f}
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			secret = testSecret
			crtCtl.GetReturnsOnCall(0, nil)
			dockerCfgContent, err = json.Marshal(testDockerConfig)
			Expect(err).NotTo(HaveOccurred())
			testSecret.Data[corev1.DockerConfigJsonKey] = dockerCfgContent
			crtCtl.UpdateReturnsOnCall(0, nil)
			crtCtl.DeleteReturnsOnCall(0, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})
})
