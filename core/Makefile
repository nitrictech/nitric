ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

include tools/tools.mk

init: check-gopath go-mod-download install-tools

.PHONY: check fmt lint
check: lint test

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	$(GOLANGCI_LINT) run --fix

lint:
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	$(GOLANGCI_LINT) run

go-mod-download:
	@echo installing go dependencies
	@go mod download

clean: check-gopath
	@rm -rf ./bin/
	@rm -rf ./lib/
	@rm -rf ./interfaces/
	@rm -f ${GOPATH}/bin/protoc-gen-go ${GOPATH}/bin/protoc-gen-go-grpc ${GOPATH}/bin/protoc-gen-validate:

# Run the integration tests
test-integration: install-tools generate-proto
	@echo Running integration tests
	@go run github.com/onsi/ginkgo/ginkgo ./tests/...

# Run the unit tests
test: install-tools generate-mocks generate-proto
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./pkg/...

test-coverage: install-tools generate-proto generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./pkg/...

check-gopath:
ifndef GOPATH
  $(error GOPATH is undefined)
endif

.PHONY: generate generate-proto generate-mocks
generate: generate-proto generate-mocks

# Generate interfaces
generate-proto: install-tools check-gopath
	@echo Generating Proto Sources
	@mkdir -p ./pkg/api/
	@$(PROTOC) --go_out=./pkg/api/ --go_opt=module=github.com/nitrictech/nitric/core/pkg/api --validate_out="lang=go,module=github.com/nitrictech/nitric/core/pkg/api:./pkg/api" --go-grpc_out=module=github.com/nitrictech/nitric/core/pkg/api:./pkg/api ../contracts/proto/*/**/*.proto -I ../contracts

# generate mock implementations
generate-mocks:
	@echo Generating Mock Clients
	@mkdir -p mocks/storage
	@mkdir -p mocks/worker
	@mkdir -p mocks/adapter
	@mkdir -p mocks/pool
	@mkdir -p mocks/nitric
	@mkdir -p mocks/sync
	@mkdir -p mocks/plugins/events
	@mkdir -p mocks/gateway
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/api/nitric/v1 FaasService_TriggerStreamServer > mocks/nitric/mock.go
	@go run github.com/golang/mock/mockgen sync Locker > mocks/sync/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/plugins/document DocumentService > mocks/document/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/plugins/secret SecretService > mocks/secret/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/plugins/storage StorageService > mocks/storage/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/plugins/queue QueueService > mocks/queue/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/plugins/gateway GatewayService > mocks/gateway/mock.go
	@go run github.com/golang/mock/mockgen -package worker github.com/nitrictech/nitric/core/pkg/worker Worker > mocks/worker/mock.go
	@go run github.com/golang/mock/mockgen -package worker github.com/nitrictech/nitric/core/pkg/worker/pool WorkerPool > mocks/pool/mock.go
	@go run github.com/golang/mock/mockgen -package worker github.com/nitrictech/nitric/core/pkg/worker/adapter Adapter > mocks/adapter/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/plugins/events EventService > mocks/plugins/events/mock.go

generate-sources: generate-proto generate-mocks