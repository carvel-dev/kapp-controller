#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ | kbld -f- | kapp deploy -a kc -f- -c -y --dangerous-override-ownership-of-existing-resources

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
