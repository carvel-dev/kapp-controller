#!/bin/bash

set -e -x -u

mkdir -p tmp/
mkdir -p .imgpkg/

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

export version="$(get_kappctrl_ver)"

# We do not want the version to be configurable in the kapp-controller package
sed 's/v0.0.0/'"$version"'/' config/deployment.yml > tmp/deployment.yml
mv tmp/deployment.yml config/deployment.yml

ytt -f config -f config-release -v dev.version="$version" --data-values-env=KCTRL | kbld --imgpkg-lock-output .imgpkg/images.yml -f- > ./tmp/release.yml

shasum -a 256 ./tmp/release.yml

echo SUCCESS
