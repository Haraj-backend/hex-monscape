# Mysql Schema

## Table `Pokemons`

Table that holds records of all available pokemons.

**Relevant Fields:**

- `id`, VARCHAR(36) => identifier of a pokemon
- `name`, VARCHAR(255) => name of a pokemon
- `max_health`, SMALLINT => maximum health (on battle start) of a pokemon
- `attack`, SMALLINT => number of damage that can be inflicted by a pokemon
- `defense`, SMALLINT => number of damage reducer for a pokemon (damage = enemy.attack - your_partner.defense)
- `speed`, SMALLINT => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
- `avatar_url`, TEXT => url for avatar image of a pokemon
- `is_partnerable`, TINYINT(1) => whether pokemon is partner or not

**Example Record:**

```json
{
    "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "name": "Pikachu",
    "max_health": 100,
    "attack": 25,
    "defense": 5,
    "speed": 15,
    "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png",
    "is_partnerable": 1
},
```

**Relevant Indexes:**

- `PRIMARY_KEY` => `id`
- `is_partnerable`

[Back to Top](#mysql-schema)

## Table `Games`

Table that holds records of every games that has been/being played

**Relevant Fields:**

- `id`, VARCHAR(36) => identifier of a game
- `player_name`, VARCHAR(255) => name of game player
- `created_at`, BIGINT(20) => unix timestamp representation of a game creation time
- `battle_won`, INT(11) => number of battle that has been won
- `scenario`, VARCHAR(30) => current scenario of the game
  - current valid values: `BATTLE_1`, `BATTLE_2`, `BATTLE_3`, `END_BATTLE`
- `partner_id`, VARCHAR(36) => identifier of the player chosen partner (pokemon)

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

**Relevant Indexes:**

- `PRIMARY_KEY` => `id`
- `FOREIGN_KEY` => `partner_id` ref to `Pokemons.id`.

[Back to Top](#mysql-schema)

## Table `Battles`

Table that holds records of running battle for each games

**Relevant Fields:**

- `game_id`, VARCHAR(36) => identifier of a game that the battle resides
- `state`, VARCHAR(30) => current state of a battle
  - current valid values: `DECIDE_TURN`, `ENEMY_TURN`, `PARTNER_TURN`, `WIN`, `LOSE`
- `partner_state_id`, VARCHAR(36) => identifier of the player's partner state
- `enemy_state_id`, VARCHAR(36) => identifier of the enemy's pokemon state
- `partner_last_damage`, Number => last inflicted damage to player's partner
- `enemy_last_damage`, Number => last inflicted damage to opposite partner

**Example Record:**

```json
{
    "game_id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
    "state": "DECIDE_TURN",
    "partner_state_id": 1,
    "enemy_state_id": 2,
    "partner_last_damage": 10,
    "enemy_last_damage": 25
}
```

**Relevant Indexes:**

- `PRIMARY_KEY` => `game_id`
- `FOREIGN_KEY` => `partner_state_id`, `enemy_state_id` ref to `Pokemon_Battle_States.id`.

[Back to Top](#mysql-schema)

## Table `Pokemon_Battle_States`

Table that holds records of all available pokemons states.

**Relevant Fields:**

- `id`, INT, AUTO INCREMENT => identifier of a pokemon's battle state
- `pokemon_id`, VARCHAR(36) => identifier of a pokemon
- `name`, VARCHAR(255) => name of a pokemon
- `max_health`, SMALLINT => maximum health (on battle start) of a pokemon
- `health`, SMALLINT => current health of a pokemon
- `attack`, SMALLINT => number of damage that can be inflicted by a pokemon
- `defense`, SMALLINT => number of damage reducer for a pokemon (damage = enemy.attack - your_partner.defense)
- `speed`, SMALLINT => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
- `avatar_url`, TEXT => url for avatar image of a pokemon

**Example Record:**

```json
{
    "id": 1,
    "pokemon_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "name": "Pikachu",
    "max_health": 100,
    "health": 75,
    "attack": 25,
    "defense": 5,
    "speed": 15,
    "avatar_url": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
},
```

**Relevant Indexes:**

- `PRIMARY_KEY` => `id`
- `FOREIGN_KEY` => `pokemon_id` ref to `Pokemons.id`.

[Back to Top](#mysql-schema)
