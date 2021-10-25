// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("Uninstall Package", func() {
	var (
		ctl     *pkgClient
		crtCtl  *fakes.CRTClusterClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.PackageOptions{
			PkgInstallName: testPkgInstallName,
			Namespace:      testNamespaceName,
			PollInterval:   testPollInterval,
			PollTimeout:    testPollTimeout,
		}
		options    = opts
		progress   *tkgpackagedatamodel.PackageProgress
		pkgInstall = kappipkg.PackageInstall{
			TypeMeta: metav1.TypeMeta{Kind: tkgpackagedatamodel.KindPackageInstall},
			ObjectMeta: metav1.ObjectMeta{
				Name:       testPkgInstallName,
				Namespace:  testNamespaceName,
				Generation: 1,
				Annotations: map[string]string{
					tkgpackagedatamodel.TanzuPkgPluginAnnotation + "-" + tkgpackagedatamodel.KindClusterRole:        "test-pkg-test-ns-cluster-role",
					tkgpackagedatamodel.TanzuPkgPluginAnnotation + "-" + tkgpackagedatamodel.KindClusterRoleBinding: "test-pkg-test-ns-cluster-rolebinding",
					tkgpackagedatamodel.TanzuPkgPluginAnnotation + "-" + tkgpackagedatamodel.KindServiceAccount:     "test-pkg-test-ns-sa",
					tkgpackagedatamodel.TanzuPkgPluginAnnotation + "-" + tkgpackagedatamodel.KindSecret:             "test-pkg-test-ns-values",
				}},
		}
	)

	JustBeforeEach(func() {
		progress = &tkgpackagedatamodel.PackageProgress{
			ProgressMsg: make(chan string, 10),
			Err:         make(chan error),
			Done:        make(chan struct{}),
		}
		ctl = &pkgClient{kappClient: kappCtl}
		go ctl.UninstallPackage(&options, progress)
		err = testReceive(progress)
	})

	Context("failure in getting installed packages due to GetPackageInstall API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturns(nil, errors.New("failure in GetPackageInstall"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageInstall"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in deleting installed package", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			kappCtl.GetPackageInstallReturns(&pkgInstall, nil)
			crtCtl.DeleteReturns(errors.New("failure in PackageInstall deletion"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in PackageInstall deletion"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in PackageInstall CR deletion", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			kappCtl.GetPackageInstallReturns(&pkgInstall, nil)
			crtCtl.DeleteReturns(nil)
			pkgInstall.Status = kappipkg.PackageInstallStatus{
				GenericStatus: kappctrl.GenericStatus{
					Conditions: []kappctrl.AppCondition{
						{Type: kappctrl.Deleting, Status: corev1.ConditionTrue},
						{Type: kappctrl.DeleteFailed, Status: corev1.ConditionTrue},
					},
					UsefulErrorMessage: testUsefulErrMsg,
					ObservedGeneration: 1,
				},
			}
			Expect(pkgInstall.Status.ObservedGeneration).To(Equal(pkgInstall.Generation))
			kappCtl.GetPackageInstallReturns(&pkgInstall, nil)
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(testUsefulErrMsg))
		})
		AfterEach(func() { options = opts })
	})
})
