#!/bin/bash

set -e

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

cat << EOF > overlay.yml
#@ load("@ytt:overlay", "overlay")
#@ load("@ytt:data", "data")
#@ def resource(kind, name):
kind: #@ kind
metadata:
  name: #@ name
#@ end
#@overlay/match by=overlay.subset(resource("Deployment", "kapp-controller"))
---
metadata:
  annotations:
    #@overlay/match missing_ok=True
    kapp-controller.carvel.dev/version: #@ data.values.dev.kapp_controller_version
EOF

./hack/build.sh && ytt -f config/ -f config-release/ -f overlay.yml -v dev.kapp_controller_version="$(get_kappctrl_ver)+develop" | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
