.DEFAULT_GOAL := help

.PHONY: help init gen

# initialize develop environment
init: ## Prepare dependencies
	@go mod download

# generate gorm code
gen: ## Generate models and query code
	@go run ./internal/cmd/gen

help: ## Show available targets
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_.-]+:.*## / {printf "  %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
