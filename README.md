<p align="center">
  <img src="docs/assets/dot-matrix-logo.png" alt="Nitric Logo"/>
</p>

![test status](https://github.com/nitrictech/go-sdk/actions/workflows/test.yaml/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/nitrictech/nitric/branch/develop/graph/badge.svg?token=20TYFIQS2P)](https://codecov.io/gh/nitrictech/nitric)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=nitrictech_membrane&metric=alert_status)](https://sonarcloud.io/dashboard?id=nitrictech_membrane)

[Nitric](https://nitric.io) is a provider independent framework and runtime for cloud-native and serverless applications. Using Nitric, applications take advantage of cloud-native services for events, queues, compute, APIs, storage, and documents without direct integration to cloud specific APIs.

This decoupling enables applications to remain portable between cloud-providers and alternate deployment options such as Kubernetes or stand-alone servers, from a single application codebase.

## Documentation

The full documentation is available at [nitric.io/docs](https://nitric.io/docs)

## Status

Nitric is currently in Public Preview, anyone can use or deploy applications, but work remains and changes are likely

## Nitric Membrane

The Membrane is at the heart of the solution. Nitric applications communicate with the Membrane via gRPC to access the following services in a provider agnostic way:

  - Events
  - Queues
  - Storage & Buckets
  - Document Store
  - Secret Store

We provide an expressive infrastructure-as-code style SDK for [Node.js](https://github.com/nitrictech/node-sdk). However, Nitric is built on gRPC, so support for many languages is possible.
 
> If you have additional languages you'd like supported, let us know in the issues, we also welcome community contributions for new language support.

## Development

### Requirements
 - Git
 - Golang (1.16)
 - Make
 - Docker
 - Google Protocol Buffers Compiler

### Getting Started

#### Install dependencies
```bash
make install-tools
```

##### Install Protocol Buffers
Download the Google Protobuf Compiler (standalone binary called `protoc`) from https://github.com/protocolbuffers/protobuf and add it to your $PATH.

> On MacOS with Homebrew, you can run `brew install protobuf`
> On Fedora, run `sudo dnf install -y protobuf protobuf-compiler protobuf-devel`

### Run unit tests
```bash
make test
```
### Run integration tests
```bash
make test-integration
```

### Build Static Membranes

#### AWS

##### Standard Binary

> Linux support only - used in container images and for production.

```bash
make aws-static
```

##### Cross-platform Binary

Useful for local testing

```bash 
make aws-static-xp
```

##### Container Images

```bash
make aws-docker
```

#### Google Cloud Platform

##### Standard Binary

> Linux support only - used in container images and for production.

```bash
make gcp-static
```

##### Cross-platform Binary

Useful for local testing

```bash 
make gcp-static-xp
```

##### Container Images

```bash
make gcp-docker
```

#### Dev Membrane

> Note: the Dev Membrane should only be used for local development and testing.

##### Standard Binary

The dev binary is always cross-platform, since it doesn't need to be optimized for production deployments.

```bash
make dev-static
```

##### Container Images

```bash
make dev-docker
```


### Run Locally

To run the membrane server locally, perform a local build of the membrane binary for the platform you're targeting, then run the resulting binary.

##### Example building and running the static Google Cloud Membrane locally

```bash
# Make the GCP Static Cross-platform binary
make gcp-static-xp

# Run the membrane binary
./bin/membrane
```

> Note: for the AWS membrane, the Lambda Gateway (default) will fail to start. Instead, set the `GATEWAY_ENVIRONMENT` environment variable so that the HTTP gateway is launched instead.

```bash
# Make the AWS Static Cross-platform binary
make aws-static-xp

# Set environment variable in subshell, then run the membrane binary
(export GATEWAY_ENVIRONMENT=http; ./bin/membrane)
```

##### Running without a child process

It can be useful to run the Membrane in a 'service only' mode, where the cloud APIs are available but you don't need/want to start a child process to handle incoming request. This can be achieved by setting the MIN_WORKERS variable to `0`:

(export MIN_WORKERS=0; ./bin/membrane)

## Project Structure

The Membrane project source code structure is outlined below:

Directory                   | Package    | Description
---------                   |----------- |------------
`/interfaces/nitric/v1`     | `v1`       | protoc generated GRPC services code 
`/pkg/adapters/grpc`        | `grpc`     | GRPC service to SDK adaptors 
`/pkg/membrane`             | `membrane` | membrane application
`/pkg/plugins/...`          | `...`      | Cloud service SDK plugins 
`/pkg/providers/...`        | `main`     | Cloud provider main application and plugin injection 
`/pkg/sdk`                  | `sdk`      | SDK service interfaces 
`/pkg/triggers`             | `triggers` | provides Nitric event triggers
`/pkg/utils`                | `utils`    | provides utility functions
`/pkg/worker`               | `worker`   | Membrane workers representing function/service connections
`/tests/mocks/...`          | `...`      | Cloud service SDK mocks 
`/tests/plugins/...`        | `...`      | Plugin services integration test suites
`/tools`                    | `tools`    | include for 3rd party build tools