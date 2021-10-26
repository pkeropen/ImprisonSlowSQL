################################## Parameter Definition And Check ##########################################
#override GIT_VERSION    		= $(shell git rev-parse --abbrev-ref HEAD)${CUSTOM} $(shell git rev-parse HEAD)
override GIT_VERSION    		= $(shell git symbolic-ref HEAD | cut -b 12-`-`git rev-parse HEAD)
override GIT_COMMIT     		= $(shell git rev-parse HEAD)
override PROJECT_NAME 			= ImprisonSlowSQL
override LDFLAGS 				= -ldflags "-X 'main.version=\"${GIT_VERSION}\"'"
override GOOS           		= linux
override OS_VERSION 			= el7
override GOARCH         		= amd64
override RPMBUILD_TARGET		= x86_64
override RELEASE 				= qa
override GO_BUILD_FLAGS 		= -mod=vendor
override GO_BUILD_TAGS			= dummyhead

# Two cases:
# 1. if there is tag on current commit, means that 
# 	 we release new version on current branch just now. 
#    Set rpm name with tag name(v1.2109.0 -> 1.2109.0).
#
# 2. if there is no tag on current commit, means that
#    current branch is on process.
#    Set rpm name with current branch name(release-1.2109.x-ee or release-1.2109.x -> 1.2109.x).
PROJECT_VERSION = $(shell if [ "$$(git tag --points-at HEAD | tail -n1)" ]; then git tag --points-at HEAD | tail -n1 | sed 's/v\(.*\)/\1/'; else git rev-parse --abbrev-ref HEAD | sed 's/release-\(.*\)/\1/' | tr '-' '\n' | head -n1; fi)

## Dynamic Parameter
GO_COMPILER_IMAGE ?= golang:1.16
RPM_BUILD_IMAGE ?= rpmbuild/centos7

## Static Parameter, should not be overwrite
GOBIN = ${shell pwd}/bin

default: install
######################################## Code Check ####################################################
## Static Code Analysis
vet:
	GOOS=$(GOOS) GOARCH=amd64 go vet $$(GOOS=${GOOS} GOARCH=${GOARCH} go list ./...)

## Unit Test
test:
	GOOS=$(GOOS) GOARCH=amd64 go test -v ./...

clean:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go clean	

install:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) ${LDFLAGS} -tags $(GO_BUILD_TAGS) -o $(GOBIN)/main ./cmd/main.go

build: clean
	go mod tidy && go mod vendor

.PHONY: help
help:
	$(warning ---------------------------------------------------------------------------------)
	$(warning Supported Variables And Values:)
	$(warning ---------------------------------------------------------------------------------)
	$(foreach v, $(.VARIABLES), $(if $(filter file,$(origin $(v))), $(info $(v)=$($(v)))))
