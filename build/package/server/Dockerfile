FROM golang:1.20.5-alpine3.18 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

RUN go install github.com/cosmtrek/air@v1.44.0
RUN mkdir -p .output

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd/server ./cmd/server
COPY ./build/package/server ./build/package/server
COPY ./internal ./internal

ENTRYPOINT ["air", "-c", "./build/package/server/.air.toml"]
