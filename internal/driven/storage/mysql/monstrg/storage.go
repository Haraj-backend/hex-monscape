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
	partner int = 1
	enemy   int = 0
)

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Monster, error) {
	return s.getPokemonsByRole(ctx, partner)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Monster, error) {
	return s.getPokemonsByRole(ctx, enemy)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Monster, error) {
	var pokemon shared.PokeRow
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

	if err := s.sqlClient.GetContext(ctx, &pokemon, query, partnerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to find partner with id %s: %v", partnerID, err)
	}

	return pokemon.ToPokemon(), nil
}

func (s *Storage) getPokemonsByRole(ctx context.Context, isPartnerable int) ([]entity.Monster, error) {
	var pokemons shared.PokeRows

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
		WHERE is_partnerable = ?
	`
	if err := s.sqlClient.SelectContext(ctx, &pokemons, query, isPartnerable); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to find pokemon's role %d: %v", isPartnerable, err)
	}

	return pokemons.ToPokemons(), nil
}
