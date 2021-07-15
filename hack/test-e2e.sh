#!/bin/bash

set -e -x -u

go clean -testcache

# create ns if not exists because the `apply -f -` won't complain if fed an empty yaml to apply due to the ns already existing.
kubectl create ns $KAPPCTRL_E2E_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
go test ./test/e2e/ -timeout 60m -test.v $@

echo E2E SUCCESS
