# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2016-2023 Open Networking Foundation (ONF) and the ONF Contributors
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
# -----------------------------------------------------------------------

.DEFAULT_GOAL := help

##-------------------##
##---]  GLOBALS  [---##
##-------------------##

##--------------------##
##---]  INCLUDES  [---##
##--------------------##
include config.mk
include makefiles/include.mk

# Variables
VERSION                    ?= $(shell cat ./VERSION)

## -----------------------------------------------------------------------
## Local Development Helpers
## -----------------------------------------------------------------------
.PHONY: local-protos
local-protos:
	@mkdir -p python/local_imports
ifdef LOCAL_PROTOS
	$(RM) -r vendor/github.com/opencord/voltha-protos
	mkdir -p vendor/github.com/opencord/voltha-protos/v5/go
	cp -r ${LOCAL_PROTOS}/go/* vendor/github.com/opencord/voltha-protos/v5/go
	$(RM) -r python/local_imports/voltha-protos
	mkdir -p python/local_imports/voltha-protos/dist
	cp ${LOCAL_PROTOS}/dist/*.tar.gz python/local_imports/voltha-protos/dist/
endif

## -----------------------------------------------------------------------
## build the library
## -----------------------------------------------------------------------
build: local-protos
	${GO} build -mod=vendor ./...

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
lint-mod:
	$(call banner-entry,Target $@)
	@echo "Running dependency check..."
	@${GO} mod verify
	@echo "Dependency check OK. Running vendor check..."
	@git status > /dev/null
	@git diff-index --quiet HEAD -- go.mod go.sum vendor || (echo "ERROR: Staged or modified files must be committed before running this test" && git status -- go.mod go.sum vendor && exit 1)
	@[[ `git ls-files --exclude-standard --others go.mod go.sum vendor` == "" ]] || (echo "ERROR: Untracked files must be cleaned up before running this test" && git status -- go.mod go.sum vendor && exit 1)

	$(HIDE)$(MAKE) --no-print-directory mod-update

	@git status > /dev/null
	@git diff-index --quiet HEAD -- go.mod go.sum vendor || (echo "ERROR: Modified files detected after running go mod tidy / go mod vendor" && git status -- go.mod go.sum vendor && git checkout -- go.mod go.sum vendor && exit 1)
	@[[ `git ls-files --exclude-standard --others go.mod go.sum vendor` == "" ]] || (echo "ERROR: Untracked files detected after running go mod tidy / go mod vendor" && git status -- go.mod go.sum vendor && git checkout -- go.mod go.sum vendor && exit 1)
	@echo "Vendor check OK."
	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: mod-update
mod-update: mod-tidy mod-vendor

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: mod-tidy
mod-tidy :
	$(call banner-enter,Target $@)
	${GO} mod tidy
	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: mod-vendor
mod-vendor : mod-tidy
mod-vendor :
	$(call banner-enter,Target $@)
	$(if $(LOCAL_FIX_PERMS),chmod 777 $(CURDIR))
	${GO} mod vendor
	$(if $(LOCAL_FIX_PERMS),chmod 755 $(CURDIR))
	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
lint: lint-mod

## -----------------------------------------------------------------------
## Coverage report: Static code analysis
## -----------------------------------------------------------------------
sca:
	@$(RM) -r ./sca-report
	@mkdir -p ./sca-report
	@echo "Running static code analysis..."
	@${GOLANGCI_LINT} run --deadline=4m --out-format junit-xml ./... | tee ./sca-report/sca-report.xml
	@echo ""
	@echo "Static code analysis OK"

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
test: local-protos
	@mkdir -p ./tests/results
	$(if $(LOCAL_FIX_PERMS),chmod 777 tests/results)
	/bin/ls -ld tests/results

	@${GO} test -mod=vendor -v -coverprofile ./tests/results/go-test-coverage.out -covermode count ./... 2>&1 | tee ./tests/results/go-test-results.out ;\
	RETURN=$$? ;\
	${GO_JUNIT_REPORT} < ./tests/results/go-test-results.out > ./tests/results/go-test-results.xml ;\
	${GOCOVER_COBERTURA} < ./tests/results/go-test-coverage.out > ./tests/results/go-test-coverage.xml ;\
	exit $$RETURN

	$(if $(LOCAL_FIX_PERMS),chmod -R 755 tests/results/*)

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
distclean:
	$(RM) -r ./sca-report ./tests

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
clean :: distclean

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
sterile :: clean

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
help ::
	@echo "Usage: make [<target>]"
	@echo "where available targets are:"
	@echo
	@echo "build                : Build the library"
	@echo "clean                : Remove files created by the build"
	@echo "distclean            : Remove build and testing artifacts and reports"
	@echo "lint-mod             : Verify the integrity of the 'mod' files"
	@echo "mod-update           : Update go.mod and the vendor directory"
	@echo "test                 : Generate reports for all go tests"
	@echo

# [EOF]

