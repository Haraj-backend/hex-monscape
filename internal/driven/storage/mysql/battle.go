package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type BattleStorage struct {
	sqlClient *sqlx.DB
}

type BattleStorageConfig struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c BattleStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewBattleStorage(cfg BattleStorageConfig) (*BattleStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &BattleStorage{sqlClient: cfg.SQLClient}, nil
}

func (s *BattleStorage) GetBattle(ctx context.Context, gameID string) (*entity.Battle, error) {
	var row battleRow
	query := `SELECT * FROM battle WHERE game_id = ?`
	err := s.sqlClient.GetContext(ctx, &row, query, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return row.toBattle(), nil
}

func (s *BattleStorage) SaveBattle(ctx context.Context, b entity.Battle) error {
	row := newBattleRow(b)
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
	_, err := s.sqlClient.NamedExecContext(ctx, query, row)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}

	return nil
}
