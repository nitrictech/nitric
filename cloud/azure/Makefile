ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

binaries: runtimebin

runtimebin:
	@echo Building Azure Runtime Server
	@CGO_ENABLED=0 go build -o bin/runtime-azure -ldflags="-extldflags=-static" ./runtime/cmd

# FIXME: proto server registration error
# There appears to be an old namespace conflict with the protobuf definitions
deploybin: runtimebin
	@echo Building Azure Deployment Server
	cp bin/runtime-azure deploy/runtime-azure
	@CGO_ENABLED=0 go build -o bin/deploy-azure -ldflags="-extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./deploy/cmd
	@rm deploy/runtime-azure

dockerimage: deploybin
	@echo Building Azure Deployment Server Docker Image
	@docker build -t nitrictech/deploy-azure:latest -f deploy/Dockerfile bin/

license-check: runtimebin
	@echo Checking Azure Runtime OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/runtime-azure

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch deploy/runtime-azure
	$(GOLANGCI_LINT) run --fix
	@rm deploy/runtime-azure

lint:
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch deploy/runtime-azure
	$(GOLANGCI_LINT) run
	@rm deploy/runtime-azure

test: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./runtime/...

test-coverage: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./runtime/...

generate-mocks:
	@echo Generating Mock Clients
	@mkdir -p mocks/key_vault
	@mkdir -p mocks/azblob
	@mkdir -p mocks/mock_event_grid
	@mkdir -p mocks/azqueue
	@mkdir -p mocks/provider
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/azure/runtime/core AzProvider > mocks/provider/azure.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/azure/runtime/storage/iface AzblobServiceUrlIface,AzblobContainerUrlIface,AzblobBlockBlobUrlIface,AzblobDownloadResponse > mocks/azblob/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/azure/runtime/secret KeyVaultClient > mocks/key_vault/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid/eventgridapi BaseClientAPI > mocks/mock_event_grid/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid/eventgridapi TopicsClientAPI > mocks/mock_event_grid/topic.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/azure/runtime/queue/iface AzqueueServiceUrlIface,AzqueueQueueUrlIface,AzqueueMessageUrlIface,AzqueueMessageIdUrlIface,DequeueMessagesResponseIface > mocks/azqueue/mock.go

generate-sources: generate-mocks
