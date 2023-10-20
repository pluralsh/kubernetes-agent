#!/usr/bin/bash

ROOT_DIR="$(cd $(dirname "${BASH_SOURCE}")/.. && pwd -P)"

ARCH=$(uname -m)
OS="linux"
case "$OSTYPE" in
  darwin*)  OS="osx" ;;
  linux*)   OS="linux" ;;
  *)        OS="linux" ;;
esac

PROTOC_VERSION="24.4"
PROTOC_ZIP="protoc-${PROTOC_VERSION}-${OS}-${ARCH}.zip"
PROTOC_GEN_GO_VERSION="v1.31.0"
PROTOC_GEN_GO_GRPC_VERSION="v1.2.0"
PROTOC_GEN_VALIDATE="v1.0.2"
PROTOC_GEN_DOC="v1.5.1"

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

function plrl::protoc-gen-go::install() {
  if [[ -z "$(which protoc-gen-go)" || "$(protoc-gen-go --version)" != "protoc-gen-go ${PROTOC_GEN_GO_VERSION}" ]]; then
    echo "Installing protoc-gen-go@${PROTOC_GEN_GO_VERSION}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}
  else
    echo "Found protoc-gen-go"
  fi
}

function plrl::protoc-gen-go-grpc::install() {
  if [[ -z "$(which protoc-gen-go-grpc)" ]]; then
    echo "Installing protoc-gen-go-grpc"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}
    go mod tidy
  else
    echo "Found protoc-gen-go-grpc"
  fi
}

function plrl::protoc-gen-validate::install() {
  if [[ -z "$(which protoc-gen-validate)" ]]; then
    echo "Installing protoc-gen-validate"
    go install github.com/envoyproxy/protoc-gen-validate@${PROTOC_GEN_VALIDATE}
    go mod tidy
  else
    echo "Found protoc-gen-validate"
  fi
}

function plrl::protoc-gen-doc::install() {
  if [[ -z "$(which protoc-gen-doc)" ]]; then
    echo "Installing protoc-gen-doc"
    go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@${PROTOC_GEN_DOC}
    go mod tidy
  else
    echo "Found protoc-gen-doc"
  fi
}

function plrl::protoc::generate() {
  local package=${1}
  local files=$(find "${package}" -name "*.proto")

  for proto in ${files}; do
    local baseDir="${proto%/*}"
    local filename="${proto##*/}"

    echo "Generating Go file: ${filename%.*}.pb.go"
    echo "Generating Go GRPC file: ${filename%.*}_grpc.pb.go"
    echo "Generating Go validate file: ${filename%.*}.validate.pb.go"
    echo "Generating Markdown docs file: ${filename%.*}_proto_docs.md"

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
plrl::protoc-gen-go::install
plrl::protoc-gen-go-grpc::install
plrl::protoc-gen-validate::install
plrl::protoc-gen-doc::install
plrl::protoc::generate "${ROOT_DIR}/cmd"
plrl::protoc::generate "${ROOT_DIR}/pkg"
plrl::protoc::generate "${ROOT_DIR}/internal"