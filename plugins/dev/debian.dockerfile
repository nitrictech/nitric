# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building
FROM golang:buster as build

WORKDIR /

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make local-plugins

# Build the default development membrane server
FROM nitric:membrane-debian

# FIXME: Build these in a build stage during the docker build
# for now will just be copied post local build
# and execute these stages through a local shell script
COPY --from=build ./lib/documents.so /plugins/documents.so
RUN chmod +rx /plugins/documents.so
COPY --from=build ./lib/eventing.so /plugins/eventing.so
RUN chmod +rx /plugins/eventing.so
COPY --from=build ./lib/storage.so /plugins/storage.so
RUN chmod +rx /plugins/storage.so
COPY --from=build ./lib/gateway.so /plugins/gateway.so
RUN chmod +rx /plugins/gateway.so

ENTRYPOINT [ "/membrane" ]