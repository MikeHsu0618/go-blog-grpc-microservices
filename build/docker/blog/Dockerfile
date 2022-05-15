# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.18-alpine as builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /workspace


COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api/ api/
COPY cmd/ cmd/
COPY configs/ configs/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o app /workspace/cmd/blog/

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /

COPY --from=builder /workspace/configs /configs
COPY --from=builder /workspace/app .

# 加入 grpc 探針
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.10 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/$GRPC_HEALTH_PROBE_VERSION/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

EXPOSE 50054
EXPOSE 8054
EXPOSE 9054
ENTRYPOINT ["/app"]