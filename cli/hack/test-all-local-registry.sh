#!/bin/bash

set -e -x -u

PORT=5000

function cleanup {
  docker stop registry-"$PORT"
  docker rm -v registry-"$PORT"
}
trap cleanup EXIT

docker run -d -p "$PORT":5000 --restart always --name registry-"$PORT" registry:2
export KCTRL_E2E_IMAGE="`ifconfig | grep inet | grep -E '\b10\.' | awk '{ print $2}'`:5000/test-repo/testimage"
./hack/test-all.sh $@
