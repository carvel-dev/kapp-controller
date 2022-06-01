// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkgrepository

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	datapackagingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	ctltpl "github.com/vmware-tanzu/carvel-kapp-controller/pkg/template"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

func (a *App) template(dirPath string) exec.CmdRunResult {
	if len(a.app.Spec.Template) == 0 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected at least one template option"))
	}

	genericOpts := ctltpl.GenericOpts{Name: a.app.Name, Namespace: a.app.Namespace}

	// first templating pass is to read all the files in
	var result exec.CmdRunResult
	var template ctltpl.Template
	template = a.templateFactory.NewYtt(*(a.app.Spec.Template[0]).Ytt, genericOpts)
	result, _ = template.TemplateDir(dirPath)
	if result.Error != nil {
		return result
	}

	// second templating applies a bunch of overlays and filters, some of which could be migrated to go code
	template = a.templateFactory.NewYtt(*(a.app.Spec.Template[1]).Ytt, genericOpts)
	result = template.TemplateStream(strings.NewReader(result.Stdout), dirPath)
	if result.Error != nil {
		return result
	}

	// intermediate phase deserializes and reserializes all of the resources
	// which guarantees that they only include fields this version of kc knows about.
	streamForSomeReason := strings.NewReader(result.Stdout)
	buf := make([]byte, streamForSomeReason.Len())
	_, err := io.ReadFull(streamForSomeReason, buf)
	if err != nil {
		result.Error = err
		return result
	}
	resources, err := FilterResources(buf)
	if err != nil {
		result.Error = err
		return result
	}

	// third templating inserts a kapp rebase rule to allow noops on identical resources provided by multiple pkgrs
	template = a.templateFactory.NewYtt(*(a.app.Spec.Template[2]).Ytt, genericOpts)
	result = template.TemplateStream(strings.NewReader(resources), dirPath)

	return result
}

// FilterResources takes a multi-doc yaml byte array of the templated
// contents of a PKGR, and filters each resource by deserializing and re-serializing
func FilterResources(yamls []byte) (string, error) {
	sch := runtime.NewScheme()
	err := scheme.AddToScheme(sch)
	if err != nil {
		return "", err
	}
	err = datapackagingv1alpha1.AddToScheme(sch)
	if err != nil {
		return "", err
	}
	decoder := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	filteredYamls := []string{}
	docs, err := yamlDocs(yamls)
	if err != nil {
		return "", err
	}
	for _, resourceYAML := range docs {

		obj, gKV, err := decoder(resourceYAML, nil, nil)
		if err != nil {
			return "", err
		}
		kind := gKV.Kind
		if gKV.Group != "data.packaging.carvel.dev" {
			return "", fmt.Errorf("got unexpected group: %s", gKV.Group)
		}

		if gKV.Version != "v1alpha1" {
			return "", fmt.Errorf("got unexpected version: %s", gKV.Version)
		}

		switch kind {
		case "Package":
			p := obj.(*datapackagingv1alpha1.Package)
			yp := printers.YAMLPrinter{}
			buffy := new(bytes.Buffer)
			yp.PrintObj(p, buffy)
			filteredYamls = append(filteredYamls, buffy.String())
		case "PackageMetadata":
			p := obj.(*datapackagingv1alpha1.PackageMetadata)
			yp := printers.YAMLPrinter{}
			buffy := new(bytes.Buffer)
			yp.PrintObj(p, buffy)
			filteredYamls = append(filteredYamls, buffy.String())
		default:
			return "", fmt.Errorf("PKGR contained unexpected kind: %s", kind)
		}

	}
	return strings.Join(filteredYamls, "\n---\n"), nil
}

func yamlDocs(yamls []byte) ([][]byte, error) {
	var docs [][]byte

	fileBytes := yamls
	reader := kyaml.NewYAMLReader(bufio.NewReaderSize(bytes.NewReader(fileBytes), 4096))

	for {
		docBytes, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		docs = append(docs, docBytes)
	}

	return docs, nil
}
