#!/bin/bash

set -e

source $(dirname "$0")/version-util.sh

cat <<EOF >overlay.yml
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
    kapp-controller.carvel.dev/version: #@ data.values.version
EOF

cat <<EOF >values.yml
#@data/values
---
dev:
  push_images: false
  image_cache: true
  platform: ""
EOF

./hack/build.sh && ytt -f config/ -f config-release/ -f values.yml > ./tmp/config.yml

ytt -f ./tmp/config.yml -f overlay.yml -v version="$(get_kappctrl_ver)+develop" | kbld -f- | kapp deploy -a kc -f- -c -y

source ./hack/secretgen-controller.sh
deploy_secretgen-controller
