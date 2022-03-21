FROM node:lts-alpine3.15 as frontend_build
WORKDIR /web

COPY ./cmd/server/web .

# FRONTEND_MODE define .env | .env.staging | .env.production that will be used
ARG FRONTEND_MODE="production"

RUN yarn install

ARG VITE_POKEBATTLE_URL
RUN yarn build:${FRONTEND_MODE}

FROM golang:1.17.7-alpine3.15 as build
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/lambda ./cmd/lambda
COPY ./internal ./internal
WORKDIR /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/lambda
RUN CGO_ENABLED=0 GOOS=linux go build -o server

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/go:1.2022.02.01.10

COPY --from=build /go/src/github.com/Haraj-backend/hex-pokebattle/cmd/lambda/server /var/task
COPY --from=frontend_build /web/dist /dist

CMD [ "./server" ]
