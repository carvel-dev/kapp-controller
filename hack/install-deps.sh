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

  if [ -x "$(command -v wget)" ]; then
    dl_bin="wget -nv -O-"
  else
    dl_bin="curl -s -L"
  fi

  if which sha256sum; then
	  echo "found sha256sum"
  else
    echo "Missing sha256sum binary"
    exit 1
  fi

  if which yq; then
	  echo "found yq"
  else
    echo "Missing yq binary"
    exit 1
  fi

  ytt_version=$(awk '/ytt_version/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
  kbld_version=$(awk '/kbld_version/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
  kapp_version=$(awk '/kapp_version/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
  imgpkg_version=$(awk '/imgpkg_version/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
  vendir_version=$(awk '/vendir_version/{print $(NF-0)}' "${DEPENDENCIES_DIR}")

  if [[ `uname` == Darwin ]]; then
    binary_type=darwin-amd64
    ytt_checksum=$(awk '/ytt_checksum_darwin/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    kbld_checksum=$(awk '/kbld_checksum_darwin/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    kapp_checksum=$(awk '/kapp_checksum_darwin/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    imgpkg_checksum=$(awk '/imgpkg_checksum_darwin/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    vendir_checksum=$(awk '/vendir_checksum_darwin/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
  else
    binary_type=linux-amd64
    ytt_checksum=$(awk '/ytt_checksum_linux/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    kbld_checksum=$(awk '/kbld_checksum_linux/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    kapp_checksum=$(awk '/kapp_checksum_linux/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    imgpkg_checksum=$(awk '/imgpkg_checksum_linux/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
    vendir_checksum=$(awk '/vendir_checksum_linux/{print $(NF-0)}' "${DEPENDENCIES_DIR}")
  fi

  echo "Installing ${binary_type} binaries..."

  echo "Installing ytt..."
  $dl_bin https://github.com/vmware-tanzu/carvel-ytt/releases/download/${ytt_version}/ytt-${binary_type} > /tmp/ytt
  echo "${ytt_checksum}  /tmp/ytt" | sha256sum -c -
  mv /tmp/ytt ${dst_dir}/ytt
  chmod +x ${dst_dir}/ytt
  echo "Installed ${dst_dir}/ytt ${ytt_version}"

  echo "Installing kbld..."
  $dl_bin https://github.com/vmware-tanzu/carvel-kbld/releases/download/${kbld_version}/kbld-${binary_type} > /tmp/kbld
  echo "${kbld_checksum}  /tmp/kbld" | sha256sum -c -
  mv /tmp/kbld ${dst_dir}/kbld
  chmod +x ${dst_dir}/kbld
  echo "Installed ${dst_dir}/kbld ${kbld_version}"

  echo "Installing kapp..."
  $dl_bin https://github.com/vmware-tanzu/carvel-kapp/releases/download/${kapp_version}/kapp-${binary_type} > /tmp/kapp
  echo "${kapp_checksum}  /tmp/kapp" | sha256sum -c -
  mv /tmp/kapp ${dst_dir}/kapp
  chmod +x ${dst_dir}/kapp
  echo "Installed ${dst_dir}/kapp ${kapp_version}"

  echo "Installing imgpkg..."
  $dl_bin https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/${imgpkg_version}/imgpkg-${binary_type} > /tmp/imgpkg
  echo "${imgpkg_checksum}  /tmp/imgpkg" | sha256sum -c -
  mv /tmp/imgpkg ${dst_dir}/imgpkg
  chmod +x ${dst_dir}/imgpkg
  echo "Installed ${dst_dir}/imgpkg ${imgpkg_version}"

  echo "Installing vendir..."
  $dl_bin https://github.com/vmware-tanzu/carvel-vendir/releases/download/${vendir_version}/vendir-${binary_type} > /tmp/vendir
  echo "${vendir_checksum}  /tmp/vendir" | sha256sum -c -
  mv /tmp/vendir ${dst_dir}/vendir
  chmod +x ${dst_dir}/vendir
  echo "Installed ${dst_dir}/vendir ${vendir_version}"
}

install
