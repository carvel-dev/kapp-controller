#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
