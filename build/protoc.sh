#!/usr/bin/bash

# Exit on error
set -e

ROOT_DIR="$(cd $(dirname "${BASH_SOURCE}")/.. && pwd -P)"

ARCH=$(uname -m)
case "$ARCH" in
  arm64*)  ARCH="aarch_64" ;;
esac

OS="linux"
case "$OSTYPE" in
  darwin*)  OS="osx" ;;
  linux*)   OS="linux" ;;
  *)        OS="linux" ;;
esac

PROTOC_VERSION="24.4"
PROTOC_ZIP="protoc-${PROTOC_VERSION}-${OS}-${ARCH}.zip"

function plrl::protoc::ensure() {
  if [[ -z "$(which protoc)" || "$(protoc --version)" != *"${PROTOC_VERSION}" ]]; then
    echo "Installing protoc@${PROTOC_VERSION}"
    curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"

    sudo unzip -o "${PROTOC_ZIP}" -d /usr/local bin/protoc
    sudo unzip -o "${PROTOC_ZIP}" -d /usr/local 'include/*'

    rm -f "${PROTOC_ZIP}"
  else
    echo "Found protoc: $(protoc --version)"
  fi
}

function plrl::protoc::generate() {
  local package=${1}
  local files=$(find "${package}" -name "*.proto")

  for proto in ${files}; do
    local baseDir="${proto%/*}"
    local filename="${proto##*/}"

    protoc \
      -I"${ROOT_DIR}" \
      -I"${ROOT_DIR}/build/proto" \
      --proto_path="${baseDir}" \
      --go_out="${ROOT_DIR}" \
      --go_opt=paths=source_relative \
      --go-grpc_out="${ROOT_DIR}" \
      --go-grpc_opt=paths=source_relative \
      --validate_out="${ROOT_DIR}" \
      --validate_opt=paths=source_relative,lang=go \
      --doc_out="${ROOT_DIR}" \
      --doc_opt=markdown,"${filename%.*}_proto_docs.md",source_relative \
      "${baseDir}/${filename}"
  done
}

plrl::protoc::ensure
plrl::protoc::generate "${ROOT_DIR}/cmd"
plrl::protoc::generate "${ROOT_DIR}/pkg"