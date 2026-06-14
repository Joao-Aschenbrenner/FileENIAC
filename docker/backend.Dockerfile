FROM golang:1.21-alpine

RUN apk add --no-cache git musl-dev gcc sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY backend/ ./backend/

CMD ["go", "run", "./backend"]
