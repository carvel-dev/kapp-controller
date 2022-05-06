#!/bin/bash

set -e

image="${KC_IMG:-carvel/kapp-controller:dev}"
docker build -t "${image}" .
if [ -n "${KC_IMG_PUSH}" ]; then
  docker push -t "${image}"
fi

./hack/build.sh && ytt -f config/ -f config-test/ -v image="${image}" | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
