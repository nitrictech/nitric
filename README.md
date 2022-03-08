<p align="center">
  <a href="https://nitric.io">
    <img src="docs/assets/nitric-logo.svg" width="120" alt="Nitric Logo"/>
  </a>
</p>

<p align="center">
  A fast & fun way to build portable cloud-native applications
</p>

<p align="center">
  <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/nitrictech/nitric?sort=semver">
  <img alt="GitHub" src="https://img.shields.io/github/license/nitrictech/nitric">
  <!-- <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/nitrictech/cli/total"> -->
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/nitrictech/nitric/Tests?label=build">
  <img alt="codecov" src="https://codecov.io/gh/nitrictech/nitric/branch/develop/graph/badge.svg?token=20TYFIQS2P">
  <!-- <a href="" target="_blank"><img src="https://img.shields.io/badge/discord-online-brightgreen.svg" alt="Discord"/></a> -->
  <a href="https://twitter.com/nitric_io">
    <img alt="Twitter Follow" src="https://img.shields.io/twitter/follow/nitric_io?label=Follow&style=social">
  </a>
</p>

## About Nitric

[Nitric](https://nitric.io) is a framework for rapid development of cloud-native and serverless applications. Define your apps in terms of the resources they need, then write the code for serverless function based APIs, event subscribers and scheduled jobs.

Apps built with Nitric can be deployed to AWS, Azure or Google Cloud all from the same code base so you can focus on your products, not your cloud provider.

Nitric makes it easy to:

- Create smart serverless functions and APIs
- Build reliable distributed apps that use events and/or queues
- Securely store, retrieve and rotate secrets
- Read and write files from buckets

## Documentation

The full documentation is available at [nitric.io/docs](https://nitric.io/docs).

We're completely opensource and encourage [code contributions](https://nitric.io/docs/contributions).

## Status

Nitric is currently in Public Preview. Anyone can use or deploy applications, but work remains and changes are likely. We’d love your feedback as we build additional functionality!

## Get in touch

- Ask questions in [GitHub discussions](https://github.com/nitrictech/nitric/discussions)

- Find us on [Twitter](https://twitter.com/nitric_io)

- Send us an [email](mailto:maintainers@nitric.io)

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

| Directory               | Package    | Description                                                |
| ----------------------- | ---------- | ---------------------------------------------------------- |
| `/interfaces/nitric/v1` | `v1`       | protoc generated GRPC services code                        |
| `/pkg/adapters/grpc`    | `grpc`     | GRPC service to SDK adaptors                               |
| `/pkg/membrane`         | `membrane` | membrane application                                       |
| `/pkg/plugins/...`      | `...`      | Cloud service SDK plugins                                  |
| `/pkg/providers/...`    | `main`     | Cloud provider main application and plugin injection       |
| `/pkg/sdk`              | `sdk`      | SDK service interfaces                                     |
| `/pkg/triggers`         | `triggers` | provides Nitric event triggers                             |
| `/pkg/utils`            | `utils`    | provides utility functions                                 |
| `/pkg/worker`           | `worker`   | Membrane workers representing function/service connections |
| `/tests/mocks/...`      | `...`      | Cloud service SDK mocks                                    |
| `/tests/plugins/...`    | `...`      | Plugin services integration test suites                    |
| `/tools`                | `tools`    | include for 3rd party build tools                          |

## Ideas for using Nitric

Examples of projects built with Nitric:

- Identity verification APIs
- FinTech APIs with complex business rules
- Move from Express.js or koa to serverless and distributed architectures
- Migrate from on-prem to the cloud (AWS, Azure or GCP)

Even though your apps are portable across clouds, they'll still make the best use of the fully-managed and serverless offerings of each cloud provider. Nitric deploys using services like Lambda, CloudRun, DynamoDB, FireStore, CosmosDB, SNS, Event Grid, PubSub... the list is super long. Nitric also makes sure IAM and other access is correctly configured in your deployed applications, so everything stays secure and just works.
