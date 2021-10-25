// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("Get Repository", func() {
	var (
		ctl     *pkgClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.RepositoryOptions{
			RepositoryName: testRepoName,
			Namespace:      testNamespaceName,
		}
		options    = opts
		repository *kappipkg.PackageRepository
	)

	JustBeforeEach(func() {
		ctl = &pkgClient{kappClient: kappCtl}
		repository, err = ctl.GetRepository(&options)
	})

	Context("failure in getting package repository due to GetPackageRepository API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(nil, errors.New("failure in GetPackageRepository"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageRepository"))
			Expect(repository).To(BeNil())
		})
		AfterEach(func() { options = opts })
	})

	Context("success in getting package repository", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(testRepository, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(repository).NotTo(BeNil())
			Expect(repository).To(Equal(testRepository))
		})
		AfterEach(func() { options = opts })
	})
})
