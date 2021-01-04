# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building

FROM golang:buster as build

WORKDIR /

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make

FROM debian:buster-slim

COPY --from=build /bin/membrane /membrane
RUN chmod +xr /membrane

ENTRYPOINT [ "/membrane" ]

