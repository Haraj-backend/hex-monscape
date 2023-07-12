# REST API

Available endpoints for this API:

- [Get Available Partners](#get-available-partners)
- [New Game](#new-game)
- [Get Game Details](#get-game-details)
- [Get Next Scenario](#get-next-scenario)
- [Start Battle](#start-battle)
- [Get Battle Info](#get-battle-info)
- [Decide Turn](#decide-turn)
- [Attack](#attack)
- [Surrender](#surrender)

## Get Available Partners

GET: `/partners`

This endpoint is used for fetching available monster partner. Player need to choose from one of these monsters when starting new game.

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
                "name": "Bluebub",
                "battle_stats": {
                    "health": 100,
                    "max_health": 100,
                    "attack": 25,
                    "defense": 5,
                    "speed": 10,
                },
                "avatar_url": "https://assets.monster.com/assets/025.png"
            }
        ]
    },
    "ts": 1644934528
}
```

[Back to Top](#rest-api)

---

## New Game

POST: `/games`

This endpoint is used for starting new game.

**Body Fields:**

- `player_name`, String => name of the player
- `partner_id`, String => id of monster partner

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
            "name": "Bluebub",
            "battle_stats": {
                "health": 100,
                "max_health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/025.png"
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

[Back to Top](#rest-api)

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
            "name": "Bluebub",
            "battle_stats": {
                "health": 100,
                "max_health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/025.png"
        },
        "created_at": 1644934528,
        "battle_won": 0,
        "scenario": "BATTLE_1"
    },
    "ts": 1644934528
}
```

**Specific Errors:**

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

[Back to Top](#rest-api)

---

## Get Scenario

GET: `/games/{game_id}/scenario`

This endpoint is used for determining what scenario should be executed by client. It should be called after every battle is done.

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

**Specific Errors:**

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

[Back to Top](#rest-api)

---

## Start Battle

PUT: `/games/{game_id}/battle`

This endpoint is used for initializing new battle for given game. The enemy that will be faced by player will be randomized by the system.

Everytime player finish from battle, health point for monster partner will be set back to full.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battle
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "game_id": "640dd7ef-be61-437d-a8ea-f12383185949",
        "state": "DECIDE_TURN",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Bluebub",
            "battle_stats": {
                "health": 100,
                "max_health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "health": 100,
                "max_health": 100,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/004.png"
        },
        "last_damage": {
            "partner": 0,
            "enemy": 0
        }
    },
    "ts": 1644934528
}
```

**Specific Errors:**

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

- Invalid Battle State (`409`)

  ```http
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_INVALID_BATTLE_STATE",
      "msg": "invalid battle state",
      "ts": 1644934528
  }
  ```

  Client receive this error when battle state when client executing action is invalid.

[Back to Top](#rest-api)

---

## Get Battle Info

GET: `/games/{game_id}/battle`

This endpoint is used for getting battle info for specified battle id. It is useful for display current battle info.

In the response, there is a field called `state`. It is represents what action should be taken by the client in the battle.

Available values for the `state` are following:

- `DECIDE_TURN` => client should call [Decide Turn](#decide-turn) endpoint
- `PARTNER_TURN` => client should call either [Attack](#attack) or [Surrender](#surrender)
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
        "game_id": "640dd7ef-be61-437d-a8ea-f12383185949",
        "state": "LOSE",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Bluebub",
            "battle_stats": {
                "health": 100,
                "max_health": 0,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "health": 100,
                "max_health": 20,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/004.png"
        },
        "last_damage": {
            "partner": 100,
            "enemy": 0
        }
    },
    "ts": 1644934528
}
```

**Specific Errors:**

- Battle Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_BATTLE_NOT_FOUND",
      "msg": "battle is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when ongoing battle is not found for given game id.

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

[Back to Top](#rest-api)

---

## Decide Turn

PUT: `/games/{game_id}/battle/turn`

This endpoint is used for deciding whether it is player or enemy turn to attack. If it is enemy turn, the monster partner will take some damage from enemy.

Turn is being randomized based on speed stats from both monster partner & enemy.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battle/turn
```

**Success Responses:**

- Enemy Attack:

  ```http
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
      "ok": true,
      "data": {
          "game_id": "640dd7ef-be61-437d-a8ea-f12383185949",
          "state": "DECIDE_TURN",
          "partner": {
              "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
              "name": "Bluebub",
              "battle_stats": {
                  "health": 80,
                  "max_health": 100,
                  "attack": 25,
                  "defense": 5,
                  "speed": 10,
              },
              "avatar_url": "https://assets.monster.com/assets/025.png"
          },
          "enemy": {
              "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
              "name": "Charmander",
              "battle_stats": {
                  "health": 100,
                  "max_health": 100,
                  "attack": 30,
                  "defense": 4,
                  "speed": 10,
              },
              "avatar_url": "https://assets.monster.com/assets/004.png"
          },
          "last_damage": {
              "partner": 20,
              "enemy": 0
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
          "game_id": "640dd7ef-be61-437d-a8ea-f12383185949",
          "state": "PARTNER_TURN",
          "partner": {
              "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
              "name": "Bluebub",
              "battle_stats": {
                  "health": 100,
                  "max_health": 100,
                  "attack": 25,
                  "defense": 5,
                  "speed": 10,
              },
              "avatar_url": "https://assets.monster.com/assets/025.png"
          },
          "enemy": {
              "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
              "name": "Charmander",
              "battle_stats": {
                  "health": 100,
                  "max_health": 100,
                  "attack": 30,
                  "defense": 4,
                  "speed": 10,
              },
              "avatar_url": "https://assets.monster.com/assets/004.png"
          },
          "last_damage": {
              "partner": 0,
              "enemy": 0
          }
      },
      "ts": 1644934528
  }
  ```

**Specific Errors:**

- Battle Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_BATTLE_NOT_FOUND",
      "msg": "battle is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when ongoing battle is not found for given game id.

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

- Invalid Battle State (`409`)

  ```http
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_INVALID_BATTLE_STATE",
      "msg": "invalid battle state",
      "ts": 1644934528
  }
  ```

  Client receive this error when battle state when client executing action is invalid.

[Back to Top](#rest-api)

---

## Attack

PUT: `/games/{game_id}/battle/attack`

This endpoint is used for inflicting damage to enemy. The resulted `state` of this action is `WIN`, `LOSE`, or `DECIDE_TURN`.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battle/attack
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "game_id": "640dd7ef-be61-437d-a8ea-f12383185949",
        "state": "DECIDE_TURN",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Bluebub",
            "battle_stats": {
                "health": 100,
                "max_health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "health": 84,
                "max_health": 100,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/004.png"
        },
        "last_damage": {
            "partner": 0,
            "enemy": 16
        }
    },
    "ts": 1644934528
}
```

**Specific Errors:**

- Battle Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_BATTLE_NOT_FOUND",
      "msg": "battle is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when ongoing battle is not found for given game id.

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

- Invalid Battle State (`409`)

  ```http
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_INVALID_BATTLE_STATE",
      "msg": "invalid battle state",
      "ts": 1644934528
  }
  ```

  Client receive this error when battle state when client executing action is invalid.

[Back to Top](#rest-api)

---

## Surrender

PUT: `/games/{game_id}/battle/surrender`

This endpoint is used by player to surrender current battle. The resulted `state` for this action is `LOSE`.

**Example Request:**

```http
PUT /games/640dd7ef-be61-437d-a8ea-f12383185949/battle/surrender
```

**Success Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "game_id": "640dd7ef-be61-437d-a8ea-f12383185949",
        "state": "LOSE",
        "partner": {
            "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
            "name": "Bluebub",
            "battle_stats": {
                "health": 100,
                "max_health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/025.png"
        },
        "enemy": {
            "id": "28933dde-b04c-46cc-9be7-5e785c62adfa",
            "name": "Charmander",
            "battle_stats": {
                "health": 84,
                "max_health": 100,
                "attack": 30,
                "defense": 4,
                "speed": 10,
            },
            "avatar_url": "https://assets.monster.com/assets/004.png"
        },
        "last_damage": {
            "partner": 0,
            "enemy": 0
        }
    },
    "ts": 1644934528
}
```

**Specific Errors:**

- Battle Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_BATTLE_NOT_FOUND",
      "msg": "battle is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when ongoing battle is not found for given game id.

- Game Not Found (`404`)

  ```http
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_GAME_NOT_FOUND",
      "msg": "game is not found",
      "ts": 1644934528
  }
  ```

  Client receive this error when game is not found.

- Invalid Battle State (`409`)

  ```http
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
      "ok": false,
      "err": "ERR_INVALID_BATTLE_STATE",
      "msg": "invalid battle state",
      "ts": 1644934528
  }
  ```

  Client receive this error when battle state when client executing action is invalid.

[Back to Top](#rest-api)

---
