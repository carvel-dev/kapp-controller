#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Note if you are not seeing generated code, most likely it's being placed into a different folder
# (e.g. Do you have GOPATH directory structure correctly named for this project?)

export GOPATH=$(cd ../../../../; pwd)
KC_PKG="github.com/vmware-tanzu/carvel-kapp-controller"

rm -rf pkg/client

# Based on vendor/k8s.io/code-generator/generate-groups.sh
# (Converted to "go runs" so that there is no dependency on installed binaries.)

echo "Generating deepcopy funcs"
rm -f $(find pkg/apis|grep zz_generated.deepcopy.go)
go run vendor/k8s.io/code-generator/cmd/deepcopy-gen/main.go \
	--input-dirs ${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apis/packaging/v1alpha1,${KC_PKG}/pkg/apis/internalpackaging/v1alpha1 \
	-O zz_generated.deepcopy \
	--bounding-dirs ${KC_PKG}/pkg/apis \
	--go-header-file ./hack/gen-boilerplate.txt

echo "Generating clientset"
go run vendor/k8s.io/code-generator/cmd/client-gen/main.go \
	--clientset-name versioned \
	--input-base '' \
	--input ${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apis/packaging/v1alpha1,${KC_PKG}/pkg/apis/internalpackaging/v1alpha1 \
	--output-package ${KC_PKG}/pkg/client/clientset \
	--go-header-file ./hack/gen-boilerplate.txt

echo "Generating listers"
go run vendor/k8s.io/code-generator/cmd/lister-gen/main.go \
	--input-dirs ${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apis/packaging/v1alpha1,${KC_PKG}/pkg/apis/internalpackaging/v1alpha1 \
	--output-package ${KC_PKG}/pkg/client/listers \
	--go-header-file ./hack/gen-boilerplate.txt

echo "Generating informers"
go run vendor/k8s.io/code-generator/cmd/informer-gen/main.go \
	--input-dirs ${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apis/packaging/v1alpha1,${KC_PKG}/pkg/apis/internalpackaging/v1alpha1 \
	--versioned-clientset-package ${KC_PKG}/pkg/client/clientset/versioned \
	--listers-package ${KC_PKG}/pkg/client/listers \
	--output-package ${KC_PKG}/pkg/client/informers \
	--go-header-file ./hack/gen-boilerplate.txt

echo GEN SUCCESS
