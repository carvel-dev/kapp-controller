#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go run ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go \
  crd \
  output:dir=./config \
  paths=./pkg/apis/...
