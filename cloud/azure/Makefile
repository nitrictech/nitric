ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

binaries: deploybin

sec:
	@touch deploy/runtime/runtime-azure
	@go run github.com/securego/gosec/v2/cmd/gosec@latest -exclude-dir=tools ./...
	@rm deploy/runtime/runtime-azure

runtimebin:
	@echo Building Azure Runtime Server
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/runtime-azure -ldflags="-s -w -extldflags=-static" ./cmd/runtime

predeploybin: runtimebin
	@cp bin/runtime-azure common/runtime/runtime-azure

deploybin: predeploybin
	@echo Building Azure Deployment Server
	@CGO_ENABLED=0 go build -o bin/deploy-azure -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploy
	@rm deploy/runtime/runtime-azure

deploybintf: generate-terraform predeploybin
	@echo Building AWS Terraform Deployment Server
	@CGO_ENABLED=0 go build -o bin/deploy-azuretf -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploytf

install: deploybin
	@echo installing azure deployment server to ${HOME}/.nitric/providers/nitric/azure-0.0.1
	@mkdir -p ${HOME}/.nitric/providers/nitric/
	@if [ "$(OS)" == "Windows_NT" ]; then \
		rm -f "${HOME}/.nitric/providers/nitric/azure-0.0.1.exe"; \
		cp bin/deploy-azure "${HOME}/.nitric/providers/nitric/azure-0.0.1.exe"; \
	else \
		rm -f "${HOME}/.nitric/providers/nitric/azure-0.0.1"; \
		cp bin/deploy-azure "${HOME}/.nitric/providers/nitric/azure-0.0.1"; \
	fi

license-check: runtimebin
	@echo Checking Azure Runtime OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/runtime-azure

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch deploy/runtime/runtime-azure
	$(GOLANGCI_LINT) run --fix
	@rm deploy/runtime/runtime-azure

lint:
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch deploy/runtime/runtime-azure
	$(GOLANGCI_LINT) run
	@rm deploy/runtime/runtime-azure

test: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./runtime/...

test-coverage: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./runtime/...

clean-mocks:
	@rm -rf mocks

generate-mocks: clean-mocks
	@echo Generating Mock Clients
	@mkdir -p mocks/key_vault
	@mkdir -p mocks/azblob
	@mkdir -p mocks/azqueue
	@mkdir -p mocks/mock_event_grid
	@mkdir -p mocks/provider
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/azure/runtime/resource AzResourceResolver > mocks/provider/azure.go
	@go run github.com/golang/mock/mockgen -package mock_azblob github.com/Azure/azure-storage-blob-go/azblob StorageError > mocks/azblob/error.go
	@go run github.com/golang/mock/mockgen -package mock_azqueue github.com/nitrictech/nitric/cloud/azure/runtime/queue/iface AzqueueServiceUrlIface,AzqueueQueueUrlIface,AzqueueMessageUrlIface,DequeueMessagesResponseIface,AzqueueMessageIdUrlIface > mocks/azqueue/mock.go
	@go run github.com/golang/mock/mockgen -package mock_azblob github.com/nitrictech/nitric/cloud/azure/runtime/storage/iface AzblobServiceUrlIface,AzblobContainerUrlIface,AzblobBlockBlobUrlIface,AzblobDownloadResponse > mocks/azblob/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/azure/runtime/secret KeyVaultClient > mocks/key_vault/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid/eventgridapi BaseClientAPI > mocks/mock_event_grid/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid/eventgridapi TopicsClientAPI > mocks/mock_event_grid/topic.go

generate-sources: generate-mocks

tidy:
	@go mod tidy

.PHONY: install license-check fmt lint test test-coverage clean-mocks generate-mocks generate-sources tidy