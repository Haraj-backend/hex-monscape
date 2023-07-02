package battlestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/battle"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type Storage struct {
	sqlClient *sqlx.DB
}

type Config struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetBattle(ctx context.Context, gameID string) (*battle.Battle, error) {
	query := `SELECT * FROM battle WHERE game_id = ?`
	var row battleRow
	if err := s.sqlClient.GetContext(ctx, &row, query, gameID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}
	return row.ToBattle(), nil
}

func (s *Storage) SaveBattle(ctx context.Context, b battle.Battle) error {
	battleRow := newBattleRow(b)
	query := `
		REPLACE INTO battle (
			game_id, state, partner_monster_id, 
			partner_name, partner_max_health, partner_health, 
			partner_attack, partner_defense, partner_speed, 
			partner_avatar_url, partner_last_damage, enemy_monster_id, 
			enemy_name, enemy_max_health, enemy_health,
			enemy_attack, enemy_defense, enemy_speed, 
			enemy_avatar_url, enemy_last_damage
		) VALUES (
			:game_id, :state, :partner_monster_id, 
			:partner_name, :partner_max_health, :partner_health, 
			:partner_attack, :partner_defense, :partner_speed, 
			:partner_avatar_url, :partner_last_damage, :enemy_monster_id, 
			:enemy_name, :enemy_max_health, :enemy_health,
			:enemy_attack, :enemy_defense, :enemy_speed, 
			:enemy_avatar_url, :enemy_last_damage
		)
	`
	_, err := s.sqlClient.NamedExecContext(ctx, query, battleRow)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}
