FROM node:lts-alpine3.15 as frontend_build
WORKDIR /web
COPY ./cmd/web .

RUN yarn install

# FRONTEND_MODE define .env | .env.staging | .env.production that will be used
ARG FRONTEND_MODE="production"
ARG VITE_API_STAGE_PATH
ARG VITE_POKEBATTLE_URL

RUN yarn build:${FRONTEND_MODE}

FROM golang:1.18-alpine3.15 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/rest-lambda-ddb ./cmd/rest-lambda-ddb
COPY ./internal ./internal

WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/rest-lambda-ddb
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/go:1.2022.02.01.10

COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/rest-lambda-ddb/server ${LAMBDA_TASK_ROOT}
COPY --from=frontend_build /web/dist /dist

CMD [ "server" ]
