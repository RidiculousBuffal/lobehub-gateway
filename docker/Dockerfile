FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
EXPOSE 8787 8788
ENV AGENT_PORT=8787
ENV DEVICE_PORT=8788
ARG BINARY_FILE=gateway-linux-amd64
COPY ${BINARY_FILE} /usr/local/bin/gateway
ENTRYPOINT ["gateway"]
