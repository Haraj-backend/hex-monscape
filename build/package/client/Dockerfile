FROM node:18.16.1-alpine3.18
WORKDIR /client

COPY ./cmd/client .
RUN yarn install
ENTRYPOINT ["yarn", "dev", "--host", "--port", "8161"]