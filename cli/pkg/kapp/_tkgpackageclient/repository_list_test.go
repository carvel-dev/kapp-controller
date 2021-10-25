// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("List Repositories", func() {
	var (
		ctl     *pkgClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.RepositoryOptions{
			Namespace:     testNamespaceName,
			AllNamespaces: false,
		}
		options        = opts
		repositories   *kappipkg.PackageRepositoryList
		repositoryList = &kappipkg.PackageRepositoryList{
			TypeMeta: metav1.TypeMeta{Kind: "PackageRepositoryList"},
			Items:    []kappipkg.PackageRepository{*testRepository},
		}
	)

	JustBeforeEach(func() {
		ctl = &pkgClient{kappClient: kappCtl}
		repositories, err = ctl.ListRepositories(&options)
	})

	Context("failure in listing package repositories due to ListPackageRepositories API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.ListPackageRepositoriesReturns(nil, errors.New("failure in ListPackageRepositories"))
			ctl = &pkgClient{kappClient: kappCtl}
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in ListPackageRepositories"))
			Expect(repositories).To(BeNil())
		})
	})

	Context("success in listing package repositories", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.ListPackageRepositoriesReturns(repositoryList, nil)
			ctl = &pkgClient{kappClient: kappCtl}
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(repositories).NotTo(BeNil())
			Expect(repositories).To(Equal(repositoryList))
		})
	})
})
