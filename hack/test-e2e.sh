#!/bin/bash

set -e -x -u

go clean -testcache

# create ns if not exists because the `apply -f -` won't complain on a no-op if the ns already exists.
kubectl create ns $KAPPCTRL_E2E_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
go test ./test/e2e/ -timeout 60m -test.v $@

echo E2E SUCCESS
