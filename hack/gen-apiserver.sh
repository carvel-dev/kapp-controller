#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

source hack/utils.sh
export GOPATH="$(go_mod_gopath_hack)"
trap "rm -rf ${GOPATH}; git checkout vendor" EXIT
KC_PKG="github.com/vmware-tanzu/carvel-kapp-controller"

# Following patch allows us to name gen-s with a name Package
# (without it generated Go code is not valid since word "package" is reserved)
git checkout vendor/k8s.io/gengo/namer/namer.go
git apply ./hack/gen-apiserver-namer.patch

rm -rf pkg/apiserver/{client,openapi}

echo "Generating clients and APIs"
go run vendor/k8s.io/code-generator/cmd/client-gen/main.go \
  --clientset-name versioned \
  --input-base "${KC_PKG}/pkg/apiserver/apis/" \
  --input "datapackaging/v1alpha1" \
  --output-package "${KC_PKG}/pkg/apiserver/client/clientset" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating listers"
go run vendor/k8s.io/code-generator/cmd/lister-gen/main.go \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --output-package "${KC_PKG}/pkg/apiserver/client/listers" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating informers"
go run vendor/k8s.io/code-generator/cmd/informer-gen/main.go \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --versioned-clientset-package "${KC_PKG}/pkg/apiserver/client/clientset/versioned" \
  --listers-package "${KC_PKG}/pkg/apiserver/client/listers" \
  --output-package "${KC_PKG}/pkg/apiserver/client/informers" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating deepcopy"
rm -f $(find pkg/apiserver|grep zz_generated.deepcopy)
go run vendor/k8s.io/code-generator/cmd/deepcopy-gen/main.go \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging" \
  -O zz_generated.deepcopy \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating conversions"
rm -f $(find pkg/apiserver|grep zz_generated.conversion)
go run vendor/k8s.io/code-generator/cmd/conversion-gen/main.go \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1,${KC_PKG}/pkg/apiserver/apis/datapackaging" \
  -O zz_generated.conversion \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating openapi"
rm -f $(find pkg/apiserver|grep zz_generated.openapi)
go run vendor/k8s.io/code-generator/cmd/openapi-gen/main.go \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --input-dirs "${KC_PKG}/pkg/apis/kappctrl/v1alpha1" \
  --input-dirs "github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr" \
  --input-dirs "k8s.io/api/core/v1" \
  --output-package "${KC_PKG}/pkg/apiserver/openapi" \
  -O zz_generated.openapi \
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

# TODO It seems this command messes around with protos in vendor directory
go-to-protobuf \
  --proto-import "${GOPATH}/src/${KC_PKG}/vendor" \
  --packages "-github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1,${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --vendor-output-base="${GOPATH}/src/${KC_PKG}/vendor" \
  --go-header-file hack/gen-boilerplate.txt

echo "GEN SUCCESS"
