#!/bin/bash

set -euo pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"

if test -z "$BASH_VERSION"; then
  echo "Please run this script using bash, not sh or any other shell." >&2
  exit 1
fi
dst_dir="${CARVEL_INSTALL_BIN_DIR:-${K14SIO_INSTALL_BIN_DIR:-/usr/local/bin}}"

go run "${DIR}/dependencies.go" install --dev --config "${DIR}/dependencies.yml" --destination "${dst_dir}"
