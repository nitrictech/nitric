ifeq ($(OS), Windows_NT)
	EXECUTABLE_EXT = .exe
else
	EXECUTABLE_EXT =
endif

.PHONY: build run

# See pkg/cmd/version.go for details
SOURCE_GIT_COMMIT ?= $(shell git rev-parse --short --dirty HEAD)
BUILD_VERSION ?= $(shell git describe --tags --abbrev=7 --dirty)
VERSION_URI = "github.com/nitrictech/nitric/cli/internal/version"
export LDFLAGS="-X $(VERSION_URI).Version=${BUILD_VERSION} \
                -X $(VERSION_URI).Commit=${SOURCE_GIT_COMMIT} \
                -X $(VERSION_URI).BuildTime=$(shell date +%Y-%m-%dT%H:%M:%S%z)"

build:
	go build -ldflags $(LDFLAGS) -o bin/suga$(EXECUTABLE_EXT) main.go

install:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@go mod verify

test:
	@echo "Running tests..."
	@go test ./pkg/...

test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./pkg/...

test-short:
	@echo "Running tests with short output..."
	@go test -short ./pkg/...

test-pretty:
	@go tool gotest.tools/gotestsum --format=pkgname ./pkg/...
