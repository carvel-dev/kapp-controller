#!/bin/bash

set -e

source $(dirname "$0")/version-util.sh

./hack/build.sh && ytt -f config/config -f config/values-schema.yml -f config-dev --data-value-yaml dev.push_images=true -v dev.version="$(get_kappctrl_ver)+develop" -v dev.image_repo="us-central1-docker.pkg.dev/cf-k8s-lifecycle-tooling-klt/kapp-controller-tests/kapp-controller" | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
