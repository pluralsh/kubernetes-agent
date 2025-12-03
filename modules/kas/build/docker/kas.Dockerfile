# Dockerfile for kas

FROM docker.io/golang:1.25 as builder

WORKDIR /src
COPY . .

RUN make TARGET_DIRECTORY=/kas -C modules/kas build-kas

FROM gcr.io/distroless/static-debian12:nonroot
LABEL source="https://github.com/pluralsh/kubernetes-agent" \
      name="Kubernetes Agent Server" \
      maintainer="Plural::sre" \
      vendor="Plural" \
      summary="Kubernetes Agent Server" \
      description="Kubernetes Agent Server supercharges your Plural CD"

COPY --from=builder /kas /

ENTRYPOINT ["/kas"]
