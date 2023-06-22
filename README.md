# Hex Monscape Go

Welcome to this repo!

In this repo you will learn the concept of [Hexagonal Architecture](./docs/hex_arch.md) and its Go implementation through simple server-client game called `Hex Monscape`.

In the game you will play as a 10 years old monster hunter that dreams to become the very best. In order to reach that, you make journey together with your monster partner to seek 3 strong wild monsters and defeat them. 🥷🏻🥷🏻🥷🏻

As Solutions Team member, your understanding towards [Hexagonal Architecture](./docs/hex_arch.md) is mandatory since it is the main architecture we used for building production-grade systems for Haraj. So if you understand this architecture very well, you will be in no time to contribute to Haraj production systems.

Table of contents:

- [Why Hexagonal Architecture?](#why-hexagonal-architecture)
- [Game Design](#game-design)
- [How to Run The Game](#how-to-run-the-game)
- [Project Methodology](#project-methodology)

## Why Hexagonal Architecture?

Some of you might be wondering why we need to learn about [Hexagonal Architecture](./docs/hex_arch.md) rather than 

## Game Design

The game concept is pretty simple, player just need to choose his/her pokemon partner & won battle for 3 times to beat the game.

Here is the flowchart of the game:

<p align="center">
    <img src="./docs/game_flow.svg" alt="Game Flow" height="400" />
</p>

Here is the flowchart for each battle in the game:

<p align="center">
    <img src="./docs/battle_flow.svg" alt="Battle Flow" height="400" />
</p>

## How to Run The Game

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

## Project Methodology