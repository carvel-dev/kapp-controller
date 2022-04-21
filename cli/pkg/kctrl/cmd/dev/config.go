// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	dpv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

type Configs struct {
	Apps        []kcv1alpha1.App
	PkgInstalls []pkgv1alpha1.PackageInstall
	Pkgs        []dpv1alpha1.Package
	Secrets     []corev1.Secret
	ConfigMaps  []corev1.ConfigMap
}

func (c *Configs) ApplyNamespace(ns string) {
	// Prefer namespace specified in the configuration
	for i, res := range c.Apps {
		if len(res.Namespace) == 0 {
			res.Namespace = ns
			c.Apps[i] = res
		}
	}
	for i, res := range c.PkgInstalls {
		if len(res.Namespace) == 0 {
			res.Namespace = ns
			c.PkgInstalls[i] = res
		}
	}
	for i, res := range c.Pkgs {
		if len(res.Namespace) == 0 {
			res.Namespace = ns
			c.Pkgs[i] = res
		}
	}
	for i, res := range c.Secrets {
		if len(res.Namespace) == 0 {
			res.Namespace = ns
			c.Secrets[i] = res
		}
	}
	for i, res := range c.ConfigMaps {
		if len(res.Namespace) == 0 {
			res.Namespace = ns
			c.ConfigMaps[i] = res
		}
	}
}

func (c *Configs) PkgsAsObjects() []runtime.Object {
	var result []runtime.Object
	for _, pkg := range c.Pkgs {
		pkg := pkg.DeepCopy()
		result = append(result, pkg)
	}
	return result
}

func NewConfigFromFiles(paths []string) (Configs, error) {
	var configs Configs

	err := parseResources(paths, func(docBytes []byte) error {
		var res resource

		err := yaml.Unmarshal(docBytes, &res)
		if err != nil {
			return fmt.Errorf("Unmarshaling doc: %s", err)
		}

		switch {
		case res.APIVersion == "v1" && res.Kind == "Secret":
			var secret corev1.Secret
			err := yaml.Unmarshal(docBytes, &secret)
			if err != nil {
				return fmt.Errorf("Unmarshaling secret: %s", err)
			}
			configs.Secrets = append(configs.Secrets, secret)

		case res.APIVersion == "v1" && res.Kind == "ConfigMap":
			var cm corev1.ConfigMap
			err := yaml.Unmarshal(docBytes, &cm)
			if err != nil {
				return fmt.Errorf("Unmarshaling config map: %s", err)
			}
			configs.ConfigMaps = append(configs.ConfigMaps, cm)

		case res.APIVersion == "kappctrl.k14s.io/v1alpha1" && res.Kind == "App":
			var app kcv1alpha1.App
			err := yaml.Unmarshal(docBytes, &app)
			if err != nil {
				return fmt.Errorf("Unmarshaling App: %s", err)
			}
			configs.Apps = append(configs.Apps, app)

		case res.APIVersion == "data.packaging.carvel.dev/v1alpha1" && res.Kind == "Package":
			var pkg dpv1alpha1.Package
			err := yaml.Unmarshal(docBytes, &pkg)
			if err != nil {
				return fmt.Errorf("Unmarshaling Package: %s", err)
			}
			configs.Pkgs = append(configs.Pkgs, pkg)

		case res.APIVersion == "data.packaging.carvel.dev/v1alpha1" && res.Kind == "PackageMetadata":
			// ignore

		case res.APIVersion == "packaging.carvel.dev/v1alpha1" && res.Kind == "PackageInstall":
			var pkgi pkgv1alpha1.PackageInstall
			err := yaml.Unmarshal(docBytes, &pkgi)
			if err != nil {
				return fmt.Errorf("Unmarshaling PackageInstall: %s", err)
			}
			configs.PkgInstalls = append(configs.PkgInstalls, pkgi)

		default:
			return fmt.Errorf("Unknown apiVersion '%s' or kind '%s' for resource",
				res.APIVersion, res.Kind)
		}
		return nil
	})
	if err != nil {
		return configs, err
	}

	if len(configs.Apps) == 0 && len(configs.PkgInstalls) == 0 {
		return configs, fmt.Errorf("Expected to find at least one App or PackageInstall, but found none")
	}
	if len(configs.Apps) > 1 || len(configs.PkgInstalls) > 1 {
		return configs, fmt.Errorf("Expected to find exactly one App or PackageInstall, but found multiple")
	}

	return configs, nil
}

type resource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

func parseResources(paths []string, resourceFunc func([]byte) error) error {
	for _, path := range paths {
		var bs []byte
		var err error

		if path == "-" {
			bs, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("Reading config from stdin: %s", err)
			}
		} else {
			bs, err = ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("Reading config '%s': %s", path, err)
			}
		}

		reader := kyaml.NewYAMLReader(bufio.NewReaderSize(bytes.NewReader(bs), 4096))

		for {
			docBytes, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("Parsing config '%s': %s", path, err)
			}
			err = resourceFunc(docBytes)
			if err != nil {
				return fmt.Errorf("Parsing resource config '%s': %s", path, err)
			}
		}
	}
	return nil
}
