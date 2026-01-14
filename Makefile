GOPATH         			:= ${HOME}/go
PATH           			:= ${GOPATH}/bin:$(PATH)
SHELL          			:= /bin/bash



.PHONY: all
all: help

.PHONY: fmt
fmt: tidy  ## Run go fmt on all go files
	@go install github.com/daixiang0/gci@latest
	@gci write \
	    -s standard \
	    -s default \
	    -s "prefix(google.golang.org)" \
	    -s "prefix(golang.org)" \
	    -s "prefix(github.com/altessa-s)" \
	    -s "prefix(github.com/altessa-s/ozzo-rules)" \
	    -s blank -s alias \
	 $$(go list -f {{.Dir}} ./...) \

.PHONY: lint
lint: tidy fmt ## Run linter
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

.PHONY: test
test: ## Run all tests
	@go test ./...

.PHONY: tidy
tidy: ## Run go mod tidy
	@go mod tidy

.PHONY: copyright
copyright: ## Add copyright header to all files
	@echo "Updating copyright header to all go files"
	go-copyright-checker check --fix .

.PHONY: precommit-install
precommit-install: ## Install pre-commit hooks
	@pip3 install pre-commit
	@pre-commit install

.PHONY: precommit-run
precommit-run: ## Run pre-commit hooks
	@pre-commit run --all-files

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
