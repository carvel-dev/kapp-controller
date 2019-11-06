#!/bin/bash

set -e -x -u

go fmt ./cmd/... ./pkg/...

go build -o controller ./cmd/controller/...
ls -la ./controller

ytt -f config/ >/dev/null

echo SUCCESS
