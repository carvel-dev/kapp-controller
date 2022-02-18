#!/bin/bash

# shamelessly adapted from https://github.com/vmware-tanzu/carvel/blob/develop/site/static/install.sh

DIR="$(dirname "${BASH_SOURCE[0]}")"
DEPENDENCIES_DIR="${DIR}/dependencies.yml"

if test -z "$BASH_VERSION"; then
  echo "Please run this script using bash, not sh or any other shell." >&2
  exit 1
fi

install() {
  set -euo pipefail

  dst_dir="${CARVEL_INSTALL_BIN_DIR:-${K14SIO_INSTALL_BIN_DIR:-/usr/local/bin}}"

  if which sha256sum; then
	  echo "found sha256sum"
  else
    echo "Missing sha256sum binary"
    exit 1
  fi

  ytt_version=$(sed -n -e 's/^ytt_version: //p' "${DEPENDENCIES_DIR}")
  kbld_version=$(sed -n -e 's/^kbld_version: //p' "${DEPENDENCIES_DIR}")
  kapp_version=$(sed -n -e 's/^kapp_version: //p' "${DEPENDENCIES_DIR}")
  imgpkg_version=$(sed -n -e 's/^imgpkg_version: //p' "${DEPENDENCIES_DIR}")
  vendir_version=$(sed -n -e 's/^vendir_version: //p' "${DEPENDENCIES_DIR}")

  if [[ `uname` == Darwin ]]; then
    binary_type=darwin-amd64
    ytt_checksum=$(sed -n -e 's/^ytt_checksum_darwin: //p' "${DEPENDENCIES_DIR}")
    kbld_checksum=$(sed -n -e 's/^kbld_checksum_darwin: //p' "${DEPENDENCIES_DIR}")
    kapp_checksum=$(sed -n -e 's/^kapp_checksum_darwin: //p' "${DEPENDENCIES_DIR}")
    imgpkg_checksum=$(sed -n -e 's/^imgpkg_checksum_darwin: //p' "${DEPENDENCIES_DIR}")
    vendir_checksum=$(sed -n -e 's/^vendir_checksum_darwin: //p' "${DEPENDENCIES_DIR}")
  else
    binary_type=linux-amd64
    ytt_checksum=$(sed -n -e 's/^ytt_checksum_linux: //p' "${DEPENDENCIES_DIR}")
    kbld_checksum=$(sed -n -e 's/^kbld_checksum_linux: //p' "${DEPENDENCIES_DIR}")
    kapp_checksum=$(sed -n -e 's/^kapp_checksum_linux: //p' "${DEPENDENCIES_DIR}")
    imgpkg_checksum=$(sed -n -e 's/^imgpkg_checksum_linux: //p' "${DEPENDENCIES_DIR}")
    vendir_checksum=$(sed -n -e 's/^vendir_checksum_linux: //p' "${DEPENDENCIES_DIR}")
  fi

  echo "Installing ${binary_type} binaries..."

  echo "Installing ytt..."
  curl -sLo /tmp/ytt https://github.com/vmware-tanzu/carvel-ytt/releases/download/${ytt_version}/ytt-${binary_type}
  echo "${ytt_checksum}  /tmp/ytt" | sha256sum -c -
  mv /tmp/ytt ${dst_dir}/ytt
  chmod +x ${dst_dir}/ytt
  echo "Installed ${dst_dir}/ytt ${ytt_version}"

  echo "Installing kbld..."
  curl -sLo /tmp/kbld https://github.com/vmware-tanzu/carvel-kbld/releases/download/${kbld_version}/kbld-${binary_type}
  echo "${kbld_checksum}  /tmp/kbld" | sha256sum -c -
  mv /tmp/kbld ${dst_dir}/kbld
  chmod +x ${dst_dir}/kbld
  echo "Installed ${dst_dir}/kbld ${kbld_version}"

  echo "Installing kapp..."
  curl -sLo /tmp/kapp https://github.com/vmware-tanzu/carvel-kapp/releases/download/${kapp_version}/kapp-${binary_type}
  echo "${kapp_checksum}  /tmp/kapp" | sha256sum -c -
  mv /tmp/kapp ${dst_dir}/kapp
  chmod +x ${dst_dir}/kapp
  echo "Installed ${dst_dir}/kapp ${kapp_version}"

  echo "Installing imgpkg..."
  curl -sLo /tmp/imgpkg https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/${imgpkg_version}/imgpkg-${binary_type}
  echo "${imgpkg_checksum}  /tmp/imgpkg" | sha256sum -c -
  mv /tmp/imgpkg ${dst_dir}/imgpkg
  chmod +x ${dst_dir}/imgpkg
  echo "Installed ${dst_dir}/imgpkg ${imgpkg_version}"

  echo "Installing vendir..."
  curl -sLo /tmp/vendir https://github.com/vmware-tanzu/carvel-vendir/releases/download/${vendir_version}/vendir-${binary_type}
  echo "${vendir_checksum}  /tmp/vendir" | sha256sum -c -
  mv /tmp/vendir ${dst_dir}/vendir
  chmod +x ${dst_dir}/vendir
  echo "Installed ${dst_dir}/vendir ${vendir_version}"
}

install
