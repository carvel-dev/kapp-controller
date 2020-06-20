#!/bin/bash

set -e -x -u

mkdir -p tmp/

ytt -f config/ -f config-release | kbld -f- > ./tmp/release.yml --lock-output ./tmp/images.yml

shared_sa_flags="--data-value-yaml dangerous_allow_shared_service_account=true"
ytt -f config/ -f config-release $shared_sa_flags | kbld -f- -f ./tmp/images.yml > ./tmp/release-dangerous-allow-shared-sa.yml

shasum -a 256 ./tmp/release*.yml

echo SUCCESS
