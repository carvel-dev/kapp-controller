#!/bin/bash

set -e

source $(dirname "$0")/version-util.sh

./hack/build.sh && ytt -f config/ -v kapp_controller_version="$(get_kappctrl_ver)+develop" | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
