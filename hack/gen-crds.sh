#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Using v0.3.0 of controller-gen
# Install: go get -u sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0
# This command generates the crds under the config folder and also
# provides a schema definition for kapp-controller crds.
controller-gen \
  crd \
  output:dir=./config \
  paths=./pkg/apis/...
