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
## :
## QUALITY CONTROL:
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
## :
## BUILD:
# ==================================================================================== #

## build: build the cmd/web application
.PHONY: build
build:
	@(cd cmd/gophetch && go get -u && go mod verify && go build -v -ldflags='-s' -o=${PROJECT_ROOT}/bin/gophetch)

# ==================================================================================== #
## :
## VERSIONING:
# ==================================================================================== #

## tag: create a new git tag version
.PHONY: tag
# Makefile command for version tagging with prompt
tag:
	@if [ -z "$$version" ]; then \
		echo "Error: version is not set. Use 'make tag version=x.y.z message=\"Some message\"'"; \
		exit 1; \
	fi
	@if [ -z "$$message" ]; then \
		echo "Error: message is not set. Use 'make tag version=x.y.z message=\"Some message\"'"; \
		exit 1; \
	fi
	@if ! echo "$$version" | egrep -q '^(v[0-9]+\.[0-9]+\.[0-9]+(-alpha|-beta(\.[0-9]+)?)?)$$'; then \
		echo "Error: Invalid version format. It should be like v1.2.3, v1.2.3-alpha, or v1.2.3-beta.1"; \
		exit 1; \
	fi
	@if ! echo "$$version" | egrep -q '^v'; then \
		version="v$$version"; \
	fi
	@latest_tag=`git describe --tags --abbrev=0 2>/dev/null` || echo "No existing tags"; \
	echo "Current version: $$latest_tag"; \
	echo "   About to tag: $$version"; \
	read -p "Do you want to continue? [y/N]: " yn; \
	case $$yn in \
		 [Yy]* ) git tag -a $$version -m "$$message"; git push origin $$version;; \
		 * ) echo "Tagging cancelled."; exit 1;; \
	esac

# ==================================================================================== #
## :
## GENERATORS:
# ==================================================================================== #

## serve-docs: serve the godoc documentation on localhost:6060
.PHONY: serve-docs
serve-docs:
	@godoc -http=:6060
