ARG BASE_IMAGE

# TODO: Need to make sure the architecture for the build matches the base image
FROM golang as base

ARG PLUGIN_DEFINITION
ENV PLUGIN_DEFINITION=${PLUGIN_DEFINITION}

# Need to install make
RUN apt-get update && apt-get install -y make

# Checkout the nitric github repo
RUN git clone --depth=1 -b gcp-cloud-run https://github.com/HomelessDinosaur/nitric /nitric
WORKDIR /nitric

RUN go work sync

WORKDIR /nitric/server

RUN make

FROM $BASE_IMAGE

ARG ORIGINAL_COMMAND
ENV ORIGINAL_COMMAND=${ORIGINAL_COMMAND}

COPY --from=base /nitric/server/bin/host /usr/local/bin/nitric

# CMD ["-c", "$ORIGINAL_COMMAND"]
ENTRYPOINT /usr/local/bin/nitric -c "$ORIGINAL_COMMAND"
