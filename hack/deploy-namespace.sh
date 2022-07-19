#!/bin/bash

set -e

./hack/build.sh && ytt -f config-namespace/ | kbld -f- | kapp deploy -a kc -f- -c -y

