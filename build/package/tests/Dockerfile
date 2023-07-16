FROM golang:1.20.5-alpine3.18 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

COPY go.mod go.sum ./
RUN go mod download -x

RUN mkdir -p .cov-report

COPY ./internal ./internal

CMD go test ./... -covermode=count -coverprofile=./.cov-report/coverage.out fmt && go tool cover -func=./.cov-report/coverage.out -o=./.cov-report/coverage.out