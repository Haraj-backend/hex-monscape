# Hex Monscape Go

Welcome to this repo!

In this repo you will learn the concept of [Hexagonal Architecture](./docs/reference/hex-architecture.md) and its implementation through simple server-client game called `Hex Monscape`.

In the game you will play as a `10 years` old monster hunter that dreams to become the very best. In order to reach that, you need to make journey together with your monster partner to seek `3` strong wild monsters and kick them in the butt. 💥💪🏻

We are applying the concept of [Hexagonal Architecture](./docs/reference/hex-architecture.md) to implement the game server while coding it using [Golang](https://go.dev/). For the web client, we implemented it using [Vue 3](https://vuejs.org/).

To see the REST API specification for this game, please see [this doc](./docs/api-design/rest-api.md).

> **Note:**
>
> As Solutions Team member, your understanding towards [Hexagonal Architecture](./docs/reference/hex-architecture.md) is mandatory since it is the main architecture we used for building Haraj production systems.
>
> So if you understand this architecture well, you will be in no time contributing to Haraj production.
>
> Please refer to [Primary References](#primary-references) section to start learning about these concepts.

## How to Run The Game

When we are using [Hexagonal Architecture](./docs/reference/hex-architecture.md) to design a system, it is quite easy to swap its infrastructure code with another technologies.

So for example, if initially we used in memory storage to store our data, we could easily swap it with MySQL storage or something else.

To demonstrate this point, there are `3` variants of game server in this project:

- Server using Memory storage
- Server using DynamoDB storage
- Server using MySQL storage

All of them will serve the same game, the only difference is the place where they store the game data.

All of these servers could be run by using this command:

```bash
> make run
```

This command will create & run the stack defined in this [docker-compose.yml](./deploy/local/deployment/docker-compose.yml). 

Wait a moment until the entire stack is running. You will see something like this in the console after it is done:

```bash
hex_mem_1     | 2022/05/11 16:29:50 server is listening on :9185...
hex_mysql_1   | 2022/05/11 16:30:21 server is listening on :9186...
hex_ddb_1     | 2022/05/11 16:30:21 server is listening on :9187...
```

After that you could access each of these servers by visiting endpoints below:

- Memory storage => http://localhost:9185
- DynamoDB storage => http://localhost:9186
- MySQL storage => http://localhost:9187

## Primary References

To start learning the concept of [Hexagonal Architecture](./docs/reference/hex-architecture.md) please use [this doc](./docs/reference/hex-architecture.md) as your primary source of learning. This is so you won't be having too much confusion when learning it from other online resources.

To know more about the game design, please refer to [this doc](./docs/reference/game-design.md).

To learn about the methodology on how to implement [Hexagonal Architecture](./docs/reference/hex-architecture.md) on a project, please refer to [this doc](./docs/reference/project-methodology.md).

## Attribution

The monster characters used in this project is designed by [Freepik](http://www.freepik.com). To be exact we are using [this asset](https://www.freepik.com/free-vector/set-funny-monsters-hand-drawn-style_1933029.htm).

## License

MIT