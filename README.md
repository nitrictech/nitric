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
 
## Running Locally

To run the membrane server locally, perform a local build of the membrane binary for the platform you're targeting, then run the resulting binary.

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

# Set Env, then Run the membrane binary
GATEWAY_ENVIRONMENT=http; ./bin/membrane
```

 