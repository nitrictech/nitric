ARG BASE_IMAGE

# Wrap any base image in this runtime wrapper
FROM ${BASE_IMAGE}

# ARG RUNTIME_URI
ARG RUNTIME_FILE

COPY ${RUNTIME_FILE} /bin/runtime
RUN chmod +x-rw /bin/runtime

# Inject original wrapped command here
CMD [%s]
ENTRYPOINT ["/bin/runtime"]
