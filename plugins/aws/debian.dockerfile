# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building
FROM golang:buster as build

WORKDIR /

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make aws-plugins

# Default dockerfile for the Nitric AWS base image
FROM nitric:membrane-debian

# FIXME: Build these in a build stage during the docker build
# for now will just be copied post local build
# and execute these stages through a local shell script
COPY --from=build ./lib/documents/dynamodb.so /plugins/documents.so
COPY --from=build ./lib/eventing/sns.so /plugins/eventing.so
COPY --from=build ./lib/storage/s3.so /plugins/storage.so
COPY --from=build ./lib/gateway/lambda.so /plugins/gateway.so

RUN chmod +rx /plugins/*

# FIXME: Do we need this here?
ENTRYPOINT [ "/membrane" ]