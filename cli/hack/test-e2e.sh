#!/bin/bash

set -e -x -o pipefail

go clean -testcache

export KCTRL_BINARY_PATH="${KCTRL_BINARY_PATH:-$PWD/kctrl}"

if [ -z "$KCTRL_E2E_NAMESPACE" ]; then
    echo "setting e2e namespace to: kctrl-test";
    export KCTRL_E2E_NAMESPACE="kctrl-test"
fi
# create ns if not exists because the `apply -f -` won't complain on a no-op if the ns already exists.
kubectl create ns $KCTRL_E2E_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

go test ./test/e2e/ -timeout 60m -test.v $@

echo E2E SUCCESS
