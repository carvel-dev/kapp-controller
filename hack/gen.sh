#!/bin/bash

set -e

gen_groups_path=./vendor/k8s.io/code-generator/generate-groups.sh

chmod +x $gen_groups_path

$gen_groups_path \
	all github.com/k14s/kapp-controller/pkg/client github.com/k14s/kapp-controller/pkg/apis kappctrl:v1alpha1 \
	--go-header-file ./hack/gen-boilerplate.txt

chmod -x $gen_groups_path
