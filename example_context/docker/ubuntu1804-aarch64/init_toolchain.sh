#!/usr/bin/env bash

set -e

TOOLS_INSTALL_DIR="$1"
TMP_DIR="/tmp/toolchain-install"

TOOLS_PACKAGE_URI="https://github.com/bacpack-system/packager-tools/releases/download/v1.0.0/bringauto-packager-tools_v1.0.0_x86-64-linux.zip"

if [[ ${TOOLS_INSTALL_DIR} = "" ]]
then
  echo "Specify tools install dir absolute path as a second argument!" >&2
  exit 1
fi

function install_tools() {
  if [[ -d ${TMP_DIR} ]]
  then
    echo "TMP dir '${TMP_DIR}' exist"
  fi
  mkdir -p "${TMP_DIR}"
  mkdir -p "${TOOLS_INSTALL_DIR}"

  pushd "${TMP_DIR}"
    wget ${TOOLS_PACKAGE_URI} \
      -O "bringauto-packager-tools.zip"
    unzip bringauto-packager-tools.zip
    rm bringauto-packager-tools.zip
    directory_name="$(echo ./*)"
    mv "${directory_name}"/* "${TOOLS_INSTALL_DIR}"
    rm -r "${directory_name}"
  popd
  rm -r "${TMP_DIR}"

  chmod +x "${TOOLS_INSTALL_DIR}/uname"
  TMP_PATH=$PATH
  echo 'PATH='"${TOOLS_INSTALL_DIR}"'/:'"${TMP_PATH}"'' >> /root/.ssh/environment

}


install_tools
