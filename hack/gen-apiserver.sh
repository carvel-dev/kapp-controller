#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

source hack/utils.sh
export GOPATH="$(go_mod_gopath_hack)"
trap "sudo rm -rf ${GOPATH}; git checkout vendor" EXIT
KC_PKG="carvel.dev/kapp-controller"

# Following patch allows us to name gen-s with a name Package
# (without it generated Go code is not valid since word "package" is reserved)
git checkout vendor/k8s.io/gengo/v2/namer/namer.go
git apply ./hack/gen-apiserver-namer.patch

rm -rf pkg/apiserver/{client,openapi}

echo "Generating clients and APIs"
go run vendor/k8s.io/code-generator/cmd/client-gen/main.go \
  --clientset-name versioned \
  --input-base "${KC_PKG}/pkg/apiserver/apis/" \
  --input "datapackaging/v1alpha1" \
  --output-dir "pkg/apiserver/client/clientset" \
  --output-pkg "${KC_PKG}/pkg/apiserver/client/clientset" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating listers"
go run vendor/k8s.io/code-generator/cmd/lister-gen/main.go \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --output-pkg "${KC_PKG}/pkg/apiserver/client/listers" \
  --output-dir "pkg/apiserver/client/listers" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating informers"
go run vendor/k8s.io/code-generator/cmd/informer-gen/main.go \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --versioned-clientset-package "${KC_PKG}/pkg/apiserver/client/clientset/versioned" \
  --listers-package "${KC_PKG}/pkg/apiserver/client/listers" \
  --output-pkg "${KC_PKG}/pkg/apiserver/client/informers" \
  --output-dir "pkg/apiserver/client/informers" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating deepcopy"
rm -f $(find pkg/apiserver|grep zz_generated.deepcopy)
go run vendor/k8s.io/code-generator/cmd/deepcopy-gen/main.go \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging" \
  --output-file zz_generated.deepcopy.go \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating conversions"
rm -f $(find pkg/apiserver|grep zz_generated.conversion)
go run vendor/k8s.io/code-generator/cmd/conversion-gen/main.go \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging" \
  --output-file zz_generated.conversion.go \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating openapi"
rm -f $(find pkg/apiserver|grep zz_generated.openapi)
go run vendor/k8s.io/kube-openapi/cmd/openapi-gen/openapi-gen.go \
  "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  "${KC_PKG}/pkg/apis/kappctrl/v1alpha1" \
  "carvel.dev/vendir/pkg/vendir/versions/v1alpha1" \
  "k8s.io/apimachinery/pkg/apis/meta/v1" \
  "k8s.io/apimachinery/pkg/runtime" \
  "k8s.io/apimachinery/pkg/util/intstr" \
  "k8s.io/api/core/v1" \
  --output-pkg "${KC_PKG}/pkg/apiserver/openapi" \
  --output-dir "pkg/apiserver/openapi" \
  --output-file zz_generated.openapi.go \
  --go-header-file hack/gen-boilerplate.txt

# Install protoc binary as directed by https://github.com/gogo/protobuf#installation
# (Chosen: https://github.com/protocolbuffers/protobuf/releases/download/v3.0.2/protoc-3.0.2-osx-x86_64.zip)
# unzip archive into ./tmp/protoc-dl directory
export PATH=$PWD/tmp/protoc-dl/bin/:$PATH
protoc --version

# Generate binaries called out by protoc binary
export GOBIN=$PWD/tmp/gen-apiserver-bin
rm -rf $GOBIN
go install \
  github.com/gogo/protobuf/protoc-gen-gogo \
  github.com/gogo/protobuf/protoc-gen-gofast \
  golang.org/x/tools/cmd/goimports \
  k8s.io/code-generator/cmd/go-to-protobuf
export PATH=$GOBIN:$PATH

rm -f $(find pkg|grep '\.proto')

go-to-protobuf \
  --proto-import "${GOPATH}/src/${KC_PKG}/vendor" \
  --packages "-carvel.dev/vendir/pkg/vendir/versions/v1alpha1,${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --apimachinery-packages "-k8s.io/apimachinery/pkg/runtime/schema,-k8s.io/apimachinery/pkg/runtime,-k8s.io/apimachinery/pkg/apis/meta/v1" \
  --go-header-file hack/gen-boilerplate.txt \
  --output-dir "${GOPATH}/src" 

echo "GEN SUCCESS"
