#!/bin/bash

set -e -x -u

# explicitly set CGO_ENABLED to help with reproducible builds see https://github.com/golang/go/issues/36230#issuecomment-567964458
export CGO_ENABLED=0

go fmt ./cmd/... ./pkg/...
go mod vendor
go mod tidy

# we set empty buildid and pass -trimpath for reproducible builds; see https://github.com/golang/go/issues/34186
go build -ldflags="-buildid=" -trimpath -mod=vendor -o controller ./cmd/main.go

ls -la ./controller

./hack/gen-crds.sh
ytt -f config/ >/dev/null

echo SUCCESS
