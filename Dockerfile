ARG GO_VERSION=1.22.2

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk add --no-cache add ga-certificates && update-ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /dist/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /dist/app /app/
