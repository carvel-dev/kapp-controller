#!/bin/bash

set -e

./hack/build.sh && ytt -f config/ --data-value-yaml=push_images=true | kbld -f- | kapp deploy -a kc -f- -c -y
