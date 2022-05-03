#!/bin/bash

set -e

export DOCKER_BUILDKIT=1
./hack/build.sh && ytt -f config/ -f config-test/ | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
