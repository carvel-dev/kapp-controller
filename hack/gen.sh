#!/bin/bash

set -e

./vendor/k8s.io/code-generator/generate-groups.sh \
	all github.com/k14s/kapp-controller/pkg/client github.com/k14s/kapp-controller/pkg/apis kappctrl:v1alpha1 \
	--go-header-file ./hack/gen-boilerplate.txt
