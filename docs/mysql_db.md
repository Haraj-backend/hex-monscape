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
- `extra_role`, ENUM('PARTNER'), *OPTIONAL* => the pokemon type, valid values: `PARTNER`

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
    "extra_role": "PARTNER"
},
```

**Relevant Indexes:**

- `PRIMARY_KEY` => `id`

[Back to Top](#mysql-schema)

## Table `Games`

Table that holds records of every games that has been/being played

**Relevant Fields:**

- `id`, VARCHAR(36) => identifier of a game
- `player_name`, VARCHAR(255) => name of game player
- `created_at`, BIGINT(20) => unix timestamp representation of a game creation time
- `battle_won`, INT(11) => number of battle that has been won
- `scenario`, ENUM('BATTLE_1', 'BATTLE_2', 'BATTLE_3', 'END_BATTLE') => current scenario of the game
  - valid values: `BATTLE_1`, `BATTLE_2`, `BATTLE_3`, `END_BATTLE`
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
- `FOREIGN_KEY` => `partner_id` ref to `Pokemons`.

[Back to Top](#mysql-schema)

## Table `Battles`

Table that holds records of running battle for each games

**Relevant Fields:**

- `game_id`, VARCHAR(36) => identifier of a game that the battle resides
- `state`, ENUM('DECIDE_TURN', 'ENEMY_TURN', 'PARTNER_TURN', 'WIN', 'LOSE') => current state of a battle
  - valid values: `DECIDE_TURN`, `ENEMY_TURN`, `PARTNER_TURN`, `WIN`, `LOSE`
- `partner_pokemon_id`, VARCHAR(36) => identifier of the player's partner
- `enemy_pokemon_id`, VARCHAR(36) => identifier of the enemy's pokemon
- `partner_last_damage`, Number => last inflicted damage to player's partner
- `enemy_last_damage`, Number => last inflicted damage to opposite partner

**Example Record:**

```json
{
    "game_id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
    "state": "DECIDE_TURN",
    "partner_pokemon_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "enemy_pokemon_id": "88a98dee-ce84-4afb-b5a8-7cc07535f73f",
    "partner_last_damage": 10,
    "enemy_last_damage": 25
}
```

**Relevant Indexes:**

- `PRIMARY_KEY` => `game_id`
- `FOREIGN_KEY` => `partner_pokemon_id`, `enemy_pokemon_id` ref to `Pokemons`.

[Back to Top](#mysql-schema)
