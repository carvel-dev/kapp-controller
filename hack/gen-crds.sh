#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go run ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go \
  crd:preserveUnknownFields=false \
  output:dir=./config \
  paths=./pkg/apis/...

# Manual steps post generation for all CRDs:
#
# 1. Remove empty status from generated CRDs (see https://github.com/kubernetes-sigs/controller-tools/issues/456):
# status:
#   acceptedNames:
#     kind: ""
#     plural: ""
#   conditions: []
#   storedVersions: []
#
# 2. Remove metadata from schema (see https://github.com/kubernetes/kubernetes/issues/80493):
# metadata:
#   description: 'Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.'
#   type: object
#
# 3. Remove annotations:
#  annotations:
#    controller-gen.kubebuilder.io/version: (devel)
#
# 4. Remove metadata creationTimestamp: null from metadata section