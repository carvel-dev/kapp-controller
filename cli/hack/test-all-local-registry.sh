#!/bin/bash

set -e -x -u

function cleanup {
  docker stop registry-"$PORT"
  docker rm -v registry-"$PORT"
}
trap cleanup EXIT

docker run -d -p 5000:5000 -e REGISTRY_VALIDATION_MANIFESTS_URLS_ALLOW='- ^https?://' --restart always --name registry-5000 registry:2
export KCTRL_E2E_IMAGE="localhost:5000/local-tests/test-repo"
./hack/test-all.sh $@