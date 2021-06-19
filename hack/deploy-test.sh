#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ -f config-test/ | kbld -f- | kapp deploy -a kc -f- -c -y
