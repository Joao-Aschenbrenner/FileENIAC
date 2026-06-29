FROM golang:1.26-alpine

RUN apk add --no-cache git musl-dev gcc sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend/ ./backend/

RUN addgroup -g 1000 fileeniac && \
    adduser -u 1000 -G fileeniac -s /bin/sh -D fileeniac
USER fileeniac

CMD ["go", "run", "./backend"]
