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

LDFLAGS := -X "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd.Version=$(GIT_TAG)"
LDFLAGS += -X "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd.Commit=$(GIT_COMMIT)"
LDFLAGS += -X "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd.BuildTime=$(BUILD_TIME)"

CI_REGISTRY ?= registry.gitlab.com
CI_PROJECT_PATH ?= gitlab-org/cluster-integration/gitlab-agent
OCI_REPO = $(CI_REGISTRY)/$(CI_PROJECT_PATH)/agentk

# Install using your package manager, as recommended by
# https://golangci-lint.run/usage/install/#local-installation
.PHONY: lint
lint:
	golangci-lint run

.PHONY: buildozer
buildozer:
	bazel run //:buildozer

.PHONY: buildifier
buildifier:
	bazel run //:buildifier

.PHONY: fmt-bazel
fmt-bazel: gazelle buildozer buildifier

.PHONY: gazelle
gazelle:
	bazel run //:gazelle

.PHONY: internal-regenerate-proto
internal-regenerate-proto:
	bazel run //build:extract_generated_proto

.PHONY: regenerate-proto
regenerate-proto: internal-regenerate-proto fmt update-bazel

.PHONY: internal-regenerate-mocks
internal-regenerate-mocks:
	PATH="${PATH}:$(shell pwd)/build" go generate -x -v \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd/agentk/agentkapp" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd/kas/kasapp" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/agent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/flux/rpc" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/agent/manifestops" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/starboard_vulnerability/agent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_agent_registrar" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_agent_tracker" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_cache" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_gitaly" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_gitlab_access" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_internalgitaly" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_k8s" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_kubernetes_api" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modserver" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modshared" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_redis" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_reverse_tunnel_rpc" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_reverse_tunnel_tunnel" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_rpc" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_stdlib" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_tool" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_usage_metrics"

.PHONY: regenerate-mocks
regenerate-mocks: internal-regenerate-mocks fmt update-bazel

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

.PHONY: update-bazel
update-bazel: gazelle

.PHONY: fmt
fmt:
	go run github.com/daixiang0/gci@v0.11.2 write cmd internal pkg -s standard -s default

.PHONY: test
test: fmt update-bazel
	bazel test -- //...

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
verify-ci: delete-generated-files internal-regenerate-proto internal-regenerate-mocks fmt update-bazel update-repos
	git add .
	git diff --cached --quiet ':(exclude).bazelrc' || (echo "Error: uncommitted changes detected:" && git --no-pager diff --cached && exit 1)

.PHONY: quick-test
quick-test:
	bazel test \
		--build_tests_only \
		-- //...

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

# Set TARGET_DIRECTORY variable to the target directory before running this target
# Optional: set GIT_TAG and GIT_COMMIT variables to supply those values manually.
# This target is used by:
# - CNG: https://gitlab.com/gitlab-org/build/CNG/-/tree/master/gitlab-kas
# - Omnibus: https://gitlab.com/gitlab-org/omnibus-gitlab/-/blob/master/config/software/gitlab-kas.rb
.PHONY: kas
kas:
	go build \
		-ldflags '$(LDFLAGS)' \
		-o '$(TARGET_DIRECTORY)' ./cmd/kas

# Set TARGET_DIRECTORY variable to the target directory before running this target
# Optional: set GIT_TAG and GIT_COMMIT variables to supply those values manually.
# This target is used by FIPS build in this repo.
.PHONY: agentk
agentk:
	go build \
		-ldflags '$(LDFLAGS)' \
		-o '$(TARGET_DIRECTORY)' ./cmd/agentk

# https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies
.PHONY: show-go-dependency-updates
show-go-dependency-updates:
	go list \
		-u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}} {{.Version}} -> {{.Update.Version}}{{end}}' -m all

.PHONY: delete-generated-files
delete-generated-files:
	find . -name '*.pb.go' -type f -delete
	find . -name '*.pb.validate.go' -type f -delete
	find . \( -name '*_pb.rb' -and -not -name 'validate_pb.rb' \) -type f -delete
	find . -name '*_proto_docs.md' -type f -delete

# Build the KAS gRPC ruby gem
.PHONY: build-gem
build-gem:
	cd pkg/ruby && gem build kas-grpc.gemspec
