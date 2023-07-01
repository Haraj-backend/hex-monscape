# DynamoDB Schema

## Table `monster`

Table that holds monster records.

**Fields:**

- `id`, String => identifier of a monster
- `name`, String => name of a monster
- `battle_stats`, Map => holds battle stats related information of a monster
  - `max_health`, Number => maximum health (on battle start) of a monster
  - `attack`, Number => number of damage that can be inflicted by a monster
  - `defense`, Number => number of damage reducer for a monster (damage = enemy.attack - your_partner.defense)
  - `speed`, Number => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
- `avatar_url`, String => url for avatar image of a monster
- `extra_role`, String, *OPTIONAL* => extra flag to define monster type, valid values: `PARTNER`

**Example Record:**

```json
{
  "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
  "name": "Bluebub",
  "battle_stats": {
    "max_health": 100,
    "attack": 25,
    "defense": 5,
    "speed": 15
  },
  "avatar_url": "https://assets.monster.com/assets/025.png",
  "extra_role": "PARTNER"
}
```

**Indexes:**

- `PRIMARY_KEY` => `id`
- `extra_role` => `extra_role`

[Back to Top](#dynamodb-schema)

## Table `game`

Table that holds records of every games that has been/being played.

**Fields:**

- `id`, String => identifier of a game
- `player_name`, String => name of game player
- `created_at`, Number => unix timestamp representation of a game creation time
- `battle_won`, Number => number of battle that has been won
- `scenario`, String => current scenario of the game, valid values: `BATTLE_1`, `BATTLE_2`, `BATTLE_3`, `END_BATTLE`
- `partner_id`, String => identifier of the player chosen partner

**Example Record:**

```json
{
  "id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
  "player_name": "Alfonso",
  "created_at": 1646205996,
  "battle_won": 2,
  "scenario": "BATTLE_3",
  "partner_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d"
}
```

**Indexes:**

- `PRIMARY_KEY` => `id`

[Back to Top](#dynamodb-schema)


## Table `battle`

Table that holds records of running battle for each games.

**Fields:**

- `game_id`, String => identifier of a game that the battle resides
- `state`, String => current state of a battle, valid values: `DECIDE_TURN`, `ENEMY_TURN`, `PARTNER_TURN`, `WIN`, `LOSE`
- `partner`, Map => holds information of player's monster
  - `id`, String => identifier of a player's partner
  - `name`, String => name of a player's partner
  - `battle_stats`, Map => holds battle stats related information of partner
    - `health`, Number => remaining health of player's partner
    - `max_health`, Number => maximum health (on battle start) of the partner
    - `attack`, Number => number of damage that can be inflicted by the partner
    - `defense`, Number => number of damage reducer for the partner (damage = enemy.attack - your_partner.defense)
    - `speed`, Number => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
  - `avatar_url`, String => url for avatar image
- `enemy`, Map => holds information of enemy's monster
  - `id`, String => identifier of a monster
  - `name`, String => name of a monster
  - `battle_stats`, Map => holds battle stats related information of a monster
  - - `health`, Number => remaining health of the enemy
    - `max_health`, Number => maximum health (on battle start) of the enemy
    - `attack`, Number => number of damage that can be inflicted by the enemy
    - `defense`, Number => number of damage reducer for the enemy (damage = your_partner.attack - enemy.defense)
    - `speed`, Number => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
  - `avatar_url`, String => url for avatar image
- `last_damage`, Map => holds information related to last damage inflicted
  - `partner`, Number => last inflicted damage to player's partner
  - `enemy`, Number => last inflicted damage to opposite partner

**Example Record:**

```json
{
  "game_id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
  "state": "DECIDE_TURN",
  "partner": {
    "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "name": "Bluebub",
    "battle_stats": {
      "max_health": 100,
      "attack": 25,
      "defense": 5,
      "speed": 15,
      "health": 75
    },
    "avatar_url": "https://assets.monster.com/assets/025.png",
    "extra_role": "PARTNER"
  },
  "enemy":{
    "id": "88a98dee-ce84-4afb-b5a8-7cc07535f73f",
    "name": "Squirtle",
    "battle_stats": {
      "max_health": 100,
      "attack": 20,
      "defense": 10,
      "speed": 15,
      "health": 60
    },
    "avatar_url": "https://assets.monster.com/assets/007.png"
  },
  "last_damage": {
    "partner": 10,
    "enemy": 25
  }
}
```

**Indexes:**

- `PRIMARY_KEY` => `game_id`


[Back to Top](#dynamodb-schema)
