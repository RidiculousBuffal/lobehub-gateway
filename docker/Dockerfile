FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod ./
COPY cmd/ ./cmd/
COPY agent-gateway-go/ ./agent-gateway-go/
COPY device-gateway-go/ ./device-gateway-go/

RUN go build -trimpath -ldflags='-s -w' -o /out/gateway ./cmd/gateway

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
EXPOSE 8787 8788
ENV AGENT_PORT=8787
ENV DEVICE_PORT=8788
COPY --from=builder /out/gateway /usr/local/bin/gateway
ENTRYPOINT ["gateway"]
