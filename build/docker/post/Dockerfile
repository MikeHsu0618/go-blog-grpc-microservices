# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.18-alpine as builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /workspace
RUN apk add git
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go install github.com/grpc-ecosystem/grpc-health-probe@latest
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY api/ api/
COPY cmd/ cmd/
COPY configs/ configs/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o app /workspace/cmd/post/

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /

COPY --from=builder /workspace/configs /configs
COPY --from=builder /workspace/app .
COPY --from=builder /go/bin/grpc-health-probe .
EXPOSE 50054
EXPOSE 8054
EXPOSE 9054
ENTRYPOINT ["/app"]