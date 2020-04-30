.DEFAULT_GOAL	:= build

SHELL := /bin/bash
PROGNAME := itacho
PKG := github.com/cookpad/itacho
VERSION := $(shell cat VERSION)

.PHONY: build
build:
	@echo "--> building ${PROGNAME} $(VERSION)"
	@go build -ldflags "-X main.version=$(VERSION)" -o "$(PROGNAME)"

.PHONY: lint
lint:
	@echo "--> running linters"
	@go vet

.PHONY: test
test:
	@echo "--> running tests"
	@go test ./...

.PHONY: build_protos
build_protos:
	@echo "--> building protos"
	@./build_protos

.PHONY: gen_legacy
gen_legacy:
	@echo "--> generating xDS response with legacy SDS"
	@rm -rf test/srv
	@mkdir test/srv
	@./itacho generate -s example/book.jsonnet -o test/srv -t CDS -v "$(shell git rev-parse HEAD)" \
	  --eds-cluster sds --use-legacy-sds
	@./itacho generate -s example/book.jsonnet -o test/srv -t RDS -v "$(shell git rev-parse HEAD)"

.PHONY: integration_test
.ONESHEL: integration_test
integration_test: build gen_legacy
	@echo "--> running integration_test"
	@pushd test && \
	  docker-compose down && \
	  docker-compose up --build -d && \
	  ruby test.rb && \
	  docker-compose down && \
	  popd
