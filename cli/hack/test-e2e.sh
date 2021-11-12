#!/bin/bash

set -e -x -u

go clean -testcache

export KAPPCTRL_BINARY_PATH="${KAPPCTRL_BINARY_PATH:-$PWD/kapp}"

go test ./test/e2e/ -timeout 60m -test.v $@

echo E2E SUCCESS
