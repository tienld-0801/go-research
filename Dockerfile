FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd

FROM alpine:latest

RUN apk add --no-cache netcat-openbsd
COPY --from=builder /app/main /api/tmp/main

WORKDIR /api/tmp

CMD ["./main"]
