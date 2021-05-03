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