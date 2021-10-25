// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("Update Repository", func() {
	var (
		ctl     *pkgClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.RepositoryOptions{
			RepositoryName:   testRepoName,
			RepositoryURL:    testRepoURL,
			CreateRepository: false,
		}
		options           = opts
		progress          *tkgpackagedatamodel.PackageProgress
		pkgRepositoryList = &kappipkg.PackageRepositoryList{
			Items: []kappipkg.PackageRepository{*testRepository},
		}
	)

	JustBeforeEach(func() {
		progress = &tkgpackagedatamodel.PackageProgress{
			ProgressMsg: make(chan string, 10),
			Err:         make(chan error),
			Done:        make(chan struct{}),
		}
		ctl = &pkgClient{kappClient: kappCtl}
		go ctl.UpdateRepository(&options, progress)
		err = testReceive(progress)
	})

	Context("failure in getting the package repository due to GetPackageRepository API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(nil, errors.New("failure in GetPackageRepository"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageRepository"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in adding package repository as a repository with the same OCI registry URL already exists", func() {
		BeforeEach(func() {
			options.CreateRepository = true
			options.RepositoryName = testSecondRepoName
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(nil, apierrors.NewNotFound(schema.GroupResource{Resource: "Repository"}, testSecondRepoName))
			kappCtl.ListPackageRepositoriesReturns(pkgRepositoryList, nil)
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("package repository URL '%s' already exists in namespace '%s'", options.RepositoryURL, options.Namespace)))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating package repository due to UpdatePackageRepository API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(testRepository, apierrors.NewNotFound(schema.GroupResource{Resource: "Repository"}, testRepoName))
			kappCtl.UpdatePackageRepositoryReturns(errors.New("failure in UpdatePackageRepository"))
			kappCtl.ListPackageRepositoriesReturns(pkgRepositoryList, nil)
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in UpdatePackageRepository"))
		})
		AfterEach(func() { options = opts })
	})

	Context("success in adding package repository", func() {
		BeforeEach(func() {
			options.CreateRepository = true
			options.RepositoryName = testSecondRepoName
			options.RepositoryURL = testSecondRepoURL
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(nil, apierrors.NewNotFound(schema.GroupResource{Resource: "Repository"}, testSecondRepoName))
			kappCtl.ListPackageRepositoriesReturns(pkgRepositoryList, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).ToNot(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating package repository with tag", func() {
		BeforeEach(func() {
			options.RepositoryURL = testSecondRepoURL
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(testRepository, apierrors.NewNotFound(schema.GroupResource{Resource: "Repository"}, testRepoName))
			kappCtl.UpdatePackageRepositoryReturns(nil)
			kappCtl.ListPackageRepositoriesReturns(pkgRepositoryList, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).ToNot(HaveOccurred())
			updateRepoCallCnt := kappCtl.UpdatePackageRepositoryCallCount()
			Expect(updateRepoCallCnt).To(BeNumerically("==", 1))
			pkgRepo := kappCtl.UpdatePackageRepositoryArgsForCall(0)
			Expect(pkgRepo.Name).Should(Equal(options.RepositoryName))
			Expect(pkgRepo.Spec.Fetch.ImgpkgBundle.Image).Should(Equal(options.RepositoryURL))
			Expect(pkgRepo.Spec.Fetch.ImgpkgBundle.TagSelection).Should(BeNil())
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating package repository without tag", func() {
		BeforeEach(func() {
			options.RepositoryURL = testThirdRepoURL
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageRepositoryReturns(testRepository, apierrors.NewNotFound(schema.GroupResource{Resource: "Repository"}, testRepoName))
			kappCtl.UpdatePackageRepositoryReturns(nil)
			kappCtl.ListPackageRepositoriesReturns(pkgRepositoryList, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).ToNot(HaveOccurred())
			updateRepoCallCnt := kappCtl.UpdatePackageRepositoryCallCount()
			Expect(updateRepoCallCnt).To(BeNumerically("==", 1))
			pkgRepo := kappCtl.UpdatePackageRepositoryArgsForCall(0)
			Expect(pkgRepo.Name).Should(Equal(options.RepositoryName))
			Expect(pkgRepo.Spec.Fetch.ImgpkgBundle.Image).Should(Equal(options.RepositoryURL))
			Expect(pkgRepo.Spec.Fetch.ImgpkgBundle.TagSelection).ShouldNot(Equal(nil))
			Expect(pkgRepo.Spec.Fetch.ImgpkgBundle.TagSelection.Semver.Constraints).Should(Equal(tkgpackagedatamodel.DefaultRepositoryImageTagConstraint))
		})
		AfterEach(func() { options = opts })
	})
})
