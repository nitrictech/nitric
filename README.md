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
  <a href="https://discord.gg/Webemece5C"><img alt="Discord" src="https://img.shields.io/discord/955259353043173427?label=discord"></a>
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

- Join us on [Discord](https://discord.gg/Webemece5C)

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
- Google Protocol Buffers Compiler (protoc)

### Getting Started

#### Install dependencies

```bash
make install-tools
```

### Run unit tests

```bash
make test
```

### Run integration tests

```bash
make test-integration
```

### Building

#### Standard Runtime Binaries

> Linux support only - used in container images and for production.

```bash
make binaries
```

##### Running without a child process

It can be useful to run the Membrane in a 'service only' mode, where the cloud APIs are available but you don't need/want to start a child process to handle incoming request. This can be achieved by setting the MIN_WORKERS variable to `0`:

(export MIN_WORKERS=0; ./bin/membrane)

## Project Structure

The Membrane project source code structure is outlined below:

| Directory            | Description                                                |
| -------------------- | ---------------------------------------------------------- |
| `/core`              | Nitric core interfaces/contracts                           |
| `/provider/common`   | Nitric provider common module                              |
| `/provider/aws`      | Nitric AWS provider                                        |
| `/provider/gcp`      | Nitric GPC provider                                        |
| `/provider/azure`    | Nitric Azure provider                                      |
| `/e2e`               | E2E and integration testing module                         |
