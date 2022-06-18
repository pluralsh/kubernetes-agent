# The same list of go build tags must be in four places:
# - Makefile
# - Workspace
# - .bazelrc
# - .golangci.yml
GO_BUILD_TAGS := tracer_static,tracer_static_jaeger

SHELL = /usr/bin/env bash -eo pipefail

GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git tag --points-at HEAD 2>/dev/null || true)
BUILD_TIME = $(shell date -u +%Y%m%d.%H%M%S)
ifeq ($(GIT_TAG), )
	GIT_TAG = "v0.0.0"
endif

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
	go generate -x -v \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd/kas/kasapp" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/starboard_vulnerability/agent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/redistool" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_agent_tracker" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_gitaly" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_gitlab_access" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_internalgitaly" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_kubernetes_api" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modagent" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modshared" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_redis" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_rpc" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_reverse_tunnel_tracker" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_rpc" \
		"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_usage_metrics"

.PHONY: regenerate-mocks
regenerate-mocks: internal-regenerate-mocks fmt update-bazel

.PHONY: update-repos
update-repos:
	go mod tidy -compat=1.18
	bazel run \
		//:gazelle -- \
		update-repos \
		-from_file=go.mod \
		-prune=true \
		-build_file_proto_mode=disable_global \
		-to_macro=build/repositories.bzl%go_repositories
	go mod tidy -compat=1.18

.PHONY: update-bazel
update-bazel: gazelle

.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports -w cmd internal pkg

.PHONY: test
test: fmt update-bazel
	bazel test -- //...

.PHONY: test-ci
test-ci:
	bazel test -- //... //cmd:push-latest

.PHONY: verify-ci
verify-ci: delete-generated-files internal-regenerate-proto internal-regenerate-mocks fmt update-bazel update-repos
	git add .
	git diff --cached --quiet ':(exclude).bazelrc' || (echo "Error: uncommitted changes detected:" && git --no-pager diff --cached && exit 1)

.PHONY: quick-test
quick-test:
	bazel test \
		--build_tests_only \
		-- //...

.PHONY: docker
docker: update-bazel
	bazel build \
		//cmd/agentk:container \
		//cmd/kas:container

# This only works from a linux machine
.PHONY: docker-race
docker-race: update-bazel
	bazel build \
		//cmd/agentk:container_race \
		//cmd/kas:container_race

# Export docker image into local Docker
.PHONY: docker-export
docker-export: update-bazel
	bazel run \
		//cmd/agentk:container \
		-- \
		--norun
	bazel run \
		//cmd/kas:container \
		-- \
		--norun

# Export docker image into local Docker
# This only works on a linux machine
.PHONY: docker-export-race
docker-export-race: update-bazel
	bazel run \
		//cmd/agentk:container_race \
		-- \
		--norun
	bazel run \
		//cmd/kas:container_race \
		-- \
		--norun

# Build and push all docker images tagged as "latest".
# This only works on a linux machine
.PHONY: release-latest
release-latest:
	bazel run //cmd:push-latest
	build/push_multiarch_agentk.sh 'latest'

# Build and push all docker images tagged with the tag on the current commit and as "stable".
# This only works on a linux machine
.PHONY: release-tag-and-stable
release-tag-and-stable:
	bazel run //cmd:push-tag-and-stable
	build/push_multiarch_agentk.sh 'stable'
	build/push_multiarch_agentk.sh "$(GIT_TAG)"

# Build and push all docker images tagged with the tag on the current commit.
# This only works on a linux machine
.PHONY: release-tag
release-tag:
	bazel run //cmd:push-tag
	build/push_multiarch_agentk.sh "$(GIT_TAG)"

# Build and push all docker images tagged with as the current commit.
# This only works on a linux machine
.PHONY: release-commit
release-commit:
	bazel run //cmd:push-commit
	build/push_multiarch_agentk.sh "$(GIT_COMMIT)"

# Set TARGET_DIRECTORY variable to the target directory before running this target
.PHONY: gdk-install
gdk-install:
	bazel run //build:extract_race_binaries_for_gdk -- "$(TARGET_DIRECTORY)"

# Set TARGET_DIRECTORY variable to the target directory before running this target
.PHONY: kas
kas:
	go build \
		-tags "${GO_BUILD_TAGS}" \
		-ldflags "-X gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd.Version=$(GIT_TAG) -X gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd.Commit=$(GIT_COMMIT) -X gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd.BuildTime=$(BUILD_TIME)" \
		-o "$(TARGET_DIRECTORY)" ./cmd/kas

# https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies
.PHONY: show-go-dependency-updates
show-go-dependency-updates:
	go list \
		-tags "${GO_BUILD_TAGS}" \
		-u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all

.PHONY: delete-generated-files
delete-generated-files:
	find . -name '*.pb.go' -type f -delete
	find . -name '*.pb.validate.go' -type f -delete
	find . -name '*_pb.rb' -type f -delete
