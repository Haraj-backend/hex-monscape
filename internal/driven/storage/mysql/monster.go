package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

const (
	enemy   int = 0
	partner int = 1
)

type MonsterStorage struct {
	sqlClient *sqlx.DB
}

type MonsterStorageConfig struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func (c MonsterStorageConfig) Validate() error {
	return validator.Validate(c)
}

func NewMonsterStorage(cfg MonsterStorageConfig) (*MonsterStorage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &MonsterStorage{sqlClient: cfg.SQLClient}, nil
}

func (s *MonsterStorage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	return s.getMonstersByRole(ctx, partner)
}

func (s *MonsterStorage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	return s.getMonstersByRole(ctx, enemy)
}

func (s *MonsterStorage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	var row monsterRow
	query := `
		SELECT
			id,
			name,
			health,
			max_health,
			attack,
			defense,
			speed,
			avatar_url
		FROM monster
		WHERE id = ?
	`
	err := s.sqlClient.GetContext(ctx, &row, query, partnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to find partner with id %s due: %w", partnerID, err)
	}

	return row.toMonster(), nil
}

func (s *MonsterStorage) getMonstersByRole(ctx context.Context, role int) ([]entity.Monster, error) {
	// build query, every monster can be an enemy, but only partnerable monsters can be a partner
	args := []interface{}{}
	query := `
		SELECT
			id,
			name,
			health,
			max_health,
			attack,
			defense,
			speed,
			avatar_url
		FROM monster
	`
	if role == partner {
		query += "WHERE is_partnerable = ?"
		args = append(args, role)
	}

	// execute query
	var rows monsterRows
	err := s.sqlClient.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return rows.toMonsters(), nil
}
