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
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/runtime-gcp -ldflags="-s -w -extldflags=-static" ./cmd/runtime

predeploybin: runtimebin
	@cp bin/runtime-gcp common/runtime/runtime-gcp

sec:
	@touch common/runtime/runtime-gcp
	@go run github.com/securego/gosec/v2/cmd/gosec@latest -exclude-dir=tools ./...
	@rm common/runtime/runtime-gcp

# There appears to be an old namespace conflict with the protobuf definitions
deploybin: predeploybin
	@echo Building GCP Deployment Server
	@CGO_ENABLED=0 go build -o bin/deploy-gcp -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploy


deploybintf: generate-terraform predeploybin
	@echo Building GCP Terraform Deployment Server
	@CGO_ENABLED=0 go build -o bin/deploy-gcptf -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploytf

install: deploybin deploybintf
	@echo installing gcp deployment server to ${HOME}/.nitric/providers/nitric/gcp-0.0.1
	@echo installing gcptf deployment server to ${HOME}/.nitric/providers/nitric/gcptf-0.0.1
	@mkdir -p ${HOME}/.nitric/providers/nitric/
	@if [ "$(OS)" == "Windows_NT" ]; then \
		rm -f "${HOME}/.nitric/providers/nitric/gcp-0.0.1.exe"; \
		rm -f "${HOME}/.nitric/providers/nitric/gcptf-0.0.1.exe"; \
		cp bin/deploy-gcp "${HOME}/.nitric/providers/nitric/gcp-0.0.1.exe"; \
		cp bin/deploy-gcptf "${HOME}/.nitric/providers/nitric/gcptf-0.0.1.exe"; \
	else \
		rm -f "${HOME}/.nitric/providers/nitric/gcp-0.0.1"; \
		rm -f "${HOME}/.nitric/providers/nitric/gcptf-0.0.1"; \
		cp bin/deploy-gcp "${HOME}/.nitric/providers/nitric/gcp-0.0.1"; \
		cp bin/deploy-gcptf "${HOME}/.nitric/providers/nitric/gcptf-0.0.1"; \
	fi

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -ignore "./deploytf/generated/**" -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch common/runtime/runtime-gcp
	$(GOLANGCI_LINT) run --fix
	@rm common/runtime/runtime-gcp

lint:
	@go run github.com/google/addlicense -ignore "./deploytf/generated/**" -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch common/runtime/runtime-gcp
	$(GOLANGCI_LINT) run
	@rm common/runtime/runtime-gcp

license-check: runtimebin
	@echo Checking GCP Runtime Server OSS Licenses
	@go tool lichen --config=./lichen.yaml ./bin/runtime-gcp

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
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/runtime/resource GcpResourceResolver > mocks/provider/gcp.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/gcloud_storage Reader,Writer,ObjectHandle,BucketHandle,BucketIterator,StorageClient,ObjectIterator > mocks/gcp_storage/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub PubsubClient,TopicIterator,Topic,PublishResult > mocks/pubsub/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/cloudtasks CloudtasksClient > mocks/cloudtasks/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/gcp/ifaces/gcloud_secret SecretManagerClient,SecretIterator > mocks/gcp_secret/mock.go

generate-sources: generate-mocks

generate-terraform:
	@cd deploytf && npx -y cdktf-cli@0.20.10 get

tidy:
	@go mod tidy

.PHONY: binaries runtimebin predeploybin deploybin deploybintf install fmt lint license-check test test-coverage generate-mocks generate-sources generate-terraform tidy