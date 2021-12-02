#!/bin/bash

set -e -x -u

go clean -testcache

export KCTRL_BINARY_PATH="${KCTRL_BINARY_PATH:-$PWD/kctrl}"

go test ./test/e2e/ -timeout 60m -test.v $@

echo E2E SUCCESS
