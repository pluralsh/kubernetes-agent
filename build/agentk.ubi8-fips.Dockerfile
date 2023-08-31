# Dockerfile for agentk

ARG BUILDER_IMAGE
ARG UBI_IMAGE=registry.access.redhat.com/ubi8/ubi-micro:8.8
ARG UID=1000

FROM ${BUILDER_IMAGE} as builder

WORKDIR /src

COPY . /src

RUN TARGET_DIRECTORY=. make agentk

FROM ${UBI_IMAGE}

LABEL source="https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent" \
      name="GitLab Agent for Kubernetes" \
      maintainer="GitLab group::environments" \
      vendor="GitLab" \
      summary="GitLab Agent for Kubernetes" \
      description="GitLab Agent for Kubernetes allows to integrate your cluster with GitLab in a secure way"

USER ${UID}

COPY --from=builder /src/agentk /usr/bin/agentk

ENTRYPOINT ["/usr/bin/agentk"]
