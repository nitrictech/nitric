# Build the deployment binary
FROM golang:buster as build

WORKDIR /app/

# Optimize this to copy go workspace mod/sum files only
COPY . .

WORKDIR cloud/gcp

RUN make deploybin

FROM nitric-deploy-common AS base

LABEL "repository"="https://github.com/nitrictech/nitric"
LABEL "homepage"="https://nitric.io"
LABEL "maintainer"="Nitric Team <team@nitric.io>"
LABEL org.opencontainers.image.description="Nitric gcp deployer, in a Docker container."

# Copy the deploy binary from the build step
COPY --from=build app/cloud/gcp/bin/deploy-gcp /usr/bin

# Install deps all in one step
RUN curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
  # Add additional apt repos all at once
  echo "deb http://packages.cloud.google.com/apt cloud-sdk-$(lsb_release -cs) main" | tee /etc/apt/sources.list.d/google-cloud-sdk.list && \
  # Install second wave of dependencies
  apt-get update -y && \
  apt-get install -y \
  google-cloud-sdk \
  google-cloud-sdk-gke-gcloud-auth-plugin && \
  # Clean up the lists work
  rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["deploy-gcp"]
