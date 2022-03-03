# DynamoDB Schema

## Table `Pokemons`

Table that holds records of all available pokemons.

**Relevant Fields:**

- `id`, String => identifier of a pokemon
- `name`, String => name of a pokemon
- `battle_stats`, Map => holds battle stats related information of a pokemon
  - `max_health`, Number => maximum health (on battle start) of a pokemon
  - `attack`, Number => number of damage that can be inflicted by a pokemon
  - `defense`, Number => number of damage reducer for a pokemon (damage = enemy.attack - your_partner.defense)
  - `speed`, Number => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
- `avatar_url`, String => url for avatar image of a pokemon
- `extra_role`, String, *OPTIONAL* => the pokemon type, valid values: `PARTNER`

**Example Record:**

```
{
    "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "name": "Pikachu",
    "battle_stats": {
        "max_health": 100,
        "attack": 25,
        "defense": 5,
        "speed": 15
    },
    "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png",
    "extra_role": "PARTNER"
},
```

**Relevant Indexes:**
- `PRIMARY_KEY` => `id`
- `extra_role`, GSI => `extra_role`

[Back to Top](#dynamodb-schema)

## Table `Games`

Table that holds records of every games that has been/being played

**Relevant Fields:**

- `id`, String => identifier of a game
- `player_name`, String => name of game player
- `created_at`, Number => unix timestamp representation of a game creation time
- `battle_won`, Number => number of battle that has been won
- `scenario`, String => current scenario of the game, valid values: `BATTLE_1`, `BATTLE_2`, `BATTLE_3`, `END_BATTLE`
- `partner_id`, String => identifier of the player chosen partner

**Example Record:**

```
{
    "id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
    "player_name": "Alfonso",
    "created_at": 1646205996,
    "battle_won": 2,
    "scenario": "BATTLE_3",
    "partner_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d"
}
```

**Relevant Indexes:**
- `PRIMARY_KEY` => `id`

[Back to Top](#dynamodb-schema)


## Table `Battles`

Table that holds records of running battle for each games

**Relevant Fields:**
- `game_id`, String => identifier of a game that the battle resides
- `state`, String => current state of a battle, valid values: `DECIDE_TURN`, `ENEMY_TURN`, `PARTNER_TURN`, `WIN`, `LOSE`
- `partner`, Map => holds information of player's pokemon
  - `health`, Number => remaining health of player's partner
- `enemy`, Map => holds information of enemy's pokemon
  - `id`, String => identifier of opposite partner
  - `health`, Number => remaining health of opposite partner
- `last_damage`, Map => holds information related to last damage inflicted
  - `partner`, Number => last inflicted damage to player's partner
  - `enemy`, Number => last inflicted damage to opposite partner

**Example Record:**

```
{
    "game_id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
    "state": "DECIDE_TURN",
    "partner": {
        "health": 90
    },
    "enemy": {
        "id": "1eb64af3-713a-4210-8c5a-883312c51fa3",
        "health": 75
    },
    "last_damage": {
        "partner": 10,
        "enemy": 25
    }
}
```

**Relevant Indexes:**
- `PRIMARY_KEY` => `game_id`


[Back to Top](#dynamodb-schema)
