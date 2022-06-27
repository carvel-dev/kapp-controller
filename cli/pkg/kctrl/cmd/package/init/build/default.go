// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"sigs.k8s.io/yaml"
)

const defaultPackageBuildYAML = `
--- 
apiVersion: kctrl.carvel.dev/v1alpha1
kind: PackageBuild
spec: 
  template:
    spec:
      export:
      - imgpkgBundle:
          image:
          useKbldImagesLock: true
          includePaths:
          - kubernetes-manifests/**/*
          - config-dev/build.yml
      app:
        spec:
          template:
          deploy:
          - kapp: {}
          fetch:
          - http: {}

  release:
  - resource: {}
`

func NewDefaultPackageBuild() (*PackageBuild, error) {
	var packageBuild PackageBuild
	err := yaml.Unmarshal([]byte(defaultPackageBuildYAML), &packageBuild)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal PackageBuild: %s", err.Error())
	}
	return &packageBuild, nil
}

//TODO Need to do MewDefaultAppTemplateSpec

const defaultPackageMetadataYAML = `
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: samplepackage.corp.com
  namespace: default
spec: 
  categories: 
  - demo
  displayName: ""
  longDescription: ""
  shortDescription: ""
`

func NewDefaultPackageMetadata() (v1alpha1.PackageMetadata, error) {
	var pkgMetadata v1alpha1.PackageMetadata
	err := yaml.Unmarshal([]byte(defaultPackageMetadataYAML), &pkgMetadata)
	if err != nil {
		return pkgMetadata, fmt.Errorf("unable to unmarshal PackageMetadata: %s", err.Error())
	}
	return pkgMetadata, nil
}

const defaultPackageYAML = `
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata: 
  name: samplepackage.corp.com.1.0.0
  namespace: default
spec: 
  licenses: 
  - "Apache 2.0"
  refName: samplepackage.corp.com
  releaseNotes: "Initial release"
  template: 
    spec: 
      deploy:
      - kapp: {}
      fetch: 
      - http: {}
      template:
      - kbld: {}
  version: "1.0.0"
`

const defaultPackageInstallYAML = `
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
`

func NewDefaultPackage() (v1alpha1.Package, error) {
	var pkg v1alpha1.Package
	err := yaml.Unmarshal([]byte(defaultPackageYAML), &pkg)
	if err != nil {
		return pkg, fmt.Errorf("unable to unmarshal Package: %s", err.Error())
	}
	return pkg, nil
}

func NewDefaultPackageInstall() (v1alpha12.PackageInstall, error) {
	var packageInstall v1alpha12.PackageInstall
	err := yaml.Unmarshal([]byte(defaultPackageInstallYAML), &packageInstall)
	if err != nil {
		return packageInstall, fmt.Errorf("unable to unmarshal PackageInstall: %s", err.Error())
	}
	return packageInstall, nil
}
