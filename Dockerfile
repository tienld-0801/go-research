FROM golang:1.23-alpine

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api/tmp/main ./cmd

CMD ["/api/tmp/main"]
