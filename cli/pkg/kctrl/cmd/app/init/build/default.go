// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"sigs.k8s.io/yaml"
)

const defaultAppBuildYAML = `
--- 
apiVersion: kctrl.carvel.dev/v1alpha1
kind: AppBuild
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
      serviceAccountName: 
      template:
      deploy:
        - kapp: {}
      fetch:
        - http: {}

  release:
  - resource: {}
`

func NewDefaultAppBuild() (AppBuild, error) {
	var appBuild AppBuild
	err := yaml.Unmarshal([]byte(defaultAppBuildYAML), &appBuild)
	if err != nil {
		return AppBuild{}, err
	}
	return appBuild, nil
}

func NewDefaultAppTemplateSpec() (*v1alpha12.AppTemplateSpec, error) {
	appBuild, err := NewDefaultAppBuild()
	if err != nil {
		return &v1alpha12.AppTemplateSpec{}, err
	}
	return appBuild.Spec.App, nil
}
