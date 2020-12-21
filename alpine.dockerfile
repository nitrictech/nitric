# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building

FROM golang:alpine as build

RUN apk update
RUN apk upgrade

# RUN apk add --update go=1.8.3-r0 gcc=6.3.0-r4 g++=6.3.0-r4
RUN apk add --no-cache git curl gcc g++ make
# build-base autoconf automake libtool

WORKDIR /

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make

FROM alpine

COPY --from=build /bin/membrane /membrane
RUN chmod +xr /membrane

ENTRYPOINT [ "/membrane" ]

