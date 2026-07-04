FROM golang:1.23-alpine AS builder
ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /src
COPY go.mod ./
COPY agent-gateway-go/ ./agent-gateway-go/
COPY device-gateway-go/ ./device-gateway-go/
COPY cmd/ ./cmd/
RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -trimpath -ldflags='-s -w' -o /out/gateway ./cmd/gateway

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
EXPOSE 8787 8788
ENV AGENT_PORT=8787
ENV DEVICE_PORT=8788
COPY --from=builder /out/gateway /usr/local/bin/gateway
ENTRYPOINT ["gateway"]
