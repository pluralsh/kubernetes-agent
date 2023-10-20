# Dockerfile for kas

FROM busybox:uclibc as busybox
FROM docker.io/golang:1.21 as builder

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /src
COPY . .

RUN GCFLAGS="all=-N -l" TARGET_DIRECTORY=/agentk make agentk

FROM gcr.io/distroless/base-debian12:nonroot
LABEL source="https://github.com/pluralsh/kas" \
      name="Kubernetes Agent" \
      maintainer="Plural::sre" \
      vendor="Plural" \
      summary="Kubernetes Agent" \
      description="Kubernetes Agent supercharges your Plural CD"

ENV KAS_FLAGS=""

# Copy the static shell into base image.
COPY --from=busybox /bin/sh /bin/sh

COPY --from=builder /go/bin/dlv /
COPY --from=builder /agentk /app
COPY --from=builder /src/build/entrypoint.debug.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
