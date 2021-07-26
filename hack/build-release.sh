#!/bin/bash

set -e -x -u

mkdir -p tmp/

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-me.sh

ytt -f config/ -f config-release -v kapp_controller_version="$(get_kappctrl_ver)" | kbld -f- > ./tmp/release.yml

shasum -a 256 ./tmp/release*.yml

echo SUCCESS
