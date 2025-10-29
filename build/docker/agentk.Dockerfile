# Dockerfile for agentk

FROM docker.io/golang:1.25 as builder

WORKDIR /src
COPY . .

RUN make TARGET_DIRECTORY=. build-agentk

FROM gcr.io/distroless/static-debian12:nonroot
LABEL source="https://github.com/pluralsh/kubernetes-agent" \
      name="Kubernetes Agent" \
      maintainer="Plural::sre" \
      vendor="Plural" \
      summary="Kubernetes Agent" \
      description="Kubernetes Agent supercharges your Plural CD"

COPY --from=builder /src/agentk /

ENTRYPOINT ["/agentk"]
