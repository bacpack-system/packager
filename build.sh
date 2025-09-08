#!/usr/bin/env bash

set -e

VERSION=$(sed -E -n 's/version=([^=]+)/\1/p' < version.txt)
MACHINE=$(uname -m | sed -E 's/_/-/')

INSTALL_DIR="./bringauto-packager_${VERSION}_${MACHINE}-linux"

if [[ -d ${INSTALL_DIR} ]]; then
  echo "${INSTALL_DIR} already exist. Delete it pls" >&2
  exit 1
fi

go get bringauto/cmd/bap-builder

pushd cmd/bap-builder
  echo "Compile bap_builder"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w'
popd

mkdir -p "${INSTALL_DIR}"

cp cmd/bap-builder/bap-builder             "${INSTALL_DIR}/"
cp -r doc                                  "${INSTALL_DIR}/"
cp -r example_context                      "${INSTALL_DIR}/"
cp README.md                               "${INSTALL_DIR}/"
cp LICENSE                                 "${INSTALL_DIR}/"


zip -r "bringauto-packager_v${VERSION}_${MACHINE}-linux.zip" ${INSTALL_DIR}/

rm -fr "${INSTALL_DIR}"
