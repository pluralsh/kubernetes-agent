FROM debian:10-slim

RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates curl python3 clang git make patch \
  && rm -rf "/var/lib/apt/lists/*" \
  && ln -sf /usr/bin/python3 /usr/local/bin/python

ENV CC="/usr/bin/clang"

ARG GO_VERSION=1.17
ARG BAZELISK_VERSION=v1.10.1
ARG BUILD_DIR=/tmp/build
ARG GO_URL="https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN mkdir ${BUILD_DIR} \
  && cd ${BUILD_DIR} \
  && curl --retry 6 -so golang.tar.gz ${GO_URL} \
  && tar -xf golang.tar.gz \
  && mv ${BUILD_DIR}/go /usr/local/go \
  && rm golang.tar.gz \
  && ln -sf /usr/local/go/bin/go /usr/local/bin/ \
  && rm -rf ${BUILD_DIR}

RUN go install github.com/bazelbuild/bazelisk@${BAZELISK_VERSION} \
    && rm -rf /root/go/pkg \
    && ln -sf /root/go/bin/bazelisk /usr/local/bin/bazel
