# Mysql Schema

## Table `monster`

Table that holds records of all available monsters.

**Fields:**

- `id`, VARCHAR(36) => identifier of a monster
- `name`, VARCHAR(255) => name of a monster
- `health`, INT(11) => health of a monster
- `max_health`, INT(11) => maximum health (on battle start) of a monster
- `attack`, INT(11) => number of damage that can be inflicted by a monster
- `defense`, INT(11) => number of damage reducer for a monster (damage = enemy.attack - your_partner.defense)
- `speed`, INT(11) => chance for getting a turn in battle, higher means more likely to get a turn in battle RNG
- `avatar_url`, TEXT => url for avatar image of a monster
- `is_partnerable`, TINYINT(1) => whether monster is partnerable or not

**Example Record:**

```json
{
    "id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "name": "Bluebub",
    "max_health": 100,
    "attack": 25,
    "defense": 5,
    "speed": 15,
    "avatar_url": "https://assets.monster.com/assets/025.png",
    "is_partnerable": 1
},
```

**Indexes:**

- `PRIMARY_KEY` => `id`
- `is_partnerable` => `is_partnerable`

[Back to Top](#mysql-schema)

## Table `game`

Table that holds records of every games that has been/being played

**Fields:**

- `id`, VARCHAR(36) => identifier of a game
- `player_name`, VARCHAR(255) => name of game player
- `created_at`, BIGINT(20) => unix timestamp representation of a game creation time
- `battle_won`, INT(11) => number of battle that has been won by player
- `scenario`, VARCHAR(30) => current scenario of the game, valid values: `BATTLE_1`, `BATTLE_2`, `BATTLE_3`, `END_BATTLE`
- `partner_id`, VARCHAR(36) => id of monster partner chosen by player

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

[Back to Top](#mysql-schema)

## Table `battle`

This table is used for holding complete information for single battle.

**Fields:**

- `game_id`, VARCHAR(36) => id of the game where this battle is occurring
- `state`, VARCHAR(30) => current battle state, valid values: `DECIDE_TURN`, `ENEMY_TURN`, `PARTNER_TURN`, `WIN`, `LOSE`
- `partner_monster_id`, VARCHAR(36) => id of the monster partner, has reference to `monsters` table
- `partner_name`, VARCHAR(255) => name of the monster partner
- `partner_max_health`, INT(11) => max health of monster partner
- `partner_health`, INT(11) => current health of monster partner
- `partner_attack`, INT(11) => attack of monster partner
- `partner_defense`, INT(11) => defense of monster partner
- `partner_speed`, INT(11) => speed of monster partner
- `partner_avatar_url`, TEXT => avatar url of monster partner
- `partner_last_damage`, Number => last inflicted damage to player's partner
- `enemy_monster_id`, VARCHAR(36) => id of the monster enemy, has reference to `monsters` table
- `enemy_name`, VARCHAR(255) => name of the monster enemy
- `enemy_max_health`, INT(11) => max health of monster enemy
- `enemy_health`, INT(11) => current health of monster enemy
- `enemy_attack`, INT(11) => attack of monster enemy
- `enemy_defense`, INT(11) => defense of monster enemy
- `enemy_speed`, INT(11) => speed of monster enemy
- `enemy_avatar_url`, TEXT => avatar url of monster enemy
- `enemy_last_damage`, Number => last inflicted damage to enemy

**Example Record:**

```json
{
    "game_id": "1a34a63d-afe6-4186-8628-13a25eaa6076",
    "state": "DECIDE_TURN",
    "partner_monster_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "partner_name": "Bluebub",
    "partner_max_health": 100,
    "partner_attack": 25,
    "partner_defense": 5,
    "partner_speed": 15,
    "partner_avatar_url": "https://assets.monster.com/assets/025.png",
    "partner_last_damage": 10,
    "enemy_monster_id": "b1c87c5c-2ac3-471d-9880-4812552ee15d",
    "enemy_name": "Bluebub",
    "enemy_max_health": 100,
    "enemy_attack": 25,
    "enemy_defense": 5,
    "enemy_speed": 15,
    "enemy_avatar_url": "https://assets.monster.com/assets/025.png",
    "enemy_last_damage": 25
}
```

**Indexes:**

- `PRIMARY_KEY` => `game_id`

[Back to Top](#mysql-schema)