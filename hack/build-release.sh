#!/bin/bash

set -e -x -u

mkdir -p tmp/

ytt -f config/ -f config-release | kbld -f- > ./tmp/release.yml

shasum -a 256 ./tmp/release*.yml

echo SUCCESS
