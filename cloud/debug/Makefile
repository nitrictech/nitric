ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: build
build:
	@echo building debug deployment server
	@CGO_ENABLED=0 go build -o bin/deploy-debug -ldflags="-extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./deploy/cmd

.PHONY: install
install: build
	@echo installing debug deployment server to ${HOME}/.nitric/providers/debug/spec-0.0.1
	@mkdir -p ${HOME}/.nitric/providers/debug/
	@if [ "$(OS)" == "Windows_NT" ]; then \
		rm -f "${HOME}/.nitric/providers/debug/spec-0.0.1.exe"; \
		cp bin/deploy-debug "${HOME}/.nitric/providers/debug/spec-0.0.1.exe"; \
	else \
		rm -f "${HOME}/.nitric/providers/debug/spec-0.0.1"; \
		cp bin/deploy-debug "${HOME}/.nitric/providers/debug/spec-0.0.1"; \
	fi

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	$(GOLANGCI_LINT) run --fix

lint:
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	$(GOLANGCI_LINT) run