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
FROM alpine:3.19 AS runtime
WORKDIR /fileeniac

RUN apk add --no-cache ca-certificates sqlite-libs libffi tzdata

COPY --from=builder /build/fileeniac /fileeniac/fileeniac

EXPOSE 8080

VOLUME ["/data"]

ENV FILEENIAC_DATA_DIR=/data

ENTRYPOINT ["/fileeniac/fileeniac"]
CMD ["serve", "-a", ":8080"]
