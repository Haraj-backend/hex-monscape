CREATE TABLE IF NOT EXISTS pokemons (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  health INT(11) NOT NULL,
  max_health INT(11) NOT NULL,
  attack INT(11) NOT NULL,
  defense INT(11) NOT NULL,
  speed INT(11) NOT NULL,
  avatar_url TEXT NOT NULL,
  is_partnerable TINYINT(1) NOT NULL,
  KEY `is_partnerable` (`is_partnerable`)
);

CREATE TABLE IF NOT EXISTS games (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  player_name VARCHAR(255) NOT NULL,
  created_at BIGINT(20) NOT NULL,
  battle_won INT(11) NOT NULL,
  scenario VARCHAR(30) NOT NULL,
  partner_id VARCHAR(36) NOT NULL,
  FOREIGN KEY (partner_id) REFERENCES pokemons(id)
);

CREATE TABLE IF NOT EXISTS battles (
  game_id VARCHAR(36) NOT NULL PRIMARY KEY,
  state VARCHAR(30) NOT NULL,
  partner_pokemon_id VARCHAR(36) NOT NULL,
  partner_name VARCHAR(255) NOT NULL,
  partner_max_health INT(11) NOT NULL,
  partner_health INT(11) NOT NULL,
  partner_attack INT(11) NOT NULL,
  partner_defense INT(11) NOT NULL,
  partner_speed INT(11) NOT NULL,
  partner_avatar_url TEXT NOT NULL,
  partner_last_damage INT(11) NOT NULL,
  enemy_pokemon_id VARCHAR(36) NOT NULL,
  enemy_name VARCHAR(255) NOT NULL,
  enemy_max_health INT(11) NOT NULL,
  enemy_health INT(11) NOT NULL,
  enemy_attack INT(11) NOT NULL,
  enemy_defense INT(11) NOT NULL,
  enemy_speed INT(11) NOT NULL,
  enemy_avatar_url TEXT NOT NULL,
  enemy_last_damage INT(11) NOT NULL,
  FOREIGN KEY (partner_pokemon_id) REFERENCES pokemons(id),
  FOREIGN KEY (enemy_pokemon_id) REFERENCES pokemons(id)
);

INSERT INTO pokemons (id, name, max_health, health, attack, defense, speed, avatar_url, is_partnerable) VALUES
  ("b1c87c5c-2ac3-471d-9880-4812552ee15d", 'Pikachu', 100, 100, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png", 1),
  ("0f9b84b6-a768-4ba9-8800-207740fc993d", 'Bulbasaur', 100, 100, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/001.png", 1),
  ("52939c53-2ef1-4bc7-b896-894c80cb2eca", 'Charmander', 100, 100, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/004.png", 1),
  ("d0090631-89c8-4529-9ff6-8daaa6a0f476", 'Squirtle', 50, 50, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/007.png", 0),
  ("be11388f-9f10-4acf-b09c-b4769771c32d", 'Pidgey', 50, 50, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/016.png", 0),
  ("2c635ff9-0b6b-4c10-9a57-3c6442cb777a", 'Rattata', 50, 50, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/019.png", 0),
  ("5a0b1af9-c6ce-4893-87b1-39d0dc07fee1", 'Ekans', 50, 50, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/023.png", 0),
  ("55b768a1-5928-4fda-8c9f-135b54eec3f7", 'Caterpie', 50, 50, 49, 49, 45, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/010.png", 0);