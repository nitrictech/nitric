binaries: runtimebin

runtimebin:
	@echo Building Azure Runtime Server
	@CGO_ENABLED=0 go build -o bin/runtime-azure -ldflags="-extldflags=-static" ./runtime/cmd

license-check: runtimebin
	@echo Checking Azure Runtime OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/runtime-azure

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
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/azure/runtime/core AzProvider > mocks/provider/azure.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/azure/runtime/storage/iface AzblobServiceUrlIface,AzblobContainerUrlIface,AzblobBlockBlobUrlIface,AzblobDownloadResponse > mocks/azblob/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/azure/runtime/secret KeyVaultClient > mocks/key_vault/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid/eventgridapi BaseClientAPI > mocks/mock_event_grid/mock.go
	@go run github.com/golang/mock/mockgen github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid/eventgridapi TopicsClientAPI > mocks/mock_event_grid/topic.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/azure/runtime/queue/iface AzqueueServiceUrlIface,AzqueueQueueUrlIface,AzqueueMessageUrlIface,AzqueueMessageIdUrlIface,DequeueMessagesResponseIface > mocks/azqueue/mock.go

generate-sources: generate-mocks