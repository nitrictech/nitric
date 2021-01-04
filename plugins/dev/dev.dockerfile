# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building
FROM golang:alpine as build

WORKDIR /

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git gcc g++ make

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make local-static

# Build the default development membrane server
FROM alpine
# FIXME: Build these in a build stage during the docker build
# for now will just be copied post local build
# and execute these stages through a local shell script
COPY --from=build ./bin/membrane /membrane
RUN chmod +rx /membrane

ENTRYPOINT [ "/membrane" ]