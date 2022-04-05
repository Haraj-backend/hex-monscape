FROM node:lts-alpine3.15 as frontend_build
WORKDIR /web
COPY ./cmd/server/web .

# FRONTEND_MODE define .env | .env.staging | .env.production that will be used
ARG FRONTEND_MODE="production"
ARG VITE_API_STAGE_PATH
ARG VITE_POKEBATTLE_URL

RUN yarn install
RUN yarn build:${FRONTEND_MODE}

FROM golang:1.17.7-alpine3.15 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/lambda-mysql ./cmd/lambda-mysql
COPY ./internal ./internal

WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/lambda-mysql
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

# copy artifacts to a clean image
FROM alpine:3.12
RUN apk update \
    && apk add \
    ca-certificates \
    tzdata \
    mysql-client \
    build-base

COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/lambda-mysql/server .
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/lambda-mysql/entrypoint.sh .
RUN /bin/chmod +x entrypoint.sh
COPY --from=frontend_build /web/dist /dist

CMD [ "./server" ]
ENTRYPOINT ["./entrypoint.sh"]
