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

  ytt_version=v0.35.1
  kbld_version=v0.30.0
  kapp_version=v0.38.0
  imgpkg_version=v0.14.0
  vendir_version=v0.22.0

  if [[ `uname` == Darwin ]]; then
    binary_type=darwin-amd64
    ytt_checksum=1f2b61d02f6d8184889719d5e0277a1ea82219f96873345157e81075ca59808e
    kbld_checksum=73274d02b0c2837d897c463f820f2c8192e8c3f63fd90c526de5f23d4c6bdec4
    kapp_checksum=2c7c9faf6b5bc564ee6a9450c1e21c16aa97c138ea59629441f8f28876bed6ad
    imgpkg_checksum=3c63f224833590526c3b48ab5db1be4ec07ece6a6eb4879888fefba14b6176be
    vendir_checksum=66cc6519c924897425c4750c197ea4c7f4e07e9275789f6a2f1a0b7db437c636
  else
    binary_type=linux-amd64
    ytt_checksum=0aa78f7b5f5a0a4c39bddfed915172880344270809c26b9844e9d0cbf6437030
    kbld_checksum=76c5c572e7a9095256b4c3ae2e076c370ef70ce9ff4eb138662f56828889a00c
    kapp_checksum=22e3d694745d5f48863018e26ecd7f3d0b8ec475adc40e081a1a39dc4d8f01bf
    imgpkg_checksum=bd53355fc3a05666681ddf2ba1dfae2da894bc1c74d86cdc545d772749abc887
    vendir_checksum=951b75467ac8be6022efe3584815ef4ea285a0e3b591eba7f775c55c4947c2ed
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
