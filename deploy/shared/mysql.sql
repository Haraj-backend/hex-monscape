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
