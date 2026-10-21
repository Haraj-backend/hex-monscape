FROM golang:1.27.3-alpine3.24 AS build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

COPY go.mod go.sum ./
RUN go mod download -x

RUN mkdir -p .cov-report

COPY ./internal ./internal

# since Go 1.22, packages without test files are counted as 0% coverage,
# so we only run tests for packages that have test files
CMD go test -v -coverprofile=./.cov-report/profile.cov $(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)
