// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kappctrl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kappipkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	versions "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/fakes"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

var _ = Describe("Update Package", func() {
	var (
		ctl     *pkgClient
		crtCtl  *fakes.CRTClusterClient
		kappCtl *fakes.KappClient
		err     error
		opts    = tkgpackagedatamodel.PackageOptions{
			PkgInstallName:  testPkgInstallName,
			Namespace:       testNamespaceName,
			Version:         "2.0.0",
			PollInterval:    testPollInterval,
			PollTimeout:     testPollTimeout,
			CreateNamespace: true,
			Install:         false,
		}
		options  = opts
		progress *tkgpackagedatamodel.PackageProgress
	)

	JustBeforeEach(func() {
		progress = &tkgpackagedatamodel.PackageProgress{
			ProgressMsg: make(chan string, 10),
			Err:         make(chan error),
			Done:        make(chan struct{}),
		}
		ctl = &pkgClient{kappClient: kappCtl}
		go ctl.UpdatePackage(&options, progress)
		err = testReceive(progress)
	})

	Context("failure in getting the installed package due to GetPackageInstall API error", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, nil, errors.New("failure in GetPackageInstall"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageInstall"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure due to PackageInstall being nil and --install flag not provided", func() {
		BeforeEach(func() {
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, nil, apierrors.NewNotFound(schema.GroupResource{Resource: "PackageInstall"}, testPkgInstallName))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("package install does not exist in the namespace"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in installing the package due to a failure in ListPackages", func() {
		BeforeEach(func() {
			options.Install = true
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, nil, nil)
			kappCtl.ListPackagesReturns(nil, errors.New("failure in ListPackages"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to list package versions"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in installing the not-already-existing package due to GetPackageByName API error", func() {
		BeforeEach(func() {
			options.Install = true
			options.PackageName = testPkgName
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			kappCtl.GetPackageInstallReturns(nil, nil)
			kappCtl.GetPackageMetadataByNameReturns(nil, errors.New("failure in GetPackageByName"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageByName"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the installed package due to GetPackageByName API error", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturns(testPkgInstall, nil)
			kappCtl.GetPackageMetadataByNameReturns(nil, errors.New("failure in GetPackageByName"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageByName"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the installed package with nil version spec", func() {
		BeforeEach(func() {
			options.Install = true
			options.PackageName = testPkgName
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			testPkgInstall.Spec.PackageRef.VersionSelection = nil
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to update package 'test-pkg' as no existing package reference/version was found in the package install"))
		})
		AfterEach(func() {
			options = opts
			testPkgInstall.Spec.PackageRef.VersionSelection = testVersionSelection
		})
	})

	Context("failure in updating the installed package with nil version spec", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			kappCtl = &fakes.KappClient{}
			testPkgInstall.Spec.PackageRef.VersionSelection = nil
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to update package 'test-pkg'"))
		})
		AfterEach(func() {
			options = opts
			testPkgInstall.Spec.PackageRef.VersionSelection = testVersionSelection
		})
	})

	Context("failure in updating the installed package due to failure in opening the provided secret value file", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			options.ValuesFile = testValuesFile
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to create secret based on values file: failed to read from data values file"))
		})
		AfterEach(func() {
			options = opts
			testPkgInstall.Spec.PackageRef.VersionSelection = testVersionSelection
		})
	})

	Context("failure in updating the installed package due to failure in creating the secret resource", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			options.ValuesFile = testValuesFile
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
			err = os.WriteFile(testValuesFile, []byte("test"), 0644)
			Expect(err).ToNot(HaveOccurred())
			secret = testSecret
			crtCtl.CreateReturns(errors.New("error on creating the secret resource"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error on creating the secret resource"))
		})
		AfterEach(func() {
			options = opts
			secret = &corev1.Secret{}
			testPkgInstall.Spec.PackageRef.VersionSelection = testVersionSelection
			err = os.Remove(testValuesFile)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("failure in updating the installed package due to failure in updating the secret resource", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			options.ValuesFile = testValuesFile
			testPkgInstall.Annotations = make(map[string]string)
			testPkgInstall.Annotations[tkgpackagedatamodel.TanzuPkgPluginAnnotation+"-Secret"] = fmt.Sprintf(tkgpackagedatamodel.SecretName, options.PkgInstallName, options.Namespace)
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
			err = os.WriteFile(testValuesFile, []byte("test"), 0644)
			Expect(err).ToNot(HaveOccurred())
			secret = testSecret
			crtCtl.UpdateReturns(errors.New("error on updating the secret resource"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error on updating the secret resource"))
		})
		AfterEach(func() {
			options = opts
			secret = &corev1.Secret{}
			testPkgInstall.Spec.PackageRef.VersionSelection = testVersionSelection
			err = os.Remove(testValuesFile)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("failure in updating the installed package due to UpdatePackageInstall API error", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
			testPkgInstall.Spec.PackageRef = &kappipkg.PackageRef{
				RefName:          testPkgInstallName,
				VersionSelection: &versions.VersionSelectionSemver{},
			}
			kappCtl.UpdatePackageInstallReturns(errors.New("failure in UpdatePackageInstall"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in UpdatePackageInstall"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the installed package due to GetPackageInstall API error in waitForResourceInstallation", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			options.Wait = true
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
			testPkgInstall.Spec.PackageRef = &kappipkg.PackageRef{
				RefName:          testPkgInstallName,
				VersionSelection: &versions.VersionSelectionSemver{},
			}
			kappCtl.GetPackageInstallReturnsOnCall(1, nil, errors.New("failure in GetPackageInstall"))
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failure in GetPackageInstall"))
		})
		AfterEach(func() { options = opts })
	})

	Context("failure in updating the installed package due to reconciliation failure", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			options.Wait = true
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
			testPkgInstall.Spec.PackageRef = &kappipkg.PackageRef{
				RefName:          testPkgInstallName,
				VersionSelection: &versions.VersionSelectionSemver{},
			}
			kappCtl.GetPackageInstallReturnsOnCall(1, testPkgInstall, nil)
			Expect(testPkgInstall.Status.ObservedGeneration).To(Equal(testPkgInstall.Generation))
			Expect(len(testPkgInstall.Status.Conditions)).To(BeNumerically("==", 2))
			testPkgInstall.Status.Conditions[1].Type = kappctrl.ReconcileFailed
			testPkgInstall.Status.UsefulErrorMessage = testUsefulErrMsg
		})
		It(testFailureMsg, func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("resource reconciliation failed: %s", testUsefulErrMsg)))
		})
		AfterEach(func() {
			options = opts
			testPkgInstall.Status.Conditions[1].Type = kappctrl.ReconcileSucceeded
		})
	})

	Context("success in installing the not-already-existing package", func() {
		BeforeEach(func() {
			options.Install = true
			options.PackageName = testPkgName
			options.Version = testPkgVersion
			kappCtl = &fakes.KappClient{}
			crtCtl = &fakes.CRTClusterClient{}
			kappCtl.GetClientReturns(crtCtl)
			kappCtl.GetPackageInstallReturnsOnCall(0, nil, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})

	Context("success in updating the installed package", func() {
		BeforeEach(func() {
			options.Version = testPkgVersion
			kappCtl = &fakes.KappClient{}
			kappCtl.GetPackageInstallReturnsOnCall(0, testPkgInstall, nil)
			kappCtl.ListPackagesReturns(testPkgVersionList, nil)
			testPkgInstall.Spec.PackageRef = &kappipkg.PackageRef{
				RefName:          testPkgInstallName,
				VersionSelection: &versions.VersionSelectionSemver{},
			}
		})
		It(testSuccessMsg, func() {
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() { options = opts })
	})
})
