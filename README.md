# Hex Monscape Go

Welcome to this repo!

In this repo you will learn the concept of [Hexagonal Architecture](./docs/hex_arch.md) and its [Go](https://go.dev/) implementation through simple server-client game called `Hex Monscape`.

We are using [Hexagonal Architecture](./docs/hex_arch.md) to implement the game server in [Go](https://go.dev/) while the client implemented using [Vue 3](https://vuejs.org/). To know the API details for this game, please see [this document](./docs/http_api.md).

In the game you will play as a 10 years old monster hunter that dreams to become the very best. In order to reach that, you need to make journey together with your monster partner to seek 3 strong wild monsters and defeat them. 🥷🏻🥷🏻🥷🏻

As Solutions Team member, your understanding towards [Hexagonal Architecture](./docs/hex_arch.md) is mandatory since it is the main architecture we used for building Haraj production systems. So if you understand this architecture well, you will be in no time contributing to Haraj production.

Please refer to [Primary References](#primary-references) section to start learning about the concepts presented in this repo.

## How to Run The Game

When we are using [Hexagonal Architecture](./docs/hex_arch.md) to design a system, it is quite easy to swap its infrastructure code with another technologies.

So for example, if initially we used in memory storage to store our data, we could easily swap it with MySQL storage or something else.

To prove this point, there are 3 variants of game server in this project:

- Server using Memory storage
- Server using DynamoDB storage
- Server using MySQL storage

All of them will serve the same game, the only difference is the place where they store the game data.

All of these servers could be run by using this command:

```bash
> make run
```

This command will create & run the stack defined in this [docker-compose.yml](./docker-compose.yml) file. 

Wait a moment until the entire stack running. You will something like this in the console after it is successfully running:

```bash
hex_mem_1     | 2022/05/11 16:29:50 server is listening on :9185...
hex_mysql_1   | 2022/05/11 16:30:21 server is listening on :9186...
hex_ddb_1     | 2022/05/11 16:30:21 server is listening on :9187...
```

After that you could access each of these servers by visiting endpoint below:

- Memory storage => http://localhost:9185
- DynamoDB storage => http://localhost:9186
- MySQL storage => http://localhost:9187

## Primary References

To start learning the concept of [Hexagonal Architecture](./docs/hex_arch.md) please use [this document](./docs/hex_arch.md) as your primary source for learning. This is so you won't be having too much confusion when learning it from other online resources.

To know more about the game design, please refer to [this document](./docs/game_design.md).

To learn about the methodology on how to implement [Hexagonal Architecture](./docs/hex_arch.md) on a project, please refer to [this document](./docs/project_method.md).

## Attribution

The monster characters used in this project is designed by [Freepik](http://www.freepik.com). To be exact we are using [this asset](https://www.freepik.com/free-vector/set-funny-monsters-hand-drawn-style_1933029.htm).

## License

MIT