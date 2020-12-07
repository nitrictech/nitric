# FIXME: Build with docker context in future
# Once the plugin SDK is visible for building

FROM golang:alpine as build
ARG NITRIC_GITHUB_TOKEN

RUN apk update
RUN apk upgrade

# RUN apk add --update go=1.8.3-r0 gcc=6.3.0-r4 g++=6.3.0-r4
RUN apk add --no-cache git gcc g++

WORKDIR /

COPY . .

# XXX: Setup private github repo access
RUN go env -w GOPRIVATE="github.com/nitric-dev"
RUN git config --global url.https://$NITRIC_GITHUB_TOKEN:x-oauth-basic@github.com/.insteadOf https://github.com/

RUN ./build.sh

FROM alpine

COPY --from=build /bin/membrane /membrane
RUN chmod +xr /membrane

ENTRYPOINT [ "/membrane" ]

