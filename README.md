# Nitric Membrane

## Architecture

## Development

### Requirements
 - Git
 - Golang
 - Make
 - Docker

### Getting Started

#### Install depdencies
```bash
make install-tools
```

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