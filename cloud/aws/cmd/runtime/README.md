# Nitric plugins for AWS

## Development

### Requirements
 - Git
 - Nitric Project
 - Golang
 - Make
 - Docker

### Getting Started

### Building Static Server Image
From the repository root run
```bash
make aws-docker-static
```

### Building Plugin Images

> __Note:__ Prior to building these plugins, the nitric pluggable server image must be built for local development


Alpine Linux
```bash
make aws-docker-alpine
```

Debian
```bash
make aws-docker-debian
```

> __Note:__ Separate distributions required between glibc/musl as dynamic linker is used for golang plugin support


