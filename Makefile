LDFLAGS := -X "github.com/pluralsh/kuberentes-agent/cmd.Version=$(GIT_TAG)"
LDFLAGS += -X "github.com/pluralsh/kuberentes-agent/cmd.Commit=$(GIT_COMMIT)"
LDFLAGS += -X "github.com/pluralsh/kuberentes-agent/cmd.BuildTime=$(BUILD_TIME)"

include tools.mk

ifndef GOPATH
$(error $$GOPATH environment variable not set)
endif

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

.PHONY: build-kas
build-kas: ## build kas
	CGO_ENABLED=0 go build \
    	-gcflags='$(GCFLAGS)' \
		-ldflags '$(LDFLAGS)' \
		-o bin/kas ./cmd/kas

.PHONY: build-agentk
build-agentk: ## build agentk
	CGO_ENABLED=0 go build \
    	-gcflags='$(GCFLAGS)' \
		-ldflags '$(LDFLAGS)' \
		-o bin/agentk ./cmd/agentk

.PHONY: build-gem
build-gem: ## build kas grpc ruby gem
	cd pkg/ruby && gem build kas-grpc.gemspec

##@ Codegen

.PHONY: delete-codegen
delete-codegen: ## delete generated files
	find . -name '*.pb.go' -type f -delete
	find . -name '*.pb.validate.go' -type f -delete
	find . \( -name '*_pb.rb' -and -not -name 'validate_pb.rb' \) -type f -delete
	find . -name '*_proto_docs.md' -type f -delete

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

SHELL = /usr/bin/env bash -eo pipefail

# git invocations must be conditional because git is not available in e.g. CNG and variables are supplied manually.
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
GIT_TAG ?= $(shell git tag --points-at HEAD 2>/dev/null || true)
BUILD_TIME = $(shell date -u +%Y%m%d.%H%M%S)
ifeq ($(GIT_TAG), )
	GIT_TAG = v0.0.0
endif

FIPS_TAG_GIT = $(CI_REGISTRY_IMAGE)/agentk-fips:$(GIT_TAG)
FIPS_TAG_GIT_ARCH = $(FIPS_TAG_GIT)-$(ARCH)

FIPS_TAG_STABLE = $(CI_REGISTRY_IMAGE)/agentk-fips:stable
FIPS_TAG_STABLE_ARCH = $(FIPS_TAG_STABLE)-$(ARCH)

CI_REGISTRY ?= registry.gitlab.com
CI_PROJECT_PATH ?= gitlab-org/cluster-integration/gitlab-agent
OCI_REPO = $(CI_REGISTRY)/$(CI_PROJECT_PATH)/agentk

.PHONY: internal-regenerate-proto
internal-regenerate-proto:
	# generate go from proto
	#bazel run //build:extract_generated_proto

.PHONY: regenerate-proto
regenerate-proto: internal-regenerate-proto

.PHONY: internal-regenerate-mocks
internal-regenerate-mocks:
	PATH="${PATH}:$(shell pwd)/build" go generate -x -v \
		"github.com/pluralsh/kuberentes-agent/cmd/agentk/agentkapp" \
		"github.com/pluralsh/kuberentes-agent/cmd/kas/kasapp" \
		"github.com/pluralsh/kuberentes-agent/internal/module/modagent" \
		"github.com/pluralsh/kuberentes-agent/internal/module/flux/agent" \
		"github.com/pluralsh/kuberentes-agent/internal/module/flux/rpc" \
		"github.com/pluralsh/kuberentes-agent/internal/module/gitops/agent/manifestops" \
		"github.com/pluralsh/kuberentes-agent/internal/module/google_profiler/agent" \
		"github.com/pluralsh/kuberentes-agent/internal/module/starboard_vulnerability/agent" \
		"github.com/pluralsh/kuberentes-agent/internal/module/reverse_tunnel/tunnel" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/redistool" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_agent_registrar" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_agent_tracker" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_cache" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_gitaly" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_gitlab_access" \
		"github.com/pluralsh/kuberentes-agent/internal/tool/testing/mock_internalgitaly" \
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

.PHONY: regenerate-mocks
regenerate-mocks: internal-regenerate-mocks

.PHONY: update-repos
update-repos:
	go mod tidy
	bazel run \
		//:gazelle -- \
		update-repos \
		-from_file=go.mod \
		-prune=true \
		-build_file_proto_mode=disable_global \
		-to_macro=build/repositories.bzl%go_repositories
	go mod tidy

.PHONY: test-ci
test-ci:
	bazel test -- //... //cmd/agentk:push //cmd/agentk:push_debug //cmd/kas:push //cmd/kas:push_debug

.PHONY: test-ci-fips
test-ci-fips:
	# No -race because it doesn't work on arm64: https://github.com/golang/go/issues/29948
	# FATAL: ThreadSanitizer: unsupported VMA range
	# FATAL: Found 39 - Supported 48
	go test -v ./...

.PHONY: verify-ci
verify-ci: internal-regenerate-proto internal-regenerate-mocks update-repos
	git add .
	git diff --cached --quiet ':(exclude).bazelrc' || (echo "Error: uncommitted changes detected:" && git --no-pager diff --cached && exit 1)

# Build and push all docker images tagged as "latest".
# This only works on a linux machine
# Flags are for oci_push rule. Docs https://docs.aspect.build/rules/rules_oci/docs/push.
.PHONY: release-latest
release-latest:
	bazel run //cmd/agentk:push -- --repository='$(OCI_REPO)' --tag=latest
	bazel run //cmd/agentk:push_debug -- --repository='$(OCI_REPO)' --tag=latest-debug

# Build and push all docker images tagged as "stable".
# This only works on a linux machine
# Flags are for oci_push rule. Docs https://docs.aspect.build/rules/rules_oci/docs/push.
.PHONY: release-stable
release-stable:
	bazel run //cmd/agentk:push -- --repository='$(OCI_REPO)' --tag=stable
	bazel run //cmd/agentk:push_debug -- --repository='$(OCI_REPO)' --tag=stable-debug

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

# Build and push all docker images tagged with the tag on the current commit.
# This only works on a linux machine
# Flags are for oci_push rule. Docs https://docs.aspect.build/rules/rules_oci/docs/push.
.PHONY: release-tag
release-tag:
	bazel run //cmd/agentk:push -- --repository='$(OCI_REPO)' --tag='$(GIT_TAG)'
	bazel run //cmd/agentk:push_debug -- --repository='$(OCI_REPO)' --tag='$(GIT_TAG)-debug'

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

# Build and push all docker images tagged with as the current commit.
# This only works on a linux machine
.PHONY: release-commit
release-commit:
	bazel run //cmd/agentk:push -- --repository='$(OCI_REPO)' --tag='$(GIT_COMMIT)'
	bazel run //cmd/agentk:push_debug -- --repository='$(OCI_REPO)' --tag='$(GIT_COMMIT)-debug'

# Set TARGET_DIRECTORY variable to the target directory before running this target
.PHONY: gdk-install
gdk-install:
	bazel run //cmd/kas:extract_kas_race
	mv 'cmd/kas/kas_race' '$(TARGET_DIRECTORY)'
