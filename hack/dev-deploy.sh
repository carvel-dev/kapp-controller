#!/bin/bash

set -ex

source $(dirname "$0")/version-util.sh

rm -rf tmp/build
mkdir -p tmp/build
CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-X 'main.Version=$(get_kappctrl_ver)+develop'" -trimpath -o tmp/build ./cmd/...

kc_latest_image=`docker image ls --filter=reference="*:kapp-controller-*" --format "{{.Repository}}:{{.Tag}}" | head -n 1`
if [ -z "$kc_latest_image" ] ;
then
  echo "Error: unable to find tag for previous image of kapp-controller"
  echo "For your first deploy please use hack/deploy.sh and then try re-running this script for subsequent deploys."
  exit 1
fi

echo "got kc latest image: $kc_latest_image"

cat << EOF > tmp/build/Dockerfile
FROM ${kc_latest_image} AS build
FROM scratch
COPY --from=build / /
COPY controller /kapp-controller
USER 1000
ENV PATH="/:\${PATH}"
ENTRYPOINT ["/kapp-controller"]
EOF
cat << EOF > tmp/build/overlay.yml
#@ load("@ytt:overlay", "overlay")

#@ def find_image_sources():
kind: Config
sources:
  - image: kapp-controller
#@ end

#@overlay/match by=overlay.subset(find_image_sources())
---
sources:
#@overlay/match by="image"
- image: kapp-controller
  path: tmp/build
  docker:
    buildx:
      pull: false
EOF

ytt -f config/ -f tmp/build/overlay.yml -v kapp_controller_version="$(get_kappctrl_ver)+develop" | kbld -f- | kapp deploy -a kc -f- -c -y

