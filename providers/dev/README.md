# Nitric plugins for Local Development

## Development

### Requirements
 - Git
 - Nitric Membrane Project
 - Golang
 - Make
 - Docker

### Getting Started

### Building Static Membrane Image
From the repository root run
```bash
make local-docker-static
```

### Building Images

> __Note:__ Prior to building these plugins, the nitric pluggable membrane image must be built for local development


Alpine Linux
```bash
make local-docker-alpine
```

Debian
```bash
make local-docker-debian
```

> __Note:__ Seperate distributions required between glibc/musl as dynamic linker is used for golang plugin support


