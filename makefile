install:
	@echo installing go dependencies
	@go mod download

install-tools: install
	@echo Installing tools from tools.go

clean:
	@rm -rf ./bin/

build:
	@echo Building Go Project...
	@CGO_ENABLED=1 GOOS=linux go build -o bin/membrane main.go

build-docker-alpine:
	@docker build . -f alpine.dockerfile -t nitric:membrane-alpine --build-arg NITRIC_GITHUB_TOKEN=${NITRIC_GITHUB_TOKEN}
build-docker-debian:
	@docker build . -f debian.dockerfile -t nitric:membrane-debian --build-arg NITRIC_GITHUB_TOKEN=${NITRIC_GITHUB_TOKEN}

build-docker: build-docker-alpine build-docker-debian
	@echo Built Docker Images