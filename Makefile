SHELL := /bin/bash

.PHONY: help clean build fmt lint run test style license godoc install_docgo install_addlicense before_commit install_golangci-lint

SOURCES:=$(shell find . -name '*.go')
DOCFILES:=$(addprefix docs/packages/, $(addsuffix .html, $(basename ${SOURCES})))


install_golangci-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

fmt: install_golangci-lint ## Run go formatting
	@echo "Running go formatting"
	golangci-lint fmt

lint: install_golangci-lint ## Run go liting
	@echo "Running go linting"
	golangci-lint run --fix

shellcheck: ## Run shellcheck
	./shellcheck.sh

abcgo: ## Run ABC metrics checker
	@echo "Run ABC metrics checker"
	./abcgo.sh

style: fmt lint shellcheck abcgo ## Run all the formatting related commands (fmt, lint, abcgo) + check shell scripts

test: clean build ## Run the unit tests
	@go test -coverprofile coverage.out $(shell go list ./...)

benchmark: ## Run benchmarks
	@echo "Running benchmarks"
	./benchmark.sh

cover: test ## Display test coverage on generated HTML pages
	@go tool cover -html=coverage.out

coverage: ## Display test coverage onto terminal
	@go tool cover -func=coverage.out

before_commit: style test license ## Checks done before commit
	./check_coverage.sh

license: install_addlicense  ## Add license in every file in repository
	addlicense -c "Red Hat, Inc" -l "apache" -v ./

docs/packages/%.html: %.go
	mkdir -p $(dir $@)
	docgo -outdir $(dir $@) $^
	addlicense -c "Red Hat, Inc" -l "apache" -v $@

godoc: export GO111MODULE=off
godoc: install_docgo install_addlicense ${DOCFILES} docs/sources.md

docs/sources.md: docs/sources.tmpl.md ${DOCFILES}
	./gen_sources_md.sh

install_docgo: export GO111MODULE=off
install_docgo:
	[[ `command -v docgo` ]] || GO111MODULE=off go get -u github.com/dhconnelly/docgo

install_addlicense: export GO111MODULE=off
install_addlicense:
	[[ `command -v addlicense` ]] || GO111MODULE=off go get -u github.com/google/addlicense

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''
