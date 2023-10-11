include .env

PROJECT_ROOT := $(shell pwd)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	@(go fmt ./... && go mod tidy -v)
	@(cd cmd/gophetch && go get -u && go fmt ./... && go mod tidy -v)

## audit: run quality control checks
.PHONY: audit
audit:
	@(go vet ./... && \
		go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./... && \
		go test -race -vet=off ./... && \
		go mod verify)

## lint: run linter
.PHONY: lint
lint:
	golangci-lint run ./...

## test: run the go tests
## : (use `make test pkg=<path-to-package>` to run a specific package, including integrations)
.PHONY: test
test:
	@if [ -z ${pkg} ]; then \
		PROJECT_ROOT=${PROJECT_ROOT} go test -coverprofile=cover.out -short ./...; \
	else \
		PROJECT_ROOT=${PROJECT_ROOT} go test -coverprofile=cover.out ${pkg}; \
	fi

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cmd/web application
.PHONY: build
build:
	@(cd cmd/gophetch && go get -u && go mod verify && go build -v -ldflags='-s' -o=${PROJECT_ROOT}/bin/gophetch)

## serve-docs: serve the docs on localhost:6060
.PHONY: serve-docs
serve-docs:
	godoc -http=:6060
