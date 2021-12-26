#!/bin/bash

set -e -x -u

go clean -testcache

# create ns if not exists because the `apply -f -` won't complain on a no-op if the ns already exists.
kubectl create ns $KAPPCTRL_E2E_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
go test ./test/e2e/kappcontroller -timeout 60m -test.v $@ | tee tmp/e2eoutput.txt

if [ "$KAPPCTRL_E2E_SECRETGEN_CONTROLLER" == "true" ]; then
  go test ./test/e2e/secretgencontroller -timeout 60m -test.v $@ | tee -a tmp/e2eoutput.txt
fi

echo E2E SUCCESS
