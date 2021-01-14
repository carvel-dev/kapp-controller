#!/bin/bash

set -e

# Note if you are not seeing generated code, most likely it's being placed into a different folder
# (e.g. Do you have GOPATH directory structure correctly named for this project?)

gen_groups_path=./vendor/k8s.io/code-generator/generate-groups.sh

chmod +x $gen_groups_path

rm -rf pkg/client

$gen_groups_path \
	all github.com/vmware-tanzu/carvel-kapp-controller/pkg/client github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis kappctrl:v1alpha1 \
	--go-header-file ./hack/gen-boilerplate.txt

chmod -x $gen_groups_path
