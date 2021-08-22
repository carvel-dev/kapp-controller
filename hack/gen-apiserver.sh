#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

export GOPATH=$(cd ../../../../; pwd)
KC_PKG="github.com/vmware-tanzu/carvel-kapp-controller"

# Following patch allows us to name gen-s with a name Package
# (without it generated Go code is not valid since word "package" is reserved)
git checkout vendor/k8s.io/gengo/namer/namer.go
patch -u vendor/k8s.io/gengo/namer/namer.go ./hack/gen-apiserver-namer.patch

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

# Undo naming change
git checkout vendor/k8s.io/gengo/namer/namer.go

# Below protogen is configured to work without GOPATH var set
unset GOPATH

# Install protoc binary as directed by https://github.com/gogo/protobuf#installation
# (Chosen: https://github.com/protocolbuffers/protobuf/releases/download/v3.0.2/protoc-3.0.2-osx-x86_64.zip)
# unzip archive into ./tmp/protoc-dl directory
export PATH=$PWD/tmp/protoc-dl/bin/:$PATH
protoc --version

# Generate binaries called out by protoc binary
rm -rf tmp/gen-apiserver-bin/
mkdir -p tmp/gen-apiserver-bin/
go build -o tmp/gen-apiserver-bin/protoc-gen-gogo vendor/github.com/gogo/protobuf/protoc-gen-gogo/main.go
go build -o tmp/gen-apiserver-bin/protoc-gen-gofast vendor/github.com/gogo/protobuf/protoc-gen-gofast/main.go
go build -o tmp/gen-apiserver-bin/goimports vendor/golang.org/x/tools/cmd/goimports/{goimports,goimports_not_gc}.go
export PATH=$PWD/tmp/gen-apiserver-bin/:$PATH

rm -f $(find pkg|grep '\.proto')

go run vendor/k8s.io/code-generator/cmd/go-to-protobuf/main.go \
  --proto-import $PWD/vendor \
  --output-base $(cd ../../../; pwd) \
  --packages "-github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions/v1alpha1,${KC_PKG}/pkg/apis/kappctrl/v1alpha1,${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --go-header-file hack/gen-boilerplate.txt

# TODO It seems that above command messes around with protos in vendor directory
git checkout vendor/

echo "GEN SUCCESS"
