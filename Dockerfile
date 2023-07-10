FROM golang:1.20.5 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 go build -o /rcc ./cmd/rcc

FROM alpine:3.17.0

COPY --from=builder /rcc /usr/local/bin/rcc

RUN mkdir -p /data

# Default store file system here for easy volume mounts
ENV RCC_STORAGE_TYPE=file
ENV RCC_FILE_STORAGE_BASE_DIR=/data

ENTRYPOINT ["/usr/local/bin/rcc"]
