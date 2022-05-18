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
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

func NewConfigFromFiles(paths []string) (kcv1alpha1.App, []corev1.Secret, []corev1.ConfigMap, error) {
	var apps []kcv1alpha1.App
	var secrets []corev1.Secret
	var configMaps []corev1.ConfigMap

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
			secrets = append(secrets, secret)

		case res.APIVersion == "v1" && res.Kind == "ConfigMap":
			var cm corev1.ConfigMap
			err := yaml.Unmarshal(docBytes, &cm)
			if err != nil {
				return fmt.Errorf("Unmarshaling config map: %s", err)
			}
			configMaps = append(configMaps, cm)

		case res.APIVersion == "kappctrl.k14s.io/v1alpha1" && res.Kind == "App":
			var app kcv1alpha1.App
			err := yaml.Unmarshal(docBytes, &app)
			if err != nil {
				return fmt.Errorf("Unmarshaling App: %s", err)
			}
			apps = append(apps, app)

		default:
			return fmt.Errorf("Unknown apiVersion '%s' or kind '%s' for resource",
				res.APIVersion, res.Kind)
		}
		return nil
	})
	if err != nil {
		return kcv1alpha1.App{}, nil, nil, err
	}

	if len(apps) == 0 {
		return kcv1alpha1.App{}, nil, nil, fmt.Errorf("Expected to find at least one App, but found none")
	}
	if len(apps) > 1 {
		return kcv1alpha1.App{}, nil, nil, fmt.Errorf("Expected to find exactly one App, but found multiple")
	}

	return apps[0], secrets, configMaps, nil
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
