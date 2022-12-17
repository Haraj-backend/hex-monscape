FROM golang:1.18-alpine3.15 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./internal ./internal