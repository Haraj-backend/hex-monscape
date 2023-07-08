# build client
FROM node:18.16.1-alpine3.18 as client_build
WORKDIR /client

COPY ./cmd/client .
RUN yarn install

# run the build script, the build result will be in ./dist
RUN yarn build:production

# build server
FROM golang:1.20.5-alpine3.18 as server_build
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd/lambda ./cmd/lambda
COPY ./internal ./internal
WORKDIR /go/src/github.com/Haraj-backend/hex-monscape/cmd/lambda
RUN CGO_ENABLED=0 GOOS=linux go build -o server

# copy both client & server artifacts to a clean image
FROM public.ecr.aws/lambda/go:1.2023.07.05.11
WORKDIR ${LAMBDA_TASK_ROOT}
COPY --from=server_build /go/src/github.com/Haraj-backend/hex-monscape/cmd/lambda/server ./server
COPY --from=client_build /client/dist ./client

CMD [ "server" ]