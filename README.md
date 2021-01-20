# Nitric Membrane

## Architecture

## Development

### Requirements
 - Git
 - Golang
 - Make
 - Docker
 - Google Protocol Buffers Compiler

### Getting Started

#### Install dependencies
```bash
make install-tools
```

##### Protocol Buffers
Download the Google Protobuf Compiler (standalone binary called `protoc`) from https://github.com/protocolbuffers/protobuf and add it to your $PATH.

> On Mac OS with Homebrew, you can run `brew install protobuf`

 

#### Run tests
```bash
make tests
```

#### Build Pluggable Membrane (images are used)
```bash
make
```

### Building Pluggable Membrane Images
Alpine Linux
```bash
make build-docker-alpine
```

Debian
```bash
make build-docker-debian
```

Or both
```bash
make build-docker
```

### For building statically compiled provider specific membranes see plugin README(s)

 - [AWS](./plugins/aws/README.md)
 - [GCP](./plugins/gcp/README.md)
 - [Dev](./plugins/dev/README.md)