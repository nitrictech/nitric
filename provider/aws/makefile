binaries: runtimebin

runtimebin:
	@echo Building AWS Runtime Server
	@CGO_ENABLED=0 go build -o bin/runtime-aws -ldflags="-extldflags=-static" ./runtime/cmd

license-check: runtimebin
	@echo Checking AWS Membrane OSS Licenses
	@go run github.com/uw-labs/lichen --config=./lichen.yaml ./bin/runtime-aws

# TODO: 
# deploybin:
# 	@echo Building Deployment Engine

test: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo ./runtime/...

test-coverage: generate-mocks
	@echo Running unit tests
	@go run github.com/onsi/ginkgo/ginkgo -cover -outputdir=./ -coverprofile=all.coverprofile ./runtime/...

generate-mocks:
	@echo Generating Mock Clients
	@mkdir -p mocks/secrets_manager
	@mkdir -p mocks/s3
	@mkdir -p mocks/sns
	@mkdir -p mocks/sfn
	@mkdir -p mocks/sqs
	@mkdir -p mocks/provider
	@mkdir -p mocks/resourcetaggingapi
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/ifaces/resourcegroupstaggingapiiface ResourceGroupsTaggingAPIAPI > mocks/resourcetaggingapi/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/ifaces/snsiface SNSAPI > mocks/sns/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/ifaces/sfniface SFNAPI > mocks/sfn/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/ifaces/secretsmanageriface SecretsManagerAPI > mocks/secrets_manager/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/ifaces/s3iface S3API,PreSignAPI > mocks/s3/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/ifaces/sqsiface SQSAPI > mocks/sqs/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/nitric/provider/aws/runtime/core AwsProvider > mocks/provider/aws.go

generate-sources: generate-mocks