# Nitric plugins for Microsoft Azure

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
make azure-docker-static
```

### Building Plugin Images

> __Note:__ Prior to building these plugins, the nitric pluggable membrane image must be built for local development


Alpine Linux
```bash
make azure-docker-alpine
```

Debian
```bash
make azure-docker-debian
```

> __Note:__ Separate distributions required between glibc/musl as dynamic linker is used for golang plugin support

### Credentials
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_TENANT_ID

AZURE_SUBSCRIPTION_ID
AZURE_RESOURCE_GROUP

#### Key Vault
AZURE_VAULT_NAME

