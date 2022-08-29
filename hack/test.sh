#!/bin/bash

set -e -x

if [ -z "$GITHUB_ACTION" ]; then
  go clean -testcache
fi

set -u

GO=go
if command -v richgo &> /dev/null
then
    GO=richgo
fi

$GO test ./pkg/... ./cmd/... ./hack/... $@

echo UNIT SUCCESS
