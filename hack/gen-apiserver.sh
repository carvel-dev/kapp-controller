#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

GOPATH=`go env GOPATH`
KC_PKG="github.com/vmware-tanzu/carvel-kapp-controller"

# Generate clientset and apis code with K8s codegen tools.
echo "Generating clients"
$GOPATH/bin/client-gen \
  --clientset-name versioned \
  --input-base "${KC_PKG}/pkg/apiserver/apis/" \
  --input "datapackaging/v1alpha1" \
  --output-package "${KC_PKG}/pkg/apiserver/client/clientset" \
  --go-header-file hack/gen-boilerplate.txt

# Generate listers with K8s codegen tools.
echo "Generating listers"
$GOPATH/bin/lister-gen \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --output-package "${KC_PKG}/pkg/apiserver/client/listers" \
  --go-header-file hack/gen-boilerplate.txt

# Generate informers with K8s codegen tools.
echo "Generating informers"
$GOPATH/bin/informer-gen \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --versioned-clientset-package "${KC_PKG}/pkg/apiserver/client/clientset/versioned" \
  --listers-package "${KC_PKG}/pkg/apiserver/client/listers" \
  --output-package "${KC_PKG}/pkg/apiserver/client/informers" \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating deepcopy"
$GOPATH/bin/deepcopy-gen \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging" \
  -O zz_generated.deepcopy \
  --go-header-file hack/gen-boilerplate.txt

echo "Generating conversions"
$GOPATH/bin/conversion-gen  \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1,${KC_PKG}/pkg/apiserver/apis/datapackaging" \
  -O zz_generated.conversion \
  --go-header-file hack/gen-boilerplate.txt

# echo "Generating openapi"
$GOPATH/bin/openapi-gen  \
  --input-dirs "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1" \
  --input-dirs "${KC_PKG}/pkg/apis/kappctrl/v1alpha1" \
  --input-dirs "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr" \
  --input-dirs "k8s.io/api/core/v1" \
  --output-package "${KC_PKG}/pkg/apiserver/openapi" \
  -O zz_generated.openapi \
  --go-header-file hack/gen-boilerplate.txt

$GOPATH/bin/go-to-protobuf \
  --proto-import vendor \
  --packages "${KC_PKG}/pkg/apiserver/apis/datapackaging/v1alpha1,${KC_PKG}/pkg/apis/kappctrl/v1alpha1" \
  --go-header-file hack/gen-boilerplate.txt

