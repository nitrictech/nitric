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

sec:
	@go run github.com/securego/gosec/v2/cmd/gosec@latest -exclude-dir=tools ./...

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	$(GOLANGCI_LINT) run --fix

lint-contracts:
	@echo "Linting contracts"
	@go run github.com/yoheimuta/protolint/cmd/protolint ../nitric/

lint-contracts-fix:
	@echo "Linting contracts"
	@go run github.com/yoheimuta/protolint/cmd/protolint -fix ../nitric/

lint: lint-contracts
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

PROTO_PREFIX:=github.com/nitrictech/nitric/core

clean-proto:
	@rm -rf ./pkg/proto/

# Generate interfaces
generate-proto: install-tools check-gopath clean-proto
	@echo Generating Proto Sources
	@mkdir -p ./pkg/api/
	@$(PROTOC) --go_out=. --go_opt=module=$(PROTO_PREFIX) --go-grpc_opt=require_unimplemented_servers=false,module=$(PROTO_PREFIX) --go-grpc_out=. ../nitric/proto/*/**/*.proto -I ../

clean-mocks:
	@rm -rf ./mocks/

# generate mock implementations
generate-mocks: clean-mocks
	@echo Generating Mock Clients
	@mkdir -p mocks/sync
	@mkdir -p mocks/gateway
	@mkdir -p mocks/workers/apis
	@mkdir -p mocks/workers/http
	@mkdir -p mocks/workers/schedules
	@mkdir -p mocks/workers/topics
	@mkdir -p mocks/workers/storage
	@mkdir -p mocks/workers/topics
	@mkdir -p mocks/workers/websockets
	@go run github.com/golang/mock/mockgen sync Locker > mocks/sync/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/gateway GatewayService > mocks/gateway/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/workers/apis ApiRequestHandler > mocks/workers/apis/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/workers/http HttpRequestHandler > mocks/workers/http/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/workers/schedules ScheduleRequestHandler > mocks/workers/schedules/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/workers/storage BucketRequestHandler > mocks/workers/storage/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/workers/topics SubscriptionRequestHandler > mocks/workers/topics/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/core/pkg/workers/websockets WebsocketRequestHandler > mocks/workers/websockets/mock.go

generate-sources: generate-proto generate-mocks