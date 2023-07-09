package monstrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/shared"
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

const (
	enemy   int = 0
	partner int = 1
)

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	return s.getMonsterByRole(ctx, partner)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	return s.getMonsterByRole(ctx, enemy)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	var row shared.MonsterRow
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

	if err := s.sqlClient.GetContext(ctx, &row, query, partnerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to find partner with id %s: %v", partnerID, err)
	}

	return row.ToMonster(), nil
}

func (s *Storage) getMonsterByRole(ctx context.Context, role int) ([]entity.Monster, error) {
	var rows shared.MonsterRows
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

	if err := s.sqlClient.SelectContext(ctx, &rows, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}

	return rows.ToMonsters(), nil
}
