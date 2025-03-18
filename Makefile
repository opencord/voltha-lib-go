# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2016-2024 Open Networking Foundation (ONF) and the ONF Contributors
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
$(if $(findstring disabled-joey,$(USER)),\
   $(eval USE_LF_MK := 1)) # special snowflake

##--------------------##
##---]  INCLUDES  [---##
##--------------------##
ifdef USE_LF_MK
  $(error should not be here)
  include lf/include.mk
else
  include lf/transition.mk
endif # ifdef USE_LF_MK

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
	$(call banner-enter,Target $@)
	${GO} build -mod=vendor ./...
	$(call banner-leave,Target $@)

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
mod-update: go-version mod-tidy mod-vendor

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: go-version
go-version :
	$(call banner-enter,Target $@)
	${GO} version
	$(call banner-leave,Target $@)

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
	$(if $(LOCAL_FIX_PERMS),chmod o+w $(CURDIR))
	${GO} mod vendor
	$(if $(LOCAL_FIX_PERMS),chmod o-w $(CURDIR))
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
	@${GOLANGCI_LINT} run --timeout=4m --out-format junit-xml ./... | tee ./sca-report/sca-report.xml
	@echo ""
	@echo "Static code analysis OK"

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
test: local-protos

	$(call banner-enter,Target $@)
	@mkdir -p ./tests/results

        # No stream redirects, exit with shell status
	$(MAKE) test-go

        # Redirect I/O, ignore shell exit status (for now)
	$(MAKE) test-go-cover

	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: test-go
test-go :

	$(call banner-enter,Target $@)
	@echo "** Testing attempt #1: exit-with-error-status: enabled"
	-$(GO) test -mod=vendor ./...
	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: test-go-cover
test-go-cover : gen-coverage-coverprofile gen-coverage-junit gen-coverage-cobertura

## -----------------------------------------------------------------------
## Intent: Generate coverprofile data
## -----------------------------------------------------------------------
cover-dir      := ./tests/results
go-cover-out   := $(cover-dir)/go-test-coverage.out
go-result-out := $(cover-dir)/go-test-results.out

.PHONY: gen-coverage-coverprofile
gen-coverage-coverprofile:

	$(call banner-enter,Target $@)
	@echo "** Testing attempt #2: exit-with-error-status: disabled"

        # Fix docker volume perms if building locally
	touch "$(go-cover-out)"
	$(if $(LOCAL_FIX_PERMS),chmod o+w "$(go-cover-out)")

        # Fix docker volume perms if building locally
	$(if $(LOCAL_FIX_PERMS),touch "$(go-result-out)")
	$(if $(LOCAL_FIX_PERMS),chmod o+w "$(go-result-out)")

        # ------------------------------------------
        # set -euo pipefail else tee masks $? return
        # ------------------------------------------
	@echo '** Running test coverage: exit-on-error is currently disabled'
	-(\
	    set -euo pipefail; \
	    $(GO) test -mod=vendor -v -coverprofile "$(go-cover-out)" -covermode count ./... 2>&1 | tee "$(go-result-out)" \
	)

	$(if $(LOCAL_FIX_PERMS),chmod o-w "$(go-result-out)")
	$(if $(LOCAL_FIX_PERMS),chmod o-w "$(go-cover-out)")

	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## Intent: Morph coverage data into junit/xml content
## -----------------------------------------------------------------------
go-results-xml  := $(cover-dir)/go-test-results.xml

.PHONY: gen-coverage-junit
gen-coverage-junit : gen-coverage-coverprofile
gen-coverage-junit:
	$(call banner-enter,Target $@)

        # Fix docker volume perms if building locally
	$(if $(LOCAL_FIX_PERMS),touch "$(go-results-xml)")
	$(if $(LOCAL_FIX_PERMS),chmod o+w "$(go-results-xml)")

	${GO_JUNIT_REPORT} < $(go-result-out) > "$(go-results-xml)"

	$(if $(LOCAL_FIX_PERMS),chmod o-w "$(go-results-xml)")
	$(call banner-leave,Target $@)

## -----------------------------------------------------------------------
## Intent: Morph coverage data into cobertura xml
## -----------------------------------------------------------------------
go-cover-xml := $(cover-dir)/go-test-coverage.xml

.PHONY: gen-coverage-cobertura
gen-coverage-cobertura : gen-coverage-junit
gen-coverage-cobertura :

	$(call banner-enter,Target $@)

        # Fix docker volume perms if building locally
	$(if $(LOCAL_FIX_PERMS),touch "$(go-cover-xml)")
	$(if $(LOCAL_FIX_PERMS),chmod o+w "$(go-cover-xml)")

	${GOCOVER_COBERTURA} < "$(go-cover-out)" > "$(go-cover-xml)"

	$(if $(LOCAL_FIX_PERMS),chmod o-w "$(go-cover-xml)")
	$(call banner-leave,Target $@)

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
	@echo '[REPORT: coverage]'
	@echo '  gen-coverage-coverprofile    Generate profiling data'
	@echo '  gen-coverage-junit           Generate junit coverage report'
	@echo '  gen-coverage-cobertura       Generate cobertura report'

# [EOF]
