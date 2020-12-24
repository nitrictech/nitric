membrane: install
	@echo Building Go Project...
	@CGO_ENABLED=1 GOOS=linux go build -o bin/membrane pluggable_membrane.go

install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	@rm -rf ./bin/
	@rm -rf ./lib/

# Run all tests
test: test-membrane test-aws-plugins test-gcp-plugins test-local-plugins
	@echo Done.

# Generate interfaces
generate-proto:
	@echo Generating Proto Sources
	@mkdir -p ./interfaces/
	@protoc --go_out=./interfaces/ --go-grpc_out=./interfaces/ -I ./contracts/proto/ ./contracts/proto/**/*.proto

# Test the membrane
test-membrane: install-tools generate-proto 
	@echo Testing Membrane
	@go run github.com/onsi/ginkgo/ginkgo -cover ./membrane/...

# BEGIN AWS Plugins
test-aws-plugins:
	@echo Testing AWS Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/aws/...

aws-static:
	@echo Building static AWS membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/aws/static_membrane.go

aws-plugins:
	@echo Building AWS plugins
	@go build -buildmode=plugin -o lib/documents/dynamodb.so ./plugins/aws/plugins/dynamodb.go
	@go build -buildmode=plugin -o lib/eventing/sns.so ./plugins/aws/plugins/sns.go
	@go build -buildmode=plugin -o lib/gateway/lambda.so ./plugins/aws/plugins/lambda.go
	@go build -buildmode=plugin -o lib/storage/s3.so ./plugins/aws/plugins/s3.go

aws-docker-alpine:
	@docker build . -f ./plugins/aws/alpine.dockerfile -t nitric:membrane-alpine-aws
aws-docker-debian:
	@docker build . -f ./plugins/aws/debian.dockerfile -t nitric:membrane-debian-aws
aws-docker-static: generate-proto
	@docker build . -f ./plugins/aws/aws.dockerfile -t nitric:membrane-aws

aws-docker: aws-docker-static aws-docker-alpine aws-docker-debian 
	@echo Built AWS Docker Images
# END AWS Plugins

# BEGIN GCP Plugins
test-gcp-plugins:
	@echo Testing GCP Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/gcp/...

gcp-static:
	@echo Building static GCP membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/gcp/static_membrane.go

gcp-plugins:
	@echo Building GCP plugins
	@go build -buildmode=plugin -o lib/documents/firestore.so ./plugins/gcp/plugins/firestore.go
	@go build -buildmode=plugin -o lib/eventing/pubsub.so ./plugins/gcp/plugins/pubsub.go
	@go build -buildmode=plugin -o lib/gateway/http.so ./plugins/gcp/plugins/http.go
	@go build -buildmode=plugin -o lib/storage/storage.so ./plugins/gcp/plugins/storage.go

gcp-docker-alpine:
	@docker build . -f ./plugins/gcp/alpine.dockerfile -t nitric:membrane-alpine-gcp
gcp-docker-debian:
	@docker build . -f ./plugins/gcp/debian.dockerfile -t nitric:membrane-debian-gcp
gcp-docker-static: generate-proto
	@docker build . -f ./plugins/gcp/gcp.dockerfile -t nitric:membrane-gcp

gcp-docker: gcp-docker-static gcp-docker-alpine gcp-docker-debian
	@echo Built GCP Docker Images
# END GCP Plugins

# BEGIN Local Plugins
local-static:
	@echo Building static Local membrane
	@CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" ./plugins/dev/static_membrane.go

local-plugins:
	@echo Building Local plugins
	@go build -buildmode=plugin -o lib/documents.so ./plugins/dev/plugins/documents.go
	@go build -buildmode=plugin -o lib/eventing.so ./plugins/dev/plugins/eventing.go
	@go build -buildmode=plugin -o lib/gateway.so ./plugins/dev/plugins/gateway.go
	@go build -buildmode=plugin -o lib/storage.so ./plugins/dev/plugins/storage.go

local-docker-alpine:
	@docker build . -f ./plugins/dev/alpine.dockerfile -t nitric:membrane-alpine-local
local-docker-debian:
	@docker build . -f ./plugins/dev/debian.dockerfile -t nitric:membrane-debian-local
local-docker-static: generate-proto
	@docker build . -f ./plugins/dev/dev.dockerfile -t nitric:membrane-local

local-docker: local-docker-static local-docker-alpine local-docker-debian
	@echo Built Local Docker Images

test-local-plugins:
	@echo Testing Local Plugins
	@go run github.com/onsi/ginkgo/ginkgo -cover ./plugins/dev/...
# END Local Plugins

membrane-docker-alpine: generate-proto
	@docker build . -f alpine.dockerfile -t nitric:membrane-alpine
membrane-docker-debian: generate-proto
	@docker build . -f debian.dockerfile -t nitric:membrane-debian

# Generate proto files locally before building docker images
# TODO: Get alpine image generating its own sources
membrane-docker: generate-proto membrane-docker-alpine membrane-docker-debian
	@echo Built Docker Images