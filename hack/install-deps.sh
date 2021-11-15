#!/bin/bash

# shamelessly adapted from https://github.com/vmware-tanzu/carvel/blob/develop/site/static/install.sh

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

  ytt_version=v0.37.0
  kbld_version=v0.31.0
  kapp_version=v0.42.0
  imgpkg_version=v0.22.0
  vendir_version=v0.23.0

  if [[ `uname` == Darwin ]]; then
    binary_type=darwin-amd64
    ytt_checksum=5eef9da11dd4f714495e2cf47041fc6fd413c8c393af4cc5bf3869e075b4e936
    kbld_checksum=d3b0a30bf3a79bedeb25d8548a91254954b99cd4a0c03f3a810b331fc4d1f071
    kapp_checksum=47102637b9cd541b4ad1d6074f77b7cec1b60c170a0eb5c92df89674207194e7
    imgpkg_checksum=bc2819d73c6680c402a1ff4a51de07a5decafaef70b34662b46460c9a5eccc4c
    vendir_checksum=aef3233128727c01ffac6561533c9f60a49307f41decc1baa18688a4a250b15f
  else
    binary_type=linux-amd64
    ytt_checksum=1aad12386f6bae1a78197acdc7ec9e60c441f82c4ca944df8d3c78625750fe59
    kbld_checksum=ba0be56d9e74b067f3e659de0b79100b0b9df86a2e3e0e6ff533b1e019c22c23
    kapp_checksum=5d5c4274a130f2fd5ad11ddd8fb3e0f647c8598ba25711360207fc6eab72f6be
    imgpkg_checksum=d34a06097a781ed15a4f04d6fe72b41ccf9dae835c821ec524a99b6c4e647f03
    vendir_checksum=c6a65e7e8e589e25bf5554e9575ab247e55ae71920d3d366ffd84b1571fe34ac
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
