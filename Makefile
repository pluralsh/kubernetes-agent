ROOT_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

include $(ROOT_DIRECTORY)/build/include/config.mk
include $(ROOT_DIRECTORY)/build/include/deploy.mk
include $(ROOT_DIRECTORY)/build/include/tools.mk

MAKEFLAGS += -j2

# List of targets that should be executed before other targets
PRE = --ensure

##@ General

.PHONY: help
help: ## show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: go-dep-updates
go-dep-updates: ## show possible go dependency updates
	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}} {{.Version}} -> {{.Update.Version}}{{end}}' -m all

##@ Run

##@ Build

.PHONY: build
build: build-kas build-agentk ## build both kas and agentk

.PHONY: build-kas
build-kas: TARGET_DIRECTORY=bin/kas
build-kas: ## build kas
	CGO_ENABLED=0 go build \
    	-gcflags='$(GCFLAGS)' \
		-ldflags '$(LDFLAGS)' \
		-o $(TARGET_DIRECTORY) ./cmd/kas

.PHONY: build-agentk
build-agentk: TARGET_DIRECTORY=bin/agentk
build-agentk: ## build agentk
	CGO_ENABLED=0 go build \
    	-gcflags='$(GCFLAGS)' \
		-ldflags '$(LDFLAGS)' \
		-o $(TARGET_DIRECTORY) ./cmd/agentk

##@ Docker

.PHONY: docker-kas
docker-kas: APP_NAME=kas
docker-kas: DOCKERFILE=${DOCKER_DIRECTORY}/kas.Dockerfile
docker-kas: --image ## build docker kas image

.PHONY: docker-kas-debug
docker-kas: APP_NAME=kas
docker-kas: DOCKERFILE=${DOCKER_DIRECTORY}/kas.debug.Dockerfile
docker-kas-debug: --image-debug ## build docker kas debug image with embedded delve

.PHONY: docker-agentk
docker-agentk: APP_NAME=agentk
docker-agentk: DOCKERFILE=${DOCKER_DIRECTORY}/agentk.Dockerfile
docker-agentk: --image ## build docker agentk

.PHONY: docker-agentk-debug
docker-agentk-debug: APP_NAME=agentk
docker-agentk-debug: DOCKERFILE=${DOCKER_DIRECTORY}/agentk.debug.Dockerfile
docker-agentk-debug: --image-debug ## build docker agentk debug image with embedded delve

##@ Codegen

.PHONY: codegen
codegen: --mocks --protoc ## regenerate proto and mocks

.PHONY: codegen-delete
codegen-delete: ## delete generated files
	find . -name '*.pb.go' -type f -delete
	find . -name '*.pb.validate.go' -type f -delete
	find . \( -name '*_pb.rb' -and -not -name 'validate_pb.rb' \) -type f -delete
	find . -name '*_proto_docs.md' -type f -delete

.PHONY: --protoc
--protoc:
	@build/protoc.sh

.PHONY: --mocks
--mocks:
	@PATH="${PATH}:$(shell pwd)/build" go generate -x -v \
		"github.com/pluralsh/kuberentes-agent/cmd/agentk/agentkapp" \
		"github.com/pluralsh/kuberentes-agent/cmd/kas/kasapp" \
		"github.com/pluralsh/kuberentes-agent/internal/module/modagent" \
		"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/tunnel" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/redistool" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_agent_registrar" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_agent_tracker" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_cache" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_k8s" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_kubernetes_api" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_modagent" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_modserver" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_modshared" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_redis" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_reverse_tunnel_rpc" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_reverse_tunnel_tunnel" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_rpc" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_stdlib" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_tool" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_usage_metrics"

##@ Tests

.PHONY: test
test: ## run tests
	go test ./... -v

.PHONY: lint
lint: $(PRE) ## run linters
	golangci-lint run ./...

.PHONY: fix
fix: $(PRE) ## fix issues found by linters
	golangci-lint run --fix ./...

################################ Old definitions

# Build and push all FIPS docker images tagged with the tag on the current commit and as "stable".
# This only works on a linux machine
# Set ARCH to the desired CPU architecture.
# Set CI_REGISTRY_IMAGE to the desired container image name.
# Set BUILDER_IMAGE to the image to be used for building the agentk binary.
.PHONY: release-tag-and-stable-fips
release-tag-and-stable-fips:
	docker buildx build --build-arg 'BUILDER_IMAGE=$(BUILDER_IMAGE)' --platform 'linux/$(ARCH)' --file build/agentk.ubi8-fips.Dockerfile --tag '$(FIPS_TAG_GIT_ARCH)' --tag '$(FIPS_TAG_STABLE_ARCH)' .
	docker push '$(FIPS_TAG_GIT_ARCH)'
	docker push '$(FIPS_TAG_STABLE_ARCH)'

.PHONY: release-tag-and-stable-fips-manifest
release-tag-and-stable-fips-manifest:
	docker manifest create '$(FIPS_TAG_GIT)' \
		--amend '$(FIPS_TAG_GIT)-amd64' \
		--amend '$(FIPS_TAG_GIT)-arm64'
	docker manifest push '$(FIPS_TAG_GIT)'
	docker manifest create '$(FIPS_TAG_STABLE)' \
		--amend '$(FIPS_TAG_STABLE)-amd64' \
		--amend '$(FIPS_TAG_STABLE)-arm64'
	docker manifest push '$(FIPS_TAG_STABLE)'

# Build and push all FIPS docker images tagged with the tag on the current commit.
# This only works on a linux machine
# Set ARCH to the desired CPU architecture.
# Set CI_REGISTRY_IMAGE to the desired container image name.
# Set BUILDER_IMAGE to the image to be used for building the agentk binary.
.PHONY: release-tag-fips
release-tag-fips:
	docker buildx build --build-arg 'BUILDER_IMAGE=$(BUILDER_IMAGE)' --platform 'linux/$(ARCH)' --file build/agentk.ubi8-fips.Dockerfile --tag '$(FIPS_TAG_GIT_ARCH)' .
	docker push '$(FIPS_TAG_GIT_ARCH)'

.PHONY: release-tag-fips-manifest
release-tag-fips-manifest:
	docker manifest create '$(FIPS_TAG_GIT)' \
		--amend '$(FIPS_TAG_GIT)-amd64' \
		--amend '$(FIPS_TAG_GIT)-arm64'
	docker manifest push '$(FIPS_TAG_GIT)'

