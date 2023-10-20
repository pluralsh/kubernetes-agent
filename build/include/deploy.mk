.PHONY: --image-debug
--image-debug: APP_VERSION = "debug"
--image-debug: --ensure-variables-set --image

.PHONY: --image
--image: APP_VERSION = "latest"
--image: --ensure-variables-set
	echo "Building '$(APP_NAME):$(APP_VERSION)'" ; \
	docker build \
		-f $(DOCKERFILE) \
		-t $(APP_NAME):$(APP_VERSION) \
		$(ROOT_DIRECTORY) ; \

.PHONY: --ensure-variables-set
--ensure-variables-set:
	@if [ -z "$(DOCKERFILE)" ]; then \
  	echo "DOCKERFILE variable not set" ; \
  	exit 1 ; \
  fi ; \
	if [ -z "$(APP_NAME)" ]; then \
		echo "APP_NAME variable not set" ; \
		exit 1 ; \
	fi ; \