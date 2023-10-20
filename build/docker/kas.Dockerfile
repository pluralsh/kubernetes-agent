# Dockerfile for kas

FROM docker.io/golang:1.21 as builder

WORKDIR /src
COPY . .

RUN make TARGET_DIRECTORY=. build-kas

FROM gcr.io/distroless/static-debian12:nonroot
LABEL source="https://github.com/pluralsh/kas" \
      name="Kubernetes Agent Server" \
      maintainer="Plural::sre" \
      vendor="Plural" \
      summary="Kubernetes Agent Server" \
      description="Kubernetes Agent Server supercharges your Plural CD"

COPY --from=builder /src/kas /

ENTRYPOINT ["/kas"]
