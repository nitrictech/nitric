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

#### Build locally (images are used)
```bash
make build
```

### Building Images
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
