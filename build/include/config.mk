### Common application/container details
PROJECT_NAME := kubernetes-agent
# Supported architectures
ARCHITECTURES := linux/amd64 linux/arm64 linux/arm linux/ppc64le linux/s390x # darwin/amd64 darwin/arm64 <- TODO: enable once it is natively supported by docker
BUILDX_ARCHITECTURES := linux/amd64,linux/arm64,linux/arm,linux/ppc64le,linux/s390x # ,darwin/amd64,darwin/arm64

### Dirs and paths
# Base paths
PARTIALS_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
# Docker files
DOCKER_DIRECTORY := $(ROOT_DIRECTORY)/build/docker
DOCKER_COMPOSE_PATH := $(DOCKER_DIRECTORY)/docker.compose.yaml
# Build
DIST_DIRECTORY := $(ROOT_DIRECTORY)/bin

# git invocations must be conditional because git is not available in e.g. CNG and variables are supplied manually.
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
GIT_TAG ?= $(shell git tag --points-at HEAD 2>/dev/null || true)
BUILD_TIME = $(shell date -u +%Y%m%d.%H%M%S)
ifeq ($(GIT_TAG), )
	GIT_TAG = v0.0.0
endif

LDFLAGS := -X "github.com/pluralsh/kuberentes-agent/cmd.Version=$(GIT_TAG)"
LDFLAGS += -X "github.com/pluralsh/kuberentes-agent/cmd.Commit=$(GIT_COMMIT)"
LDFLAGS += -X "github.com/pluralsh/kuberentes-agent/cmd.BuildTime=$(BUILD_TIME)"

### GOPATH check
ifndef GOPATH
$(error $$GOPATH environment variable not set)
endif

ifeq (,$(findstring $(GOPATH)/bin,$(PATH)))
$(error $$GOPATH/bin directory is not in your $$PATH)
endif