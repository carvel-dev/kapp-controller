#!/bin/bash

set -ex

CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-X 'main.Version=develop' -buildid=" -trimpath -o tmp/kapp-controller ./cmd/main.go

run_image_start="$(yq eval '.image_repo' config/values.yml)"
kc_latest_image="$(docker image ls $run_image_start --format "{{.Repository}}@{{.Digest}}" | head -1)"
if [ -z "$kc_latest_image" ] ;
then
  echo "Error: unable to find tag for previous image of kapp-controller"
  echo "For your first deploy please use hack/deploy.sh and then try re-running this script for subsequent deploys."
  exit 1
fi
tar -cf tmp/kc.tar -C tmp kapp-controller
crane append -b "$kc_latest_image" -f tmp/kc.tar -t "${run_image_start}:latest"

ytt -f config/deployment.yml -f config/values.yml --data-value image="${run_image_start}:latest" | kbld -f- | kubectl apply -f-
