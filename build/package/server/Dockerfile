FROM node:lts-alpine3.15 as frontend_build
WORKDIR /web

COPY ./cmd/server/web .

# FRONTEND_MODE define .env | .env.staging | .env.production that will be used
ARG FRONTEND_MODE="production"

RUN yarn install
RUN VITE_POKEBATTLE_URL=http://localhost:3000 yarn build:${FRONTEND_MODE}

FROM golang:1.17.7-alpine3.15 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/server
RUN go build -o server

FROM alpine:3.15
RUN apk add ca-certificates tzdata
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/server/server .
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/server/partners.json .
COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/server/enemies.json .

# Copy compiled frontned to current directory
COPY --from=frontend_build /web/dist ./dist

ENTRYPOINT [ "./server" ]
