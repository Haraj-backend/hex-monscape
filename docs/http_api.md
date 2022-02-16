# HTTP API

Available endpoints for this API:

- [Get Available Pokemon Partner](#get-availaible-pokemon-partner)
- [New Game](#new-game)
- [Get Game Details](#get-game-details)
- [Get Next Scenario](#get-next-scenario)
- [Start Battle](#start-battle)
- [Get Battle Info](#get-battle-info)
- [Decide Turn](#decide-turn)
- [Attack](#attack)
- [Surrender](#surrender)

## Get Availaible Pokemon Partner

GET: `/partners`

This endpoint is used for fetching available pokemon partner. Player need to choose from one of these pokemons when starting new game.

**Example Request:**

```http
GET /partners
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "partners": [
            {
                "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
                "name": "Pikachu",
                "battle_stats": {
                    "init_health": 100,
                    "health": 100,
                    "attack": 25,
                    "defense": 5,
                    "speed": 10,
                },
                "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
            }
        ]
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---

## New Game

POST: `/games`

This endpoint is used for starting new game.

**Body Fields:**

- `player_name`, String => name of the player
- `partner_id`, String => id of pokemon partner

**Example Request:**

```http
POST /games
Content-Type: application/json

{
    "player_name": "Riandy R.N",
    "partner_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d"
}
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "id": "640dd7ef-be61-437d-a8ea-f12383185949",
        "player_name": "Riandy R.N",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Pikachu",
            "battle_stats": {
                "init_health": 100,
                "health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
        }
    },
    "ts": 1644934528
}
```

**Specific Errors:**

- Partner Not Found (`404`)

    ```http
    HTTP/1.1 404 Not Found
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_PARTNER_NOT_FOUND",
        "msg": "given `partner_id` is not found",
        "ts": 1644934528
    }
    ```

    This error will be received by client when given `partner_id` is not found.

[Back to Top](#http-api)

---

## Get Game Details

GET: `/games/{game_id}`

This endpoint is used for getting game details.

**Example Request:**

```http
GET /games/640dd7ef-be61-437d-a8ea-f12383185949
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "id": "640dd7ef-be61-437d-a8ea-f12383185949",
        "player_name": "Riandy R.N",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Pikachu",
            "battle_stats": {
                "init_health": 100,
                "health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
        },
        "created_at": 1644934528,
        "battle_won": 0,
        "scenario": "BATTLE_1"
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---

## Get Next Scenario

GET: `/games/{game_id}/scenario`

This endpoint is used for determining what next scenario to execute given current game status. It should be called everytime after a battle is done.

Possible scenarios:

- `BATTLE_1` => player need to beat the first battle
- `BATTLE_2` => player need to beat second battle
- `BATTLE_3` => player need to beat third battle
- `END_GAME` => player has won all 3 battles offered in the game, client should show ending scene

**Example Request:**

```http
GET /games/640dd7ef-be61-437d-a8ea-f12383185949/scenario
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "scenario": "BATTLE_1"
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---

## Start Battle

POST: `/games/{game_id}/battles`

This endpoint is used for initializing new battle for given game. The enemy that will be faced by player will be randomized by the system.

Everytime player finish from battle, health point for pokemon partner will be set back to full.

**Example Request:**

```http
POST /games/640dd7ef-be61-437d-a8ea-f12383185949/battles
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "id": "0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a",
        "state": "DECIDE_TURN",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Pikachu",
            "battle_stats": {
                "init_health": 100,
                "health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "init_health": 100,
                "health": 100,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png"
        }
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---

## Get Battle Info

GET: `/games/{game_id}/battles/{battle_id}`

This endpoint is used for getting battle info for specified battle id. It is useful for displaying current battle info.

In the response, there is a field called `state`. It is represents what action should be taken by the client in the battle.

Available values for the `state` are following:

- `DECIDE_TURN` => client should call [Decide Turn](#decide-turn) endpoint
- `PLAYER_TURN` => client should call either [Attack](#attack) or [Surrender](#surrender)
- `WIN` => player won the battle, client should clear the battle scene
- `LOSE` => enemy won the battle, client should clear the battle scene

**Example Request:**

```http
GET /games/640dd7ef-be61-437d-a8ea-f12383185949/battles/0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "id": "0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a",
        "state": "LOSE",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Pikachu",
            "battle_stats": {
                "init_health": 100,
                "health": 0,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "init_health": 100,
                "health": 20,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png"
        }
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---

## Decide Turn

PUT: `/games/{game_id}/battles/{battle_id}/turn`

This endpoint is used for deciding whether it is player or enemy turn to attack. If it is enemy turn, the pokemon partner will take some damage from enemy.

Turn is being randomized based on speed stats from both pokemon partner & enemy.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battles/0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a/turn
```

**Success Responses:**

- Enemy Attack:

    ```http
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "ok": true,
        "data": {
            "id": "0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a",
            "state": "DECIDE_TURN",
            "partner": {
                "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
                "name": "Pikachu",
                "battle_stats": {
                    "init_health": 100,
                    "health": 80,
                    "attack": 25,
                    "defense": 5,
                    "speed": 10,
                },
                "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
            },
            "enemy": {
                "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
                "name": "Charmander",
                "battle_stats": {
                    "init_health": 100,
                    "health": 100,
                    "attack": 30,
                    "defense": 4,
                    "speed": 10,
                },
                "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png"
            }
        },
        "ts": 1644934528
    }
    ```

- Player Turn:

    ```http
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "ok": true,
        "data": {
            "id": "0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a",
            "state": "PLAYER_TURN",
            "partner": {
                "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
                "name": "Pikachu",
                "battle_stats": {
                    "init_health": 100,
                    "health": 100,
                    "attack": 25,
                    "defense": 5,
                    "speed": 10,
                },
                "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
            },
            "enemy": {
                "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
                "name": "Charmander",
                "battle_stats": {
                    "init_health": 100,
                    "health": 100,
                    "attack": 30,
                    "defense": 4,
                    "speed": 10,
                },
                "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png"
            }
        },
        "ts": 1644934528
    }
    ```

[Back to Top](#http-api)

---

## Attack

PUT: `/games/{game_id}/battles/{battle_id}/attack`

This endpoint is used for inflicting damage to enemy. The resulted `state` of this action is `WIN`, `LOSE`, or `DECIDE_TURN`.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battles/0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a/attack
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "id": "0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a",
        "state": "DECIDE_TURN",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Pikachu",
            "battle_stats": {
                "init_health": 100,
                "health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "init_health": 100,
                "health": 84,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png"
        }
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---

## Surrender

PUT: `/games/{game_id}/battles/{battle_id}/surrender`

This endpoint is used by player to surrender current battle. The resulted `state` for this action is `LOSE`.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battles/0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a/surrender
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "id": "0f4d64d4-fd2d-4da6-bb6c-488fb4e60c2a",
        "state": "LOSE",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Pikachu",
            "battle_stats": {
                "init_health": 100,
                "health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "init_health": 100,
                "health": 84,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png"
        }
    },
    "ts": 1644934528
}
```

[Back to Top](#http-api)

---