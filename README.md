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

This app is powered by docker. So make sure to install it before running below command:

```bash
> make run
```

Upon success, your console should output message like following:

```bash
2022/02/20 15:53:42 server is listening on :9186...
```

## How to Run in serverless mode

We will run the app in serverless by using SAM

```bash
> make deploy-local
```

Upon success, docker compose should have log like following:

```
Mounting EntrypointFunction at http://127.0.0.1:3000/ [DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT]
Mounting EntrypointFunction at http://127.0.0.1:3000/{proxy+} [DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT]
You can now browse to the above endpoints to invoke your functions. You do not need to restart/reload SAM CLI while working on your functions, changes will be reflected instantly/automatically. You only need to restart SAM CLI if you update your AWS SAM template
2022-03-21 17:58:42  * Running on http://127.0.0.1:3000/ (Press CTRL+C to quit)
```

After the message is shown, you could access http://localhost:3000 using your browser to play the game. You may need to refresh the page several times before the page could displayed properly. This is a problem related to API Gateway return 502 for some assets in several first requests.
