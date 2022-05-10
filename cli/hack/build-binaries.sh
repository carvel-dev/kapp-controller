#!/bin/bash

set -e -x -u

./hack/build.sh

function get_latest_git_tag {
  git describe --tags | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+'
}

VERSION="${1:-`get_latest_git_tag`}"

# makes builds reproducible
export CGO_ENABLED=0
LDFLAGS="-X github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/version.Version=$VERSION"

GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kctrl-darwin-amd64 ./cmd/kctrl/...
GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kctrl-darwin-arm64 ./cmd/kctrl/...
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kctrl-linux-amd64 ./cmd/kctrl/...
GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kctrl-linux-arm64 ./cmd/kctrl/...
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kctrl-windows-amd64.exe ./cmd/kctrl/...

shasum -a 256 ./kctrl-darwin-amd64 ./kctrl-darwin-arm64 ./kctrl-linux-amd64 ./kctrl-linux-arm64 ./kctrl-windows-amd64.exe
