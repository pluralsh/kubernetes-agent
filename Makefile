ROOT_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

include $(ROOT_DIRECTORY)/hack/include/config.mk
include $(TOOLS_MAKEFILE)

.PHONY: help
help:
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ============================ GLOBAL ============================ #
#
# A global list of targets executed for every module, it includes:
# - all modules in 'modules' directory except 'modules/common'
# - all common modules in 'modules/common' directory except 'modules/tools'
#
# ================================================================ #

.PHONY: build
build: clean ## Builds all modules
	@$(MAKE) --no-print-directory -C $(MODULES_DIR) TARGET=build

.PHONY: check
check: clean ## Runs all available checks
	@$(MAKE) --no-print-directory -C $(MODULES_DIR) TARGET=check

.PHONY: clean
clean: --clean ## Clean up all temporary directories
	@$(MAKE) --no-print-directory -C $(MODULES_DIR) TARGET=clean

.PHONY: fix
fix: clean ## Runs all available fix scripts
	@$(MAKE) --no-print-directory -C $(MODULES_DIR) TARGET=fix

.PHONY: test
test: clean ## Runs all available test scripts
	@$(MAKE) --no-print-directory -C $(MODULES_DIR) TARGET=test

# ============================ Local ============================ #

.PHONY: schema
schema: clean
	@echo "[root] Regenerating schemas"
	@(cd $(API_DIR) && make --no-print-directory schema)
	@echo "[root] Schema regenerated successfully"

.PHONY: tools
tools: clean ## Installs required tools

# Starts development version of the application.
#
# URL: http://localhost:8080
#
# Note: Make sure that the port 8080 (Web HTTP) is free on your localhost
.PHONY: serve
serve: clean --ensure-kind-cluster --ensure-metrics-server ## Starts development version of the application on http://localhost:8080
	@KUBECONFIG=$(KIND_CLUSTER_INTERNAL_KUBECONFIG_PATH) \
	SYSTEM_BANNER=$(SYSTEM_BANNER) \
	SYSTEM_BANNER_SEVERITY=$(SYSTEM_BANNER_SEVERITY) \
	SIDECAR_HOST=$(SIDECAR_HOST) \
	docker compose -f $(DOCKER_COMPOSE_DEV_PATH) --project-name=$(PROJECT_NAME) up \
		--build \
		--force-recreate \
		--remove-orphans \
		--no-attach gateway \
		--no-attach scraper \
		--no-attach metrics-server

# Starts production version of the application.
#
# HTTPS: https://localhost:8443
# HTTP: http://localhost:8080
#
# Note: Make sure that the ports 8443 (Gateway HTTPS) and 8080 (Gateway HTTP) are free on your localhost
.PHONY: run
run: clean --ensure-kind-cluster --ensure-metrics-server ## Starts production version of the application on https://localhost:8443 and https://localhost:8000
	@KUBECONFIG=$(KIND_CLUSTER_INTERNAL_KUBECONFIG_PATH) \
	SYSTEM_BANNER=$(SYSTEM_BANNER) \
	SYSTEM_BANNER_SEVERITY=$(SYSTEM_BANNER_SEVERITY) \
	SIDECAR_HOST=$(SIDECAR_HOST) \
	VERSION="v0.0.0-prod" \
	WEB_BUILDER_ARCH=$(ARCH) \
	docker compose -f $(DOCKER_COMPOSE_PATH) --project-name=$(PROJECT_NAME) up \
		--build \
		--remove-orphans \
		--no-attach gateway \
		--no-attach scraper \
		--no-attach metrics-server

.PHONY: image
image:
ifndef NO_BUILD
		@KUBECONFIG=$(KIND_CLUSTER_INTERNAL_KUBECONFIG_PATH) \
		SYSTEM_BANNER=$(SYSTEM_BANNER) \
		SYSTEM_BANNER_SEVERITY=$(SYSTEM_BANNER_SEVERITY) \
		SIDECAR_HOST=$(SIDECAR_HOST) \
		VERSION="v0.0.0-prod" \
		WEB_BUILDER_ARCH=$(ARCH) \
		docker compose -f $(DOCKER_COMPOSE_PATH) --project-name=$(PROJECT_NAME) build \
		--no-cache
endif

# ============================ Private ============================ #

.PHONY: --clean
--clean:
	@echo "[root] Cleaning up"
	@rm -rf $(TMP_DIR)
