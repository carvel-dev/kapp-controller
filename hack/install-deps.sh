#!/bin/zsh

# shamelessly adapted from https://github.com/vmware-tanzu/carvel/blob/develop/site/static/install.sh

install() {
  set -euo pipefail

  dst_dir="${K14SIO_INSTALL_BIN_DIR:-/usr/local/bin}"

  wget --version 1>/dev/null 2>&1 || (echo "Missing wget binary" && exit 1)
  shasum -v 1>/dev/null 2>&1 || (echo "Missing shasum binary" && exit 1)

  carvel_tools=(ytt kbld kapp imgpkg vendir)

  declare -A versions
  declare -A checksums

  versions[ytt]=v0.35.1
  versions[kbld]=v0.30.0
  versions[kapp]=v0.37.0
  versions[imgpkg]=v0.14.0
  versions[vendir]=v0.21.1

  if [[ `uname` == Darwin ]]; then
    binary_type=darwin-amd64
    checksums[ytt]=1f2b61d02f6d8184889719d5e0277a1ea82219f96873345157e81075ca59808e
    checksums[kbld]=73274d02b0c2837d897c463f820f2c8192e8c3f63fd90c526de5f23d4c6bdec4
    checksums[kapp]=da6411b79c66138cd7437beb268675edf2df3c0a4a8be07fb140dd4ebde758c1
    checksums[kwt]=555d50d5bed601c2e91f7444b3f44fdc424d721d7da72955725a97f3860e2517
    checksums[imgpkg]=3c63f224833590526c3b48ab5db1be4ec07ece6a6eb4879888fefba14b6176be
    checksums[vendir]=579d661291f669a4f618c602e48d456693762e2ba23d4d8b64d12ceea05dd2cd
  else
    binary_type=linux-amd64
    checksums[ytt]=0aa78f7b5f5a0a4c39bddfed915172880344270809c26b9844e9d0cbf6437030
    checksums[kbld]=76c5c572e7a9095256b4c3ae2e076c370ef70ce9ff4eb138662f56828889a00c
    checksums[kapp]=f845233deb6c87feac7c82d9b3f5e03ced9a4672abb1a14d4e5b74fe53bc4538
    checksums[kwt]=92a1f18be6a8dca15b7537f4cc666713b556630c20c9246b335931a9379196a0
    checksums[imgpkg]=bd53355fc3a05666681ddf2ba1dfae2da894bc1c74d86cdc545d772749abc887
    checksums[vendir]=7d9ffd06a888bf13e16ad964d7a0d0f6b7c23e8cad9774084c563cda81b91184
  fi

  echo "Installing ${binary_type} binaries..."

  for i in ${carvel_tools[@]}; do
    echo "Installing ${i}..."
    wget -nv -O- https://github.com/vmware-tanzu/carvel-${i}/releases/download/${versions[${i}]}/${i}-${binary_type} > /tmp/${i}
    echo "${checksums[${i}]}  /tmp/${i}" | shasum -c -
    mv /tmp/${i} ${dst_dir}/${i}
    chmod +x ${dst_dir}/${i}
    echo "Installed ${dst_dir}/${i} ${versions[${i}]}"
  done
}

install
