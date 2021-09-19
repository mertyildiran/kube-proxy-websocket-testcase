FROM golang:1.16-alpine AS builder

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

WORKDIR /app/build

COPY go.mod .
COPY go.sum .
COPY server.go .
RUN go mod download
RUN go build -o server server.go

FROM alpine:3.13.5

WORKDIR /app

COPY --from=builder ["/app/build/server", "."]

ENV GIN_MODE=release

ENTRYPOINT "/app/server"
