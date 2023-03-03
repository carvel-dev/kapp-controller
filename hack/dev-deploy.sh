#!/bin/bash

set -ex

source $(dirname "$0")/version-util.sh

CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-X 'main.Version=$(get_kappctrl_ver)+develop'" -trimpath -o controller-linux-amd64 ./cmd/controller/...

kc_latest_image=`docker image ls --filter=reference="*:kapp-controller-*" --format "{{.Repository}}:{{.Tag}}" | head -n 1`
if [ -z "$kc_latest_image" ] ;
then
  echo "Error: unable to find tag for previous image of kapp-controller"
  echo "For your first deploy please use hack/deploy.sh and then try re-running this script for subsequent deploys."
  exit 1
fi

echo "got kc latest image: $kc_latest_image"

cat << EOF > Dockerfile.dev
FROM ${kc_latest_image} AS build
FROM scratch
COPY --from=build / /
COPY controller-linux-amd64 /kapp-controller
USER 1000
ENV PATH="/:\${PATH}"
ENTRYPOINT ["/kapp-controller"]
EOF

ytt -f config/config -f config/values-schema.yml -f config-dev -v dev.version="$(get_kappctrl_ver)+develop" --data-value-yaml dev.rapid_deploy=true | kbld -f- | kapp deploy -a kc -f- -c -y
