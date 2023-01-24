#!/bin/bash

set -e -x -u

mkdir -p tmp/
mkdir -p .imgpkg/

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

export version="$(get_kappctrl_ver)"

ytt -f config/ -f config-release -v dev.kapp_controller_version="$(get_kappctrl_ver)" --data-values-env=KCTRL | kbld --imgpkg-lock-output .imgpkg/images.yml -f- > ./tmp/pre-release.yml

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

ytt -f ./tmp/pre-release.yml -f overlay.yml -v version="$(get_kappctrl_ver)" > ./tmp/release.yml
shasum -a 256 ./tmp/release.yml

echo SUCCESS
