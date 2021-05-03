# Copyright 2021 Nitric Pty Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

