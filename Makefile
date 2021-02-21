PROJECTNAME := "$(shell basename "$(PWD)")"

GOBASE := "$(shell pwd)"
GOBIN := $(GOBASE)

# Strips the binary out of the symbol table and debug information
# https://lukeeckley.com/post/useful-go-build-flags/
LDFLAGS_RELEASE ?= -s -w

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

GO_BUILD_COMMAND := go build -v -ldflags "\
	-X '$(IMPORT_PATH_BASE)/pkg/build.GitBranch=$(GIT_BRANCH)' \
	-X '$(IMPORT_PATH_BASE)/pkg/build.GitCommit=$(GIT_COMMIT)' \
	-X '$(IMPORT_PATH_BASE)/pkg/build.BuildTime=$(BUILD_TIME)'"\
	$(GO_RACE) \
	-o

# Stop annying warnings https://github.com/mattn/go-sqlite3/issues/803
export CGO_CFLAGS ?= -Wno-return-local-addr

init:
## install: Run go modules download
	@echo "Downloading dependencies:"
	go mod download
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
	$(GO_BUILD_COMMAND) "$(GOBIN)/LinuxMonitorControl" main.go || exit 1
	@echo ""

update-packages: tidy
## update-packages: Update the go packages to their latest version
	go get -v -u all

test:
## test: Run all tests
	go test -race -cover ./...


help:
## help: This helpful list of commands
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/-/'
