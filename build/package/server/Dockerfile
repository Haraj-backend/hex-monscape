FROM golang:1.20.5-alpine3.18 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd ./cmd
COPY ./build ./build
COPY ./internal ./internal

COPY cmd/rest-mem/partners.json /partners.json
COPY cmd/rest-mem/enemies.json /enemies.json

RUN go install github.com/cosmtrek/air@v1.27.3
RUN mkdir -p tmp

ENTRYPOINT ["air", "-c", "build/package/rest-mem/.air.toml"]
