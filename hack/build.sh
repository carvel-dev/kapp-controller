#!/bin/bash

set -e -x -u

# makes builds reproducible
export CGO_ENABLED=0
repro_flags="-ldflags=-buildid= -trimpath"

go fmt ./cmd/... ./pkg/...
go mod vendor
go mod tidy

go build $repro_flags -mod=vendor -o controller ./cmd/main.go
ls -la ./controller

./hack/gen-crds.sh
ytt -f config/ >/dev/null

echo SUCCESS
