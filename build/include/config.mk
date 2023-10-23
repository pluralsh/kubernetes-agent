### Common application/container details
PROJECT_NAME := kubernetes-agent

### Dirs and paths
# Base paths
PARTIALS_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
# Docker files
DOCKER_DIRECTORY := $(ROOT_DIRECTORY)/build/docker
DOCKER_COMPOSE_PATH := $(DOCKER_DIRECTORY)/compose.yaml
DOCKER_COMPOSE_DEBUG_PATH := $(DOCKER_DIRECTORY)/compose.debug.yaml
# Build
DIST_DIRECTORY := $(ROOT_DIRECTORY)/bin
# Secret dir for run targets
SECRET_DIRECTORY := $(ROOT_DIRECTORY)/.secret

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

.PHONY: --certificate
--certificate: --secrets --kubeconfig
	@openssl req -x509 -newkey rsa:4096 -keyout $(SECRET_DIRECTORY)/cert.key -out $(SECRET_DIRECTORY)/cert.pub -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname" 2>/dev/null
	@echo Changing permissions for generated certs and secrets
	@sudo chmod -R +rx $(SECRET_DIRECTORY)
	@sudo chown -R 65532:65532 $(SECRET_DIRECTORY)

.PHONY: --kubeconfig
--kubeconfig:
	@kind export kubeconfig --kubeconfig $(SECRET_DIRECTORY)/kubeconfig
	@sed -i 's/127.0.0.1/host.docker.internal/' $(SECRET_DIRECTORY)/kubeconfig
	@sed -i 's/.*certificate-authority-data.*/    insecure-skip-tls-verify: true/' $(SECRET_DIRECTORY)/kubeconfig

.PHONY: --secrets
--secrets:
	@head -c 512 /dev/urandom | LC_CTYPE=C tr -cd 'a-zA-Z0-9' | head -c 32 | base64 > $(SECRET_DIRECTORY)/api_listen_secret
	@head -c 512 /dev/urandom | LC_CTYPE=C tr -cd 'a-zA-Z0-9' | head -c 32 | base64 > $(SECRET_DIRECTORY)/private_api_secret
	@head -c 512 /dev/urandom | LC_CTYPE=C tr -cd 'a-zA-Z0-9' | head -c 32 | base64 > $(SECRET_DIRECTORY)/redis_server_secret

### GOPATH check
ifndef GOPATH
$(error $$GOPATH environment variable not set)
endif

ifeq (,$(findstring $(GOPATH)/bin,$(PATH)))
$(error $$GOPATH/bin directory is not in your $$PATH)
endif