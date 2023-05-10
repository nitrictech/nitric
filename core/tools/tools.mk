ifeq ($(OS),Windows_NT)
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
	PROTOC_PLATFORM := win64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
	    PROTOC_PLATFORM := win64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
	    PROTOC_PLATFORM := win32
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        UNAME_P := $(shell uname -m)
	ifeq ($(UNAME_P),x86_64)
	    PROTOC_PLATFORM := linux-x86_64
	endif
        ifneq ($(filter %86,$(UNAME_P)),)
	    PROTOC_PLATFORM := linux-x86_32
        endif
    endif
    ifeq ($(UNAME_S),Darwin)
        UNAME_P := $(shell uname -m)
        ifeq ($(UNAME_P),X86_64)
	    PROTOC_PLATFORM := osx-x86_64
        endif
        ifeq ($(UNAME_P),arm64)
        PROTOC_PLATFORM := osx-aarch_64
        endif
    endif
endif

ifndef PROTOC_PLATFORM
    $(error unsupported platform $(UNAME_S):$(UNAME_P))
endif

TOOLS_DIR := ./tools
TOOLS_BIN := $(TOOLS_DIR)/bin

PROTOC_VERSION := 21.4
PROTOC_RELEASES_PATH := https://github.com/protocolbuffers/protobuf/releases/download
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-$(PROTOC_PLATFORM).zip
PROTOC_DOWNLOAD := $(PROTOC_RELEASES_PATH)/v$(PROTOC_VERSION)/$(PROTOC_ZIP)
PROTOC := $(TOOLS_BIN)/protoc

install-tools: $(PROTOC) ${GOPATH}/bin/protoc-gen-go ${GOPATH}/bin/protoc-gen-go-grpc ${GOPATH}/bin/protoc-gen-validate ./contracts/validate/validate.proto

./contracts/validate/validate.proto:
	@echo fetching envoyproxy validate contract
	@mkdir -p ./contracts/validate
	@curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.6.1/validate/validate.proto --output ./contracts/validate/validate.proto

$(PROTOC): $(TOOLS_DIR)/$(PROTOC_ZIP)
	unzip -o -d "$(TOOLS_DIR)" $< && touch $@  # avoid Prerequisite is newer than target `tools/bin/protoc'.
	rm $(TOOLS_DIR)/readme.txt

$(TOOLS_DIR)/$(PROTOC_ZIP):
	curl --location $(PROTOC_DOWNLOAD) --output $@

${GOPATH}/bin/protoc-gen-go: go.sum
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

${GOPATH}/bin/protoc-gen-go-grpc: go.sum
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

${GOPATH}/bin/protoc-gen-validate: go.sum
	go install github.com/envoyproxy/protoc-gen-validate

