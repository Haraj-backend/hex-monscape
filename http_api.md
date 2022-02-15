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
                    "health": 100,
                    "attack": 25,
                    "defense": 5,
                    "speed": 10,
                },
                "avatar": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
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
                "health": 100,
                "attack": 25,
                "defense": 5,
                "speed": 10,
            },
            "avatar": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
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

[Back to Top](#http-api)

---

## Get Next Scenario

GET: `/games/{game_id}/scenario`

[Back to Top](#http-api)

---

## Start Battle

POST: `/games/{game_id}/battles`

[Back to Top](#http-api)

---

## Get Battle Info

GET: `/games/{game_id}/battles/{battle_id}`

[Back to Top](#http-api)

---

## Decide Turn

PUT: `/games/{game_id}/battles/{battle_id}/turn`

[Back to Top](#http-api)

---

## Attack

PUT: `/games/{game_id}/battles/{battle_id}/attack`

[Back to Top](#http-api)

---

## Surrender

PUT: `/games/{game_id}/battles/{battle_id}/surrender`

[Back to Top](#http-api)

---