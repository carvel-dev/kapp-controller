#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ -f cmd/main.go | kbld -f- | kapp deploy -a kc -f- -c -y
