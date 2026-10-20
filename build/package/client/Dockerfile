FROM node:24.22.0-alpine3.24
WORKDIR /client

COPY ./cmd/client .
RUN yarn install
ENTRYPOINT ["yarn", "dev", "--host", "--port", "8161"]