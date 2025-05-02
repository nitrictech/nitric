ifeq (/,${HOME})
GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache/
else
GOLANGCI_LINT_CACHE=${HOME}/.cache/golangci-lint
endif
GOLANGCI_LINT ?= GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) go run github.com/golangci/golangci-lint/cmd/golangci-lint

binaries: deploybin

sec:
	@touch common/runtime/runtime-aws
	@go run github.com/securego/gosec/v2/cmd/gosec@latest -exclude-dir=tools ./...
	@rm common/runtime/runtime-aws

# build runtime binary directly into the deploy director so it can be embedded directly into the deployment engine binary
# We only build a linux amd64 binary here to be packaged for cloud runtimes with docker
# More binaries can be distributed in future if required
runtimebin:
	@echo Building AWS Runtime Server
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/runtime-aws -ldflags="-s -w -extldflags=-static" ./cmd/runtime

predeploybin: runtimebin
	@cp bin/runtime-aws common/runtime/runtime-aws

deploybin: predeploybin
	@echo Building AWS Deployment Server
	@CGO_ENABLED=0 go build -o bin/deploy-aws -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploy

deploybintf: generate-terraform predeploybin
	@echo Building AWS Terraform Deployment Server
	@CGO_ENABLED=0 go build -o bin/deploy-awstf -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploytf

install: deploybin deploybintf
	@echo installing aws deployment server to ${HOME}/.nitric/providers/nitric/aws-0.0.1
	@echo installing awstf deployment server to ${HOME}/.nitric/providers/nitric/awstf-0.0.1
	@mkdir -p ${HOME}/.nitric/providers/nitric/
	@if [ "$(OS)" == "Windows_NT" ]; then \
		rm -f "${HOME}/.nitric/providers/nitric/aws-0.0.1.exe"; \
		rm -f "${HOME}/.nitric/providers/nitric/awstf-0.0.1.exe"; \
		cp bin/deploy-aws "${HOME}/.nitric/providers/nitric/aws-0.0.1.exe"; \
		cp bin/deploy-awstf "${HOME}/.nitric/providers/nitric/awstf-0.0.1.exe"; \
	else \
		rm -f "${HOME}/.nitric/providers/nitric/aws-0.0.1"; \
		rm -f "${HOME}/.nitric/providers/nitric/awstf-0.0.1"; \
		cp bin/deploy-aws "${HOME}/.nitric/providers/nitric/aws-0.0.1"; \
		cp bin/deploy-awstf "${HOME}/.nitric/providers/nitric/awstf-0.0.1"; \
	fi

license-check: runtimebin
	@echo Checking AWS Runtime Server OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/runtime-aws

sourcefiles := $(shell find . -type f -name "*.go" -o -name "*.dockerfile")

fmt:
	@go run github.com/google/addlicense -ignore "./deploytf/generated/**" -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch common/runtime/runtime-aws
	$(GOLANGCI_LINT) run --fix
	@rm common/runtime/runtime-aws

lint:
	@go run github.com/google/addlicense -ignore "./deploytf/generated/**" -check -c "Nitric Technologies Pty Ltd." -y "2021" $(sourcefiles)
	@touch common/runtime/runtime-aws
	$(GOLANGCI_LINT) run
	@rm common/runtime/runtime-aws

test: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./runtime/...

test-coverage: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./runtime/...

clean-mocks:
	@rm -rf ./mocks

generate-mocks: clean-mocks
	@echo Generating Mock Clients
	@mkdir -p mocks/secrets_manager
	@mkdir -p mocks/s3
	@mkdir -p mocks/sns
	@mkdir -p mocks/sfn
	@mkdir -p mocks/sqs
	@mkdir -p mocks/provider
	@mkdir -p mocks/resourcetaggingapi
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/ifaces/resourcegroupstaggingapiiface ResourceGroupsTaggingAPIAPI > mocks/resourcetaggingapi/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/ifaces/snsiface SNSAPI > mocks/sns/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/ifaces/sfniface SFNAPI > mocks/sfn/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/ifaces/secretsmanageriface SecretsManagerAPI > mocks/secrets_manager/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/ifaces/s3iface S3API,PreSignAPI > mocks/s3/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/ifaces/sqsiface SQSAPI > mocks/sqs/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/cloud/aws/runtime/resource AwsResourceResolver > mocks/provider/aws.go

generate-terraform:
	@cd deploytf && npx -y cdktf-cli@0.20.12 get

generate-sources: generate-mocks

tidy:
	@go mod tidy

.PHONY: binaries sec deploybin deploybintf install license-check fmt lint test test-coverage clean-mocks generate-mocks generate-terraform generate-sources tidy