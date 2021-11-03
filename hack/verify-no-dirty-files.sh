#!/bin/bash

set -e

./hack/build.sh
./hack/gen.sh
./hack/gen-apiserver.sh

if ! git diff --exit-code >/dev/null; then
  echo "Error: Running ./hack/verify-no-dirty-files.sh resulted in non zero exit code from git diff."
  echo "Please run './hack/build.sh', './hack/gen.sh', and './hack/gen-apiserver.sh' and 'git add' the generated file(s)."
  echo "Showing diff:"
  git diff --exit-code
  exit 1
fi
