# Build the deployment binary
FROM golang:buster as build

WORKDIR /app/

# Optimize this to copy go workspace mod/sum files only
COPY . .

WORKDIR cloud/azure

RUN make deploybin

FROM nitric-deploy-common

LABEL "repository"="https://github.com/nitrictech/nitric"
LABEL "homepage"="https://nitric.io"
LABEL "maintainer"="Nitric Team <team@nitric.io>"
LABEL org.opencontainers.image.description="Nitric azure deployer, in a Docker container."

COPY --from=build app/cloud/azure/bin/deploy-azure /usr/bin

# Install deps all in one step
RUN curl -fsSL https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
  # Add additional apt repos all at once
  echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/azure.list && \
  # Install second wave of dependencies
  apt-get update -y && \
  apt-get install -y \
  # Pin azure-cli to 2.33.1 as workaround for https://github.com/pulumi/pulumi-docker-containers/issues/106
  "azure-cli=2.33.1-1~bullseye" && \
  # Clean up the lists work
  rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["deploy-azure"]
