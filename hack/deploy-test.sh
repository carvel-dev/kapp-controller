#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ -f config-test/ | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
