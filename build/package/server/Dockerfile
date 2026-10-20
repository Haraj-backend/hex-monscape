FROM golang:1.27.3-alpine3.24 AS build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

RUN go install github.com/air-verse/air@v1.67.4
RUN mkdir -p .output

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd/server ./cmd/server
COPY ./build/package/server ./build/package/server
COPY ./internal ./internal

ENTRYPOINT ["air", "-c", "./build/package/server/.air.toml"]
