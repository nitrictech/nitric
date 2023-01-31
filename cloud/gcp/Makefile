ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

binaries: deploybin

# build runtime binary directly into the deploy director so it can be embedded directly into the deployment engine binary
runtimebin:
	@echo Building GCP Runtime Server
	@CGO_ENABLED=0 go build -o bin/runtime-gcp -ldflags="-extldflags=-static" ./runtime/cmd

# FIXME: proto server registration error
# There appears to be an old namespace conflict with the protobuf definitions
deploybin: runtimebin
	@echo Building GCP Deployment Server
	@cp bin/runtime-gcp deploy/runtime-gcp
	@CGO_ENABLED=0 go build -o bin/deploy-gcp -ldflags="-extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./deploy/cmd
	@rm deploy/runtime-gcp

dockerimage: deploybin
	@echo Building AWS Deployment Server Docker Image
	@docker build -t nitrictech/deploy-gcp:latest -f deploy/Dockerfile bin/

license-check: runtimebin
	@echo Checking GCP Membrane OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/runtime-gcp

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch deploy/runtime-aws
	$(GOLANGCI_LINT) run --fix
	@rm deploy/runtime-aws

lint:
	@touch deploy/runtime-gcp
	@go run github.com/google/addlicense -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch deploy/runtime-aws
	$(GOLANGCI_LINT) run

test: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./runtime/...

test-coverage: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./runtime/...

# generate mock implementations
generate-mocks:
	@echo Generating Mock Clients
	@mkdir -p mocks/gcp_storage
	@mkdir -p mocks/gcp_secret
	@mkdir -p mocks/pubsub
	@mkdir -p mocks/cloudtasks
	@mkdir -p mocks/provider
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/runtime/core GcpProvider > mocks/provider/gcp.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/gcloud_storage Reader,Writer,ObjectHandle,BucketHandle,BucketIterator,StorageClient,ObjectIterator > mocks/gcp_storage/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub PubsubClient,TopicIterator,Topic,PublishResult > mocks/pubsub/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/cloudtasks CloudtasksClient > mocks/cloudtasks/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/gcloud_secret SecretManagerClient,SecretIterator > mocks/gcp_secret/mock.go

generate-sources: generate-mocks
