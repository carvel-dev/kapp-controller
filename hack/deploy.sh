#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ -v push_images=true --data-value-yaml dangerous_enable_pprof=true | kbld -f- | kapp deploy -a kc -f- -c -y
