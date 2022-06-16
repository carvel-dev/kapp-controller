#!/bin/bash

set -e -x -u

# explicitly set CGO_ENABLED to help with reproducible builds see https://github.com/golang/go/issues/36230#issuecomment-567964458
export CGO_ENABLED=0

go mod vendor
go mod tidy
go fmt ./cmd/... ./pkg/...

go build -trimpath -mod=vendor -o controller ./cmd/controller/...

ls -la ./controller

./hack/gen-crds.sh
ytt -f config/ >/dev/null

# compile tests, but do not run them: https://github.com/golang/go/issues/15513#issuecomment-839126426
go test --exec=echo ./... >/dev/null

echo SUCCESS
