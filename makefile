ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

include tools/tools.mk

init: check-gopath go-mod-download install-tools
	@echo Installing git hooks
	@find .git/hooks -type l -exec rm {} \; && find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

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

license-check-dev: dev-static
	@echo Checking Dev Membrane OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/membrane

license-check-aws: aws-static
	@echo Checking AWS Membrane OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/membrane

license-check-gcp: gcp-static
	@echo Checking GCP Membrane OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/membrane

license-check-azure: azure-static
	@echo Checking Azure Membrane OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/membrane

license-check: install-tools license-check-dev license-check-aws license-check-gcp license-check-azure
	@echo Checking OSS Licenses

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
	@$(PROTOC) --go_out=./pkg/api/ --validate_out="lang=go:./pkg/api" --go-grpc_out=./pkg/api -I ./contracts/proto ./contracts/proto/*/**/*.proto -I ./contracts

# BEGIN AWS Plugins
aws-static: generate-proto
	@echo Building static AWS membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/aws/membrane.go

# Cross-platform Build
aws-static-xp: generate-proto
	@echo Building static AWS membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/aws/membrane.go
# END AWS Plugins

# BEGIN Azure Plugins
azure-static: generate-proto
	@echo Building static Azure membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/azure/membrane.go

# Cross-platform Build
azure-static-xp: generate-proto
	@echo Building static Azure membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/azure/membrane.go
# END Azure Plugins

gcp-static: generate-proto
	@echo Building static GCP membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/gcp/membrane.go

# Cross-platform Build
gcp-static-xp: generate-proto
	@echo Building static GCP membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/gcp/membrane.go
# END GCP Plugins

# BEGIN Local Plugins
# Cross-platform build only, this membrane is not for production use.
dev-static: generate-proto
	@echo Building static Local membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/dev/membrane.go
# END Local Plugins

# BEGIN DigitalOcean Plugins
do-static: generate-proto
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/do/membrane.go
# END DigitalOcean Plugins

build-all-binaries: clean generate-proto
	@echo Building all provider membranes
	@CGO_ENABLED=0 go build -o bin/membrane-gcp -ldflags="-extldflags=-static" ./pkg/providers/gcp/membrane.go
	@CGO_ENABLED=0 go build -o bin/membrane-aws -ldflags="-extldflags=-static" ./pkg/providers/aws/membrane.go
	@CGO_ENABLED=0 go build -o bin/membrane-azure -ldflags="-extldflags=-static" ./pkg/providers/azure/membrane.go
	@CGO_ENABLED=0 go build -o bin/membrane-do -ldflags="-extldflags=-static" ./pkg/providers/do/membrane.go
	@CGO_ENABLED=0 go build -o bin/membrane-dev -ldflags="-extldflags=-static" ./pkg/providers/dev/membrane.go

# generate mock implementations
generate-mocks:
	@echo Generating Mock Clients
	@mkdir -p mocks/secret_manager
	@mkdir -p mocks/secrets_manager
	@mkdir -p mocks/key_vault
	@mkdir -p mocks/s3
	@mkdir -p mocks/sns
	@mkdir -p mocks/sqs
	@mkdir -p mocks/azblob
	@mkdir -p mocks/mock_event_grid
	@mkdir -p mocks/azqueue
	@mkdir -p mocks/worker
	@mkdir -p mocks/nitric
	@mkdir -p mocks/sync
	@mkdir -p mocks/provider
	@mkdir -p mocks/resourcetaggingapi
	@go run github.com/golang/mock/mockgen github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface ResourceGroupsTaggingAPIAPI > mocks/resourcetaggingapi/mock.go
	@go run github.com/golang/mock/mockgen github.com/aws/aws-sdk-go/service/sns/snsiface SNSAPI > mocks/sns/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/providers/aws/core AwsProvider > mocks/provider/aws.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/providers/azure/core AzProvider > mocks/provider/azure.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/api/nitric/v1 FaasService_TriggerStreamServer > mocks/nitric/mock.go
	@go run github.com/golang/mock/mockgen sync Locker > mocks/sync/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/secret/secret_manager SecretManagerClient > mocks/secret_manager/mock.go
	@go run github.com/golang/mock/mockgen github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface SecretsManagerAPI > mocks/secrets_manager/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/storage/azblob/iface AzblobServiceUrlIface,AzblobContainerUrlIface,AzblobBlockBlobUrlIface,AzblobDownloadResponse > mocks/azblob/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/secret/key_vault KeyVaultClient > mocks/key_vault/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/document DocumentService > mocks/document/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/secret SecretService > mocks/secret/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/storage StorageService > mocks/storage/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/queue QueueService > mocks/queue/mock.go
	@go run github.com/golang/mock/mockgen -package worker github.com/nitrictech/nitric/pkg/worker GrpcWorker > mocks/worker/mock.go
	@go run github.com/golang/mock/mockgen github.com/aws/aws-sdk-go/service/s3/s3iface S3API > mocks/s3/mock.go
	@go run github.com/golang/mock/mockgen github.com/aws/aws-sdk-go/service/sqs/sqsiface SQSAPI > mocks/sqs/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid/eventgridapi BaseClientAPI > mocks/mock_event_grid/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid/eventgridapi TopicsClientAPI > mocks/mock_event_grid/topic.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/pkg/plugins/queue/azqueue/iface AzqueueServiceUrlIface,AzqueueQueueUrlIface,AzqueueMessageUrlIface,AzqueueMessageIdUrlIface,DequeueMessagesResponseIface > mocks/azqueue/mock.go

generate-sources: generate-proto generate-mocks