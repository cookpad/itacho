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

.PHONY: integration_test
.ONESHEL: integration_test
integration_test:
	@echo "--> running integration_test"
	@pushd test && \
	  docker-compose down && \
	  docker-compose up --build -d && \
	  bundle install && \
	  bundle exec ruby test.rb && \
	  docker-compose down && \
	  popd
