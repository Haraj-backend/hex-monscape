FROM golang:1.20.5-alpine3.18 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

COPY go.mod go.sum ./
RUN go mod download -x

RUN mkdir -p .cov-report

COPY ./internal ./internal

CMD go test -v -coverprofile=./.cov-report/profile.cov ./...
