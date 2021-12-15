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
  kbld_version=v0.32.0
  kapp_version=v0.43.0
  imgpkg_version=v0.24.0
  vendir_version=v0.24.0

  if [[ `uname` == Darwin ]]; then
    binary_type=darwin-amd64
    ytt_checksum=5eef9da11dd4f714495e2cf47041fc6fd413c8c393af4cc5bf3869e075b4e936
    kbld_checksum=5fc8a491327294717611974c6ab3da2bda3f3809ef3147c1e8472ac62af3ee18
    kapp_checksum=c5c7f34399293ccda62dc7b809535a8e2f2afb9901147429281c0f4884b13483
    imgpkg_checksum=f0c87c8caefb3d2a82e648779b36783403fe5c93930df2d5cbf4968713933392
    vendir_checksum=f3a738d1fe55803ad5faba495f662c48efa230976ccad7a159587dcf9b020f63
  else
    binary_type=linux-amd64
    ytt_checksum=1aad12386f6bae1a78197acdc7ec9e60c441f82c4ca944df8d3c78625750fe59
    kbld_checksum=de546ac46599e981c20ad74cd2deedf2b0f52458885d00b46b759eddb917351a
    kapp_checksum=f8669039dfba001081c94576c898d10aba28ecceffcd98708e8f2c87c13109e4
    imgpkg_checksum=cfcfcb5afc5e3d28ce1f2f67971a4dcd18f514dadf8a63d70c864e49c9ddca7e
    vendir_checksum=b7bfd227aa2e6df602f8e79edf725bb0a944b68d207005f42f46f061c4ecd55a
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
