# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building
FROM golang as build

WORKDIR /

RUN apt-get update && apt-get install -y protobuf-compiler

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
COPY makefile makefile
COPY tools/tools.go tools/tools.go
RUN make install-tools

COPY . .

RUN make gcp-static

# Build the default development membrane server
FROM alpine
# FIXME: Build these in a build stage during the docker build
# for now will just be copied post local build
# and execute these stages through a local shell script
COPY --from=build ./bin/membrane /membrane
RUN chmod +rx /membrane

ENTRYPOINT [ "/membrane" ]