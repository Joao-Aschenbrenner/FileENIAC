# ENIACSYSTEMS/FileENIAC - Multi-stage build for Linux backend
# https://github.com/ENIACSystems/FileENIAC
# =============================================================================
FROM golang:1.26-alpine AS builder
WORKDIR /app

RUN apk add --no-cache musl-dev gcc sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY backend/ ./backend/
COPY main.go ./

ENV CGO_ENABLED=1
RUN go build -trimpath -ldflags="-s -w" -o /build/fileeniac ./backend/

# =============================================================================
FROM alpine:3.21 AS runtime
WORKDIR /fileeniac

RUN apk add --no-cache ca-certificates sqlite-libs libffi tzdata && \
    addgroup -g 1000 fileeniac && \
    adduser -u 1000 -G fileeniac -s /bin/sh -D fileeniac

COPY --from=builder /build/fileeniac /fileeniac/fileeniac
RUN chown fileeniac:fileeniac /fileeniac/fileeniac && chmod 0755 /fileeniac/fileeniac

USER fileeniac

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

VOLUME ["/data"]

ENV FILEENIAC_DATA_DIR=/data

ENTRYPOINT ["/fileeniac/fileeniac"]
CMD ["serve", "-a", ":8080"]
