#!/bin/bash

set -e -x -u

# explicitly set CGO_ENABLED to help with reproducible builds see https://github.com/golang/go/issues/36230#issuecomment-567964458
export CGO_ENABLED=0

go fmt ./cmd/... ./pkg/...
go mod vendor
go mod tidy

# pull in function to get the version for use below
source $(dirname "$0")/version-me.sh

# helpful ldflags reference: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
# we set empty buildid and pass -trimpath for reproducible builds; see https://github.com/golang/go/issues/34186
ldflags="-X 'main.Version=$(get_kappctrl_ver)' -buildid="
go build -ldflags="${ldflags}" -trimpath -mod=vendor -o controller ./cmd/main.go

ls -la ./controller

./hack/gen-crds.sh
./hack/ytt-me.sh >/dev/null

echo SUCCESS
