BASE_PATH := "$(shell pwd)"

DDCUTIL_BASE_PATH := $(BASE_PATH)/ddcutil
DDCUTIL_BIN_PATH := $(DDCUTIL_BASE_PATH)/bin/ddcutil

BIN_BASE_PATH := $(BASE_PATH)
BIN_NAME := LinuxMonitorControl
BIN_PATH := $(BIN_BASE_PATH)/$(BIN_NAME)

# Strips the binary out of the symbol table and debug information
# https://lukeeckley.com/post/useful-go-build-flags/
LDFLAGS_RELEASE ?= -s -w
export LDFLAGS ?=

# Build time in UTC time
BUILD_TIME ?= $(shell date --utc +"%Y.%m.%d.%H%M")

GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

# The root package path
IMPORT_PATH_BASE ?= github.com/vampy/LinuxMonitorControl

# Call make with: WITH_RACE=1
GO_RACE :=
ifdef WITH_RACE
	GO_RACE := -race
endif

GO_BUILD_COMMAND := go build -v -ldflags "$(LDFLAGS)\
	-X '$(IMPORT_PATH_BASE)/pkg/build.GitBranch=$(GIT_BRANCH)' \
	-X '$(IMPORT_PATH_BASE)/pkg/build.GitCommit=$(GIT_COMMIT)' \
	-X '$(IMPORT_PATH_BASE)/pkg/build.BuildTime=$(BUILD_TIME)'"\
	$(GO_RACE) \
	-o

init:
## install: Run go modules download
	@echo "Downloading dependencies:"
	go mod download
	@echo ""

dependencies:
## ddcutil: Download & ddcutil ddcutil
	@echo "Building and downloading ddcutil:"
	$(DDCUTIL_BASE_PATH)/build.sh
	@echo ""

tidy:
## tidy: Run go modules tidy
	go mod tidy -v

build-packages:
## build-packages: Build all the packages (NO BINARIES)
	go build -v ./...

build: init
## build: Builds the binary
	@echo "Building:"
	$(GO_BUILD_COMMAND) $(BIN_PATH) main.go || exit 1
	@echo ""

build-release:
	# $(eval LDFLAGS += $(LDFLAGS_RELEASE))
	$(MAKE) build

update-packages: tidy
## update-packages: Update the go packages to their latest version
	go get -v -u all

test:
## test: Run all tests
	go test -race -cover ./...


release: DST_NAME = $(BIN_NAME)-$(GIT_COMMIT)
release: DST_BASE_PATH = $(BASE_PATH)/releases/
release: DST_DIR_PATH = $(DST_BASE_PATH)/$(DST_NAME)
release: dependencies build-release
	@echo ""
	@echo "Packaging for release:"
	mkdir -p $(DST_DIR_PATH)
	cp --force $(BIN_PATH) $(DST_DIR_PATH)
	cp --force $(DDCUTIL_BIN_PATH) $(DST_DIR_PATH)
	rm -f "$(DST_BASE_PATH)/$(DST_NAME).zip"
	zip --junk-paths -r $(DST_BASE_PATH)/$(DST_NAME).zip $(DST_DIR_PATH)

help:
## help: This helpful list of commands
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/-/'
