# Dockerfile for kas

FROM busybox:uclibc as busybox
FROM docker.io/golang:1.25 as builder

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /src
COPY . .

RUN GCFLAGS="all=-N -l" make TARGET_DIRECTORY=/kas build-kas

FROM gcr.io/distroless/base-debian12:nonroot
LABEL source="https://github.com/pluralsh/kubernetes-agent" \
      name="Kubernetes Agent Server" \
      maintainer="Plural::sre" \
      vendor="Plural" \
      summary="Kubernetes Agent Server" \
      description="Kubernetes Agent Server supercharges your Plural CD"

ENV KAS_FLAGS=""

# Copy the static shell into base image.
COPY --from=busybox /bin/sh /bin/sh

COPY --from=builder /go/bin/dlv /
COPY --from=builder /kas /app
COPY --from=builder /src/build/docker/entrypoint.debug.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
