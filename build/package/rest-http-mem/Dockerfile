FROM node:lts-alpine3.15 as frontend_build
WORKDIR /web

COPY ./cmd/web .

# FRONTEND_MODE define .env | .env.staging | .env.production that will be used
ARG FRONTEND_MODE="production" 

RUN yarn install
RUN yarn build:${FRONTEND_MODE}

FROM golang:1.17.7-alpine3.15 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/rest-http-mem
RUN go build -o server

FROM alpine:3.15
RUN apk add ca-certificates tzdata
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/rest-http-mem/server .
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/rest-http-mem/partners.json .
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/rest-http-mem/enemies.json .

# Copy compiled frontned to current directory
COPY --from=frontend_build /web/dist ./dist

ENTRYPOINT [ "./server" ]