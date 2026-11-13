package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type GameStorage struct {
	sqlClient *sqlx.DB
}

type GameStorageConfig struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c GameStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewGameStorage(cfg GameStorageConfig) (*GameStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &GameStorage{sqlClient: cfg.SQLClient}, nil
}

func (s *GameStorage) GetGame(ctx context.Context, gameID string) (*entity.Game, error) {
	var row gameRow
	query := `
		SELECT
			g.id,
			g.player_name,
			g.created_at,
			g.battle_won,
			g.scenario,
			p.id as 'partner.id',
			p.name as 'partner.name',
			p.avatar_url as 'partner.avatar_url',
			p.health as 'partner.health',
			p.max_health as 'partner.max_health',
			p.attack as 'partner.attack',
			p.defense as 'partner.defense',
			p.speed as 'partner.speed'
		FROM game g
		LEFT JOIN monster p on partner_id = p.id
		WHERE g.id = ?
	`
	err := s.sqlClient.GetContext(ctx, &row, query, gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to find game with id %s due: %w", gameID, err)
	}

	return row.toGame(), nil
}

func (s *GameStorage) SaveGame(ctx context.Context, game entity.Game) error {
	row := newGameRow(&game)
	query := `
		REPLACE INTO game (
			id, player_name, created_at, battle_won, scenario, partner_id
		) VALUES (
			:id, :player_name, :created_at, :battle_won, :scenario, :partner_id
		)
	`
	_, err := s.sqlClient.NamedExecContext(ctx, query, map[string]interface{}{
		"id":          row.ID,
		"player_name": row.PlayerName,
		"created_at":  row.CreatedAt,
		"battle_won":  row.BattleWon,
		"scenario":    row.Scenario,
		"partner_id":  row.Partner.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}

	return nil
}
