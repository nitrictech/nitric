# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building

FROM golang:buster as build
ARG NITRIC_GITHUB_TOKEN

WORKDIR /

# XXX: Setup private github repo access
RUN go env -w GOPRIVATE="github.com/nitric-dev"
RUN git config --global url.https://$NITRIC_GITHUB_TOKEN:x-oauth-basic@github.com/.insteadOf https://github.com/

# Cache dependencies in seperate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ./build.sh

FROM debian:buster-slim

COPY --from=build /bin/membrane /membrane
RUN chmod +xr /membrane

ENTRYPOINT [ "/membrane" ]

