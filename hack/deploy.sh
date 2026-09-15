#!/bin/bash

set -e

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

./hack/build.sh
rendered=$(ytt -f config/config -f config/values-schema.yml -f config-dev -v dev.version="$(get_kappctrl_ver)+develop" | kbld -f-)

# kbld builds locally and retags the image as "kbld:<image>-<sha>" (see
# RetagStable in carvel-dev/kbld). When targeting minikube, side-load that
# straight into the node instead of relying on `eval $(minikube docker-env)`,
# which breaks under the containerd runtime.
# See https://github.com/carvel-dev/kapp-controller/issues/1858
if [ "${KAPPCTRL_MINIKUBE_LOAD_IMAGES:-}" = "true" ]; then
  for img in $(echo "${rendered}" | grep -oE 'kbld:[a-zA-Z0-9_.-]+' | sort -u); do
    minikube image load "${img}"
  done
fi

echo "${rendered}" | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
