# The MIT License (MIT)
#
# Copyright (c) 2018, Secure64 Software Corporation
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.


####################################
# Edit these macros.
HELP_DESC              := Builds the gorun golang package.
GO_REPO_SITE           := github.com
GO_REPO_SUBDIR         := apatters/go-run
####################################
# Do not touch anything below here.

define HELP
$(HELP_DESC)

Default target: all

Available targets

help:		Output this text.
all:		Build all tests.
test:		Run tests.
setup:		Prepares the build directory for go builds.
fmt:		Runs go fmt on all GO source directories.
lint:		Runs various linters on all GO source directories.
vendor:		Runs the vendor tool on each go project which creates/updates
		the vendors directory for the project. Warning: This target may
		change your vendor source. Remember to commit the result.
update_vendor:  Updates vendoring with missing packages and cleanup.
clean:			Remove everything but source (use before commits).

Overridable variables (Set on make command line, e.g.
'make test GO_TEST_FLAGS=-v'):

GO_TEST_FLAGS:		Additional 'go test' flags.
GO_LINT_FLAGS:		Additional go linter flags.
endef


.DEFAULT_GOAL := _default
SHELL := /bin/bash
.PHONY: help setup test fmt vendor update_vendor clean _default

TGTS_DIR               := $(CURDIR)/_tgts

GO_SRC_DIR             := $(CURDIR)
GO_VENDOR_DIR          := $(CURDIR)/vendor
GO_BUILD_DIR           := $(CURDIR)/_gobuild
GO_BUILD_BIN_DIR       := $(GO_BUILD_DIR)/bin
GO_BUILD_CACHE_DIR     := $(GO_BUILD_DIR)/cache
GO_BUILD_ETC_DIR       := $(GO_BUILD_DIR)/etc
GO_BUILD_DIRS          := $(GO_BUILD_DIR) $(GO_BUILD_BIN_DIR) $(GO_BUILD_ETC_DIR)
GO_PATH                := $(GO_BUILD_DIR)
GO_BIN                 := $(GO_BUILD_BIN_DIR)
GO_CACHE               := $(GO_BUILD_CACHE_DIR)

GO_BUILD_FLAGS         :=
GO_INSTALL_FLAGS       :=
GO_TEST_FLAGS          :=

ifeq ($(TRAVIS), true)
TRAVIS_GO_TEST_FLAGS=-v
endif

GO_LINTER_NAME         := golangci-lint
GO_LINTER              := $(shell which $(GO_LINTER_NAME) 2>/dev/null)
GO_LINTER_INCLUDES     :=
GO_LINTER_EXCLUDES     := \
	gochecknoglobals \
	gochecknoinits
GO_LINTER_OPTIONS      := \
			run \
			--enable-all \
			--tests \
			--deadline=600s \
			$(addprefix -E, $(GO_LINTER_INCLUDES)) $(addprefix -D, $(GO_LINTER_EXCLUDES))

define check_module_support
	set -e ; \
	go_version=$$(go version | cut -f3 -d " ") ; \
	go_major_version=$$(go version | sed -re 's!^.*go([0-9]+)\.([0-9]+).*!\1!') ; \
	go_minor_version=$$(go version | sed -re 's!^.*go([0-9]+)\.([0-9]+).*!\2!') ; \
	if [ $${go_version} != "devel" -a $${go_major_version} -lt 2 -a $${go_minor_version} -lt 11 ]; then \
		echo "error: Vendoring is not supported for go versions < 1.11"  >&2 ; \
		exit 1 ; \
	fi
endef

help:
	$(info $(HELP))
	@true

$(TGTS_DIR)/dirs.tgt:
	mkdir -p $(TGTS_DIR)
	touch $@

dirs: $(TGTS_DIR)/dirs.tgt

$(TGTS_DIR)/setup.tgt: $(TGTS_DIR)/dirs.tgt
	mkdir -p $(GO_BUILD_DIRS)
	touch $@

setup: $(TGTS_DIR)/setup.tgt

test: $(TGTS_DIR)/setup.tgt
	cd $(GO_SRC_DIR); GOCACHE=$(GO_CACHE) GOBIN=$(GO_BIN) go test -race $(TRAVIS_GO_TEST_FLAGS) $(GO_TEST_FLAGS) $$(go list ./... | grep -v vendor)

all: test

fmt:
	cd $(CURDIR); GOCACHE=$(GO_CACHE) go fmt $$(go list ./... | grep -v vendor)

lint: $(TGTS_DIR)/setup.tgt
	@if [ -z "$(GO_LINTER)" ]; then \
		echo "error: Go linter '$(GO_LINTER_NAME)' not available." >&2 ; \
		exit 1 ; \
	fi
	cd $(GO_SRC_DIR); $(GO_LINTER_NAME) $(GO_LINTER_OPTIONS) $(GO_LINT_FLAGS)

vendor: setup
	$(check_module_support)
	rm -rf $(GO_BUILD_CACHE_DIR) $(GO_SRC_DIR)/go.mod $(GO_SRC_DIR)/go.sum $(GO_VENDOR_DIR)
	cd $(GO_SRC_DIR); GOCACHE=$(GO_CACHE) go mod init $(GO_REPO_SITE)/$(GO_REPO_SUBDIR)
	cd $(GO_SRC_DIR); GOCACHE=$(GO_CACHE) go mod vendor
	cd $(GO_SRC_DIR); GOCACHE=$(GO_CACHE) go mod tidy

update_vendor: setup
	$(check_module_support)
	cd $(GO_SRC_DIR); GOCACHE=$(GO_CACHE) go mod tidy
	cd $(GO_SRC_DIR); GOCACHE=$(GO_CACHE) go mod vendor

clean:
	rm -rf $(GO_BUILD_DIR)
	rm -rf $(TGTS_DIR)

_default: all
