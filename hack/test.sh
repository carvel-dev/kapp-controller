#!/bin/bash

set -e -x

if [ -z "$GITHUB_ACTION" ]; then
  go clean -testcache
fi

set -u

export KAPPCTRL_API_PORT="90210"

go test ./pkg/... ./cmd/... -test.v $@

echo UNIT SUCCESS
