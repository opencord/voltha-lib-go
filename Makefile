#
# Copyright 2016 the original author or authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# set default shell
SHELL = bash -e -o pipefail

# Variables
VERSION                    ?= $(shell cat ./VERSION)

# Default is GO111MODULE=auto, which will refuse to use go mod if running
# go less than 1.13.0 and this repository is checked out in GOPATH. For now,
# force module usage. This affects commands executed from this Makefile, but
# not the environment inside the Docker build (which does not build from
# inside a GOPATH).
export GO111MODULE=on

.PHONY: local-protos

# This should to be the first and default target in this Makefile
help:
	@echo "Usage: make [<target>]"
	@echo "where available targets are:"
	@echo
	@echo "build                : Build the library"
	@echo "clean                : Remove files created by the build"
	@echo "distclean            : Remove build and testing artifacts and reports"
	@echo "lint-style           : Verify code is properly gofmt-ed"
	@echo "lint-sanity          : Verify that 'go vet' doesn't report any issues"
	@echo "lint-mod             : Verify the integrity of the 'mod' files"
	@echo "lint                 : Shorthand for lint-style & lint-sanity"
	@echo "mod-update           : Update go.mod and the vendor directory"
	@echo "test                 : Generate reports for all go tests"
	@echo

## Local Development Helpers
local-protos:
	@mkdir -p python/local_imports
ifdef LOCAL_PROTOS
	rm -rf vendor/github.com/opencord/voltha-protos
	mkdir -p vendor/github.com/opencord/voltha-protos/v3/go
	cp -r ${LOCAL_PROTOS}/go/* vendor/github.com/opencord/voltha-protos/v3/go
	rm -rf python/local_imports/voltha-protos
	mkdir -p python/local_imports/voltha-protos/dist
	cp ${LOCAL_PROTOS}/dist/*.tar.gz python/local_imports/voltha-protos/dist/
endif

## build the library
build: local-protos
	go build -mod=vendor ./...

## lint and unit tests

lint-style:
ifeq (,$(shell which gofmt))
	go get -u github.com/golang/go/src/cmd/gofmt
endif
	@echo "Running style check..."
	@gofmt_out="$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))" ;\
	if [ ! -z "$$gofmt_out" ]; then \
	  echo "$$gofmt_out" ;\
	  echo "Style check failed on one or more files ^, run 'go fmt' to fix." ;\
	  exit 1 ;\
	fi
	@echo "Style check OK"

lint-sanity:
	@echo "Running sanity check..."
	@go vet -mod=vendor ./...
	@echo "Sanity check OK"

lint-mod:
	@echo "Running dependency check..."
	@go mod verify
	@echo "Dependency check OK"

lint: lint-style lint-sanity lint-mod

# Rules to automatically install golangci-lint
GOLANGCI_LINT_TOOL?=$(shell which golangci-lint)
ifeq (,$(GOLANGCI_LINT_TOOL))
GOLANGCI_LINT_TOOL=$(GOPATH)/bin/golangci-lint
golangci_lint_tool_install:
	# Same version as installed by Jenkins ci-management
	# Note that install using `go get` is not recommended as per https://github.com/golangci/golangci-lint
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOPATH)/bin v1.17.0
else
golangci_lint_tool_install:
endif

# Rules to automatically install go-junit-report
GO_JUNIT_REPORT?=$(shell which go-junit-report)
ifeq (,$(GO_JUNIT_REPORT))
GO_JUNIT_REPORT=$(GOPATH)/bin/go-junit-report
go_junit_install:
	go get -u github.com/jstemmer/go-junit-report
else
go_junit_install:
endif

# Rules to automatically install gocover-covertura
GOCOVER_COBERTURA?=$(shell which gocover-cobertura)
ifeq (,$(GOCOVER_COBERTURA))
	@GOCOVER_COBERTURA=$(GOPATH)/bin/gocover-cobertura
gocover_cobertura_install:
	go get -u github.com/t-yuki/gocover-cobertura
else
gocover_cobertura_install:
endif

sca: golangci_lint_tool_install
	rm -rf ./sca-report
	@mkdir -p ./sca-report
	$(GOLANGCI_LINT_TOOL) run --out-format junit-xml ./... 2>&1 | tee ./sca-report/sca-report.xml

test: go_junit_install gocover_cobertura_install
	@mkdir -p ./tests/results
	@go test -mod=vendor -v -coverprofile ./tests/results/go-test-coverage.out -covermode count ./... 2>&1 | tee ./tests/results/go-test-results.out ;\
	RETURN=$$? ;\
	$(GO_JUNIT_REPORT) < ./tests/results/go-test-results.out > ./tests/results/go-test-results.xml ;\
	$(GOCOVER_COBERTURA) < ./tests/results/go-test-coverage.out > ./tests/results/go-test-coverage.xml ;\
	exit $$RETURN

clean:

distclean: clean
	rm -rf ./sca_report ./tests

mod-update: build
	go mod tidy
	go mod vendor

# end file
