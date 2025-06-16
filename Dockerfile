FROM golang:1.24.4-alpine

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api/tmp/main ./cmd

CMD ["/api/tmp/main"]
