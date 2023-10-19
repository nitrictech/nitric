# Copyright Nitric Pty Ltd.

# SPDX-License-Identifier: Apache-2.0

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG BASE_IMAGE

FROM ${BASE_IMAGE}

ARG RUNTIME_FILE

COPY ${RUNTIME_FILE} /bin/runtime
RUN chmod +x-rw /bin/runtime

ARG OTELCOL_CONTRIB_URI

ADD ${OTELCOL_CONTRIB_URI} /usr/bin/
RUN tar -xzf /usr/bin/otelcol*.tar.gz &&\
    rm /usr/bin/otelcol*.tar.gz &&\
	mv /otelcol-contrib /usr/bin/

ARG OTELCOL_CONFIG
RUN mkdir /etc/otelcol
RUN echo -e "${OTELCOL_CONFIG}" > /etc/otelcol/config.yaml
RUN chmod -R a+r /etc/otelcol

ARG NITRIC_TRACE_SAMPLE_PERCENT

ENV NITRIC_TRACE_SAMPLE_PERCENT ${NITRIC_TRACE_SAMPLE_PERCENT}
ENV OTELCOL_BIN /usr/bin/otelcol-contrib
ENV OTELCOL_CONFIG /etc/otelcol/config.yaml

CMD [%s]
ENTRYPOINT ["/bin/runtime"]
