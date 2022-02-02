#!/usr/bin/env bash

set -e

./hack/build.sh
./hack/gen.sh
./hack/gen-apiserver.sh
