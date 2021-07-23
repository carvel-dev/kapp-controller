#!/bin/bash

set -e

./hack/build.sh && ./hack/ytt-me.sh | kbld -f- | kapp deploy -a kc -f- -c -y
