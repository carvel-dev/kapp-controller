#!/bin/bash

set -e -x -u

mkdir -p tmp/
mkdir -p .imgpkg/

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

export version="$(get_kappctrl_ver)"

yq eval '.metadata.annotations."kapp-controller.carvel.dev/version" = env(version)' -i "config/deployment.yml"

ytt -f config/ -f config-release -v dev.kapp_controller_version="$(get_kappctrl_ver)" --data-values-env=KCTRL | kbld --imgpkg-lock-output .imgpkg/images.yml -f- > ./tmp/release.yml
shasum -a 256 ./tmp/release.yml

echo SUCCESS
