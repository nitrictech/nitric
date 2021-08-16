init: install-tools
	@echo Installing git hooks
	@find .git/hooks -type l -exec rm {} \; && find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

fmt:
	@echo Formatting Code
	@gofmt -s -w ./**/*.go

lint:
	@echo Formatting Code
	@golint ./...

install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get %

# Install integration testing tools
install-test-tools:
	@wget https://s3.us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.tar.gz
	@sudo mkdir -p /usr/local/dynamodb
	@sudo tar -xf dynamodb_local_latest.tar.gz -C /usr/local/dynamodb
	@rm dynamodb_local_latest.tar.gz

clean:
	@rm -rf ./bin/
	@rm -rf ./lib/
	@rm -rf ./interfaces/

# Run the integration tests
test-integration: install-tools generate-proto
	@echo Running integration tests
	@go run github.com/onsi/ginkgo/ginkgo ./tests/...

# Run the unit tests
test: install-tools generate-mocks generate-proto
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./pkg/...

test-coverage: install-tools generate-mocks generate-proto
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./pkg/...

license-check-dev: dev-static
	@echo Checking Dev Membrane OSS Licenses
	@lichen --config=./lichen.yaml ./bin/membrane

license-check-aws: aws-static
	@echo Checking AWS Membrane OSS Licenses
	@lichen --config=./lichen.yaml ./bin/membrane

license-check-gcp: gcp-static
	@echo Checking GCP Membrane OSS Licenses
	@lichen --config=./lichen.yaml ./bin/membrane

license-check-azure: azure-static
	@echo Checking Azure Membrane OSS Licenses
	@lichen --config=./lichen.yaml ./bin/membrane

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

license-header-add:
	@echo Add License Headers to Source Files
	@addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)

license-header-check:
	@echo Checking License Headers to Source Files
	@addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)

license-check: install-tools license-check-dev license-check-aws license-check-gcp license-check-azure
	@echo Checking OSS Licenses

check-gopath:
ifndef GOPATH
  $(error GOPATH is undefined)
endif

# Generate interfaces
generate-proto: install-tools check-gopath
	@echo Generating Proto Sources
	@GO111MODULE=off go get github.com/envoyproxy/protoc-gen-validate
	@mkdir -p ./interfaces/
	@protoc --go_out=./interfaces/ --validate_out="lang=go:./interfaces/" --go-grpc_out=./interfaces/ -I ./contracts/proto ./contracts/proto/*/**/*.proto -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate

# BEGIN AWS Plugins
aws-static: generate-proto
	@echo Building static AWS membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/aws/membrane.go

# Cross-platform Build
aws-static-xp: generate-proto
	@echo Building static AWS membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/aws/membrane.go

# # Service Factory Plugin for Pluggable Membrane
# aws-plugin:
# 	@echo Building AWS Service Factory Plugin
# 	@go build -buildmode=plugin -o lib/plugins/aws.so ./providers/aws/plugin.go

aws-docker-static:
	@docker build . -f ./pkg/providers/aws/aws.dockerfile -t nitricimages/membrane-aws

aws-docker: aws-docker-static
	@echo Built AWS Docker Images
# END AWS Plugins

# BEGIN Azure Plugins
azure-static: generate-proto
	@echo Building static Azure membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/azure/membrane.go

# Cross-platform Build
azure-static-xp: generate-proto
	@echo Building static Azure membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/azure/membrane.go

# # Service Factory Plugin for Pluggable Membrane
# azure-plugin:
# 	@echo Building Azure Service Factory Plugin
# 	@go build -buildmode=plugin -o lib/plugins/azure.so ./pkg/providers/azure/plugin.go

azure-docker-static:
	@docker build . -f ./pkg/providers/azure/azure.dockerfile -t nitricimages/membrane-azure

azure-docker: azure-docker-static # azure-docker-alpine azure-docker-debian
	@echo Built Azure Docker Images
# END Azure Plugins

gcp-static: generate-proto
	@echo Building static GCP membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/gcp/membrane.go

# Cross-platform Build
gcp-static-xp: generate-proto
	@echo Building static GCP membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/gcp/membrane.go

# # Service Factory Plugin for Pluggable Membrane
# gcp-plugin:
# 	@echo Building GCP Service Factory Plugin
# 	@go build -buildmode=plugin -o lib/plugins/gcp.so ./providers/gcp/plugin.go

gcp-docker-static:
	@docker build . -f ./pkg/providers/gcp/gcp.dockerfile -t nitricimages/membrane-gcp

gcp-docker: gcp-docker-static # gcp-docker-alpine gcp-docker-debian
	@echo Built GCP Docker Images
# END GCP Plugins

# BEGIN Local Plugins
# Cross-platform build only, this membrane is not for production use.
dev-static: generate-proto
	@echo Building static Local membrane
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/dev/membrane.go

# # Service Factory Plugin for Pluggable Membrane
# dev-plugin:
# 	@echo Building Dev Service Factory Plugin
# 	@go build -buildmode=plugin -o lib/plugins/dev.so ./pkg/providers/dev/plugin.go

dev-docker-static:
	@docker build . -f ./pkg/providers/dev/dev.dockerfile -t nitricimages/membrane-local

dev-docker: dev-docker-static
	@echo Built Local Docker Images
# END Local Plugins

# BEGIN DigitalOcean Plugins
do-static: generate-proto
	@CGO_ENABLED=0 go build -o bin/membrane -ldflags="-extldflags=-static" ./pkg/providers/do/membrane.go

do-docker-static:
	@docker build . -f ./pkg/providers/do/do.dockerfile -t nitricimages/membrane-do

do-docker: do-docker-static
	@echo Built Digital Ocean Docker Images
# END DigitalOcean Plugins

# membrane-docker-alpine: generate-proto
# 	@docker build . -f alpine.dockerfile -t nitric:membrane-alpine
# membrane-docker-debian: generate-proto
# 	@docker build . -f debian.dockerfile -t nitric:membrane-debian

# # Generate proto files locally before building docker images
# # TODO: Get alpine image generating its own sources
# membrane-docker: generate-proto membrane-docker-alpine membrane-docker-debian
# 	@echo Built Docker Images

# generate mock implementations
generate-mocks:
	@echo Generating Mock Clients
	@mkdir -p mocks/mock_secret_manager
	@mkdir -p mocks/secrets_manager
	@go run github.com/golang/mock/mockgen github.com/nitric-dev/membrane/pkg/plugins/secret/secret_manager SecretManagerClient > mocks/mock_secret_manager/mock.go
	@go run github.com/golang/mock/mockgen github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface SecretsManagerAPI > mocks/secrets_manager/mock.go