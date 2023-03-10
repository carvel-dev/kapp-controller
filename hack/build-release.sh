#!/bin/bash

set -e -x -u

mkdir -p tmp/
mkdir -p .imgpkg/

# makes the get_kappctrl_ver function available (scrapes version from git tag)
source $(dirname "$0")/version-util.sh

export version="$(get_kappctrl_ver)"

# We do not want the version to be configurable in the kapp-controller package
sed 's/v0.0.0/'"$version"'/' config/config/deployment.yml > tmp/deployment.yml
mv tmp/deployment.yml config/config/deployment.yml

ytt -f config/config -f config/values-schema.yml -f config-release -v dev.version="$version" --data-values-env=KCTRL | kbld --imgpkg-lock-output .imgpkg/images.yml -f- > ./tmp/release.yml

# Update image url in kapp-controller package overlays
image_url=`yq e '.spec.template.spec.containers[] | select(.name == "kapp-controller") | .image' ./tmp/release.yml`
sed 's|image: kapp-controller|image: '"$image_url"'|' config/overlays/update-deployment.yml > tmp/update-deployment.yml
mv tmp/update-deployment.yml config/overlays/update-deployment.yml

shasum -a 256 ./tmp/release.yml

echo SUCCESS
