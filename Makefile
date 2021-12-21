################################## Parameter Definition And Check ##########################################
PKG 						    = "$(PROJECT_NAME)"
PROJECT_NAME 				    = ImprisonSlowSQL
buildDate 						= $(shell TZ=Asia/Shanghai date +%FT%T%z)

override GIT_COMMIT     		= $(shell git rev-parse --short HEAD || echo "GitNotFound")
override VERSION_DIR 			= ImprisonSlowSQL/pkg/version
override GOOS           		= linux
override OS_VERSION 			= el7
override GOARCH         		= amd64
override RPMBUILD_TARGET		= x86_64
override RELEASE 				= qa
override GO_BUILD_FLAGS 		= -mod=vendor
override GO_BUILD_TAGS			= dummyhead
override PROJECT_VERSION		= $(shell cat VERSION | grep 'version' | awk -F ' ' '{print $$2}' | awk '{gsub(/ /,"")}1')
override LDFLAGS 				= "-X \"${VERSION_DIR}.version=${PROJECT_VERSION}\" -X \"${VERSION_DIR}.gitTag=${gitTag}\" -X \"${VERSION_DIR}.buildDate=${buildDate}\" -X \"${VERSION_DIR}.gitCommit=${gitCommit}\""

## Dynamic Parameter
GO_COMPILER_IMAGE ?= golang:1.16
RPM_BUILD_IMAGE ?= rpmbuild/centos7

## Static Parameter, should not be overwrite
GOBIN = ${shell pwd}/bin

# init environment variables
export PATH        := $(shell go env GOPATH)/bin:$(PATH)
export GOPATH      := $(shell go env GOPATH)
export GO111MODULE := on


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
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -ldflags ${LDFLAGS} -tags $(GO_BUILD_TAGS) -o $(GOBIN)/$(PROJECT_NAME) ./cmd/main.go

build: clean
	go mod tidy && go mod vendor

debug: install
	./cmd/$(PROJECT_NAME) -h 192.168.0.1 -p 12345

.PHONY: dlv-build
dlv-build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -ldflags ${LDFLAGS} -tags $(GO_BUILD_TAGS) -gcflags "all=-N -l" -o $(GOBIN)/$(PROJECT_NAME) ./cmd/main.go

.PHONY: dlv
# make dlv 远程调试
dlv: dlv-build
	cd ./cmd && dlv debug --headless --listen=:2345 --api-version=2 -- -i 127.0.0.1 -p 58_v29rC -d shop-fstv


.PHONY: help
help:
	$(warning ---------------------------------------------------------------------------------)
	$(warning Supported Variables And Values:)
	$(warning ---------------------------------------------------------------------------------)
	$(foreach v, $(.VARIABLES), $(if $(filter file,$(origin $(v))), $(info $(v)=$($(v)))))
