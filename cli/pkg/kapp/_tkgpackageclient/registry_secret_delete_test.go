// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("Delete Secret", func() {
	var (
		ctl           *pkgClient
		crtCtl        *fakes.CRTClusterClient
		kappCtl       *fakes.KappClient
		err           error
		isSecretFound bool
		opts          = tkgpackagedatamodel.RegistrySecretOptions{
			ExportToAllNamespaces: false,
			Namespace:             testNamespaceName,
			Password:              testPassword,
			Server:                testRegistry,
			SecretName:            testSecretName,
			Username:              testUsername,
			SkipPrompt:            true,
		}
		options = opts
	)

	JustBeforeEach(func() {
		ctl = &pkgClient{kappClient: kappCtl}
		isSecretFound, err = ctl.DeleteRegistrySecret(&options)
	})

	Context("success when trying to delete SecretExport returns NotFound error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.DeleteReturnsOnCall(0, apierrors.NewNotFound(schema.GroupResource{Resource: tkgpackagedatamodel.KindSecretExport}, testSecretName))
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in deleting SecretExport due to Delete API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.DeleteReturnsOnCall(0, errors.New("failure in delete SecretExport"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in delete SecretExport"))
			Expect(isSecretFound).To(BeTrue())
		})
		AfterEach(func() { options = opts })
	})

	Context("success when trying to delete Secret returns NotFound error", func() {
		BeforeEach(func() {
			options.ExportToAllNamespaces = true
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.DeleteReturnsOnCall(0, nil)
			crtCtl.DeleteReturnsOnCall(1, apierrors.NewNotFound(schema.GroupResource{Resource: tkgpackagedatamodel.KindSecret}, testSecretName))
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(isSecretFound).To(BeFalse())
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in deleting Secret due to Delete API error", func() {
		BeforeEach(func() {
			options.ExportToAllNamespaces = true
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.DeleteReturnsOnCall(0, nil)
			crtCtl.DeleteReturnsOnCall(1, errors.New("failure in delete Secret"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in delete Secret"))
			Expect(isSecretFound).To(BeTrue())
		})
		AfterEach(func() { options = opts })
	})

	Context("success in deleting Secret and SecretExport", func() {
		BeforeEach(func() {
			options.ExportToAllNamespaces = true
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			crtCtl.DeleteReturnsOnCall(0, nil)
			crtCtl.DeleteReturnsOnCall(1, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(isSecretFound).To(BeTrue())
		})
		AfterEach(func() { options = opts })
	})
})
