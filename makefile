membrane: install
	@echo Building Go Project...
	@CGO_ENABLED=1 GOOS=linux go build -o bin/membrane main.go

install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	@rm -rf ./bin/

test: install-tools
	@echo Running tests...
	@go run github.com/onsi/ginkgo/ginkgo -cover ./membrane/...

generate-proto:
	@echo Generating Proto Sources
	@mkdir -p ./interfaces/
	@protoc --go_out=./interfaces/ --go-grpc_out=./interfaces/ -I ./contracts/proto/ ./contracts/proto/**/*.proto

aws-plugins:
	@echo Building AWS plugins
	@go build -buildmode=plugin -o lib/documents/dynamodb.so ./plugins/aws/documents/dynamodb.go
	@go build -buildmode=plugin -o lib/eventing/sns.so ./plugins/aws/eventing/sns.go
	@go build -buildmode=plugin -o lib/gateway/lambda.so ./plugins/aws/gateway/lambda.go
	@go build -buildmode=plugin -o lib/storage/s3.so ./plugins/aws/storage/s3.go

aws-docker-debian:
	@docker build . -f ./plugins/aws/debian.dockerfile -t nitric:membrane-debian-aws

build-docker-alpine:
	@docker build . -f alpine.dockerfile -t nitric:membrane-alpine
build-docker-debian:
	@docker build . -f debian.dockerfile -t nitric:membrane-debian

membrane-docker: generate-proto build-docker-alpine build-docker-debian
	@echo Built Docker Images