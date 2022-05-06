#!/bin/bash

set -e -x -u

mkdir -p tmp/

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

ytt -f config/ -v image="${IMAGE}" -v kapp_controller_version="$(get_kappctrl_ver)" > ./tmp/release.yml

shasum -a 256 ./tmp/release*.yml | tee ./tmp/checksums.txt

echo SUCCESS
