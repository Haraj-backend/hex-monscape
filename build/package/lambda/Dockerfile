# build client
FROM node:24.22.0-alpine3.24 AS client_build
WORKDIR /client

COPY ./cmd/client .
RUN yarn install

# run the build script, the build result will be in ./dist
RUN yarn build:production

# build server
FROM golang:1.27.3-alpine3.24 AS server_build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd/lambda ./cmd/lambda
COPY ./internal ./internal
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape/cmd/lambda
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap

# copy both client & server artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023.2026.10.14.18
WORKDIR ${LAMBDA_TASK_ROOT}
COPY --from=server_build /go/src/github.com/Haraj-backend/hex-monscape/cmd/lambda/bootstrap ./bootstrap
COPY --from=client_build /client/dist ./client

ENTRYPOINT [ "./bootstrap" ]