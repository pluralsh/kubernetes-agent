.PHONY: --ensure-tools
--ensure-tools: --ensure-go-tools

.PHONY: --ensure-go-tools
--ensure-go-tools:
	@echo "ensuring required tools availability..."
	@cat internal/tool/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %