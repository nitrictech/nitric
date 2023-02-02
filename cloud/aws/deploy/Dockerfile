# Build the deployment binary
FROM golang:buster as build

WORKDIR /app/

# Optimize this to copy go workspace mod/sum files only
COPY . .

WORKDIR cloud/aws

RUN make deploybin

FROM nitric-deploy-common AS base

LABEL "repository"="https://github.com/nitrictech/nitric"
LABEL "homepage"="https://nitric.io"
LABEL "maintainer"="Nitric Team <team@nitric.io>"
LABEL org.opencontainers.image.description="Nitric aws deployer, in a Docker container."

# Copy the deploy binary from the build step
COPY --from=build app/cloud/aws/bin/deploy-aws /usr/bin

# Install deps all in one step
RUN curl -fsSLo /usr/bin/aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.10.3/2018-07-26/bin/linux/amd64/aws-iam-authenticator && \
  chmod +x /usr/bin/aws-iam-authenticator && \
  # AWS v2 cli
  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
  unzip awscliv2.zip && \
  ./aws/install && \
  rm -rf aws && \
  rm awscliv2.zip && \
  # Add additional apt repos all at once
  echo "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"      | tee /etc/apt/sources.list.d/docker.list           && \
  # Install second wave of dependencies
  apt-get install -y \
  docker-ce && \
  # Clean up the lists work
  rm -rf /var/lib/apt/lists/*

RUN pulumi plugin install resource aws 5.26.0

ENTRYPOINT ["deploy-aws"]
