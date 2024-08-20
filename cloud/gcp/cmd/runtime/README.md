# Nitric plugins for GCP

## Development

### Requirements
 - Git
 - Nitric Project
 - Golang
 - Make
 - Docker

### Getting Started

### Building Static Image
From the repository root run
```bash
make gcp-docker-static
```

### Building Plugin Images


> __Note:__ Prior to building these plugins, the nitric pluggable server image must be built for local development


Alpine Linux
```bash
make gcp-docker-alpine
```

Debian
```bash
make gcp-docker-debian
```

> __Note:__ Separate distributions required between glibc/musl as dynamic linker is used for golang plugin support

