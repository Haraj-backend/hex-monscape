# Hex PokeBattle

Hex PokeBattle is a simple web game implemented using [Hexagonal Architecture](./docs/hex_arch.md).

It is intended to become Solutions Team reference when they want to implement web app using Hexagonal Architecture.

To learn more about API for this game check out [HTTP API](./docs/http_api.md) doc.

To learn the methodology of how to create web app using Hexagonal Architecture, check out [Project Methodology](./docs/project_method.md) doc.

## Game Concept

The game concept is pretty simple, player just need to choose his/her pokemon partner & won battle for 3 times to beat the game.

Here is the flowchart of the game:

<p align="center">
    <img src="./docs/game_flow.svg" alt="Game Flow" height="400" />
</p>

Here is the flowchart for each battle in the game:

<p align="center">
    <img src="./docs/battle_flow.svg" alt="Battle Flow" height="400" />
</p>

## How to Run

There are 3 variants of server in this project:

- Server using Memory storage
- Server using DynamoDB storage
- Server using MySQL storage

These variants could be run by using this command:

```bash
> make run
```

This command will create & run the stack defined in this [docker-compose.yml](./docker-compose.yml) file. 

Wait a moment until the entire stack setup done. You will something like this in the console log after the setup is done:

```bash
hex_mem_1     | 2022/05/11 16:29:50 server is listening on :9186...
hex_mysql_1   | 2022/05/11 16:30:21 Running in server mode at :9186
hex_ddb_1     | 2022/05/11 16:30:21 Running in server mode at :9186
```

After that we could access each of variants by visiting the urls below:

- Memory storage => http://localhost:9185
- DynamoDB storage => http://localhost:9186
- MySQL storage => http://localhost:9187