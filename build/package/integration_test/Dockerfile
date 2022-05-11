FROM golang:1.17.7-alpine3.15 as build
WORKDIR /go/src/github.com/ghazlabs/hex-mathrush

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./internal ./internal