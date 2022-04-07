package pokestrg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

const (
	partner int = 1
	enemy   int = 0
)

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	return s.getPokemonsByRole(ctx, partner)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	return s.getPokemonsByRole(ctx, enemy)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
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
		FROM pokemons
		WHERE id = ?
	`

	if err := s.db.GetContext(ctx, &pokemon, query, partnerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find partner with id %s", partnerID)
		}
		return nil, fmt.Errorf("unable to find partner with id %s: %v", partnerID, err)
	}

	return pokemon.ToPokemon(), nil
}

func (s *Storage) getPokemonsByRole(ctx context.Context, isPartnerable int) ([]entity.Pokemon, error) {
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
		FROM pokemons
		WHERE is_partnerable = ?
	`

	err := s.db.SelectContext(ctx, &pokemons, query, isPartnerable)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("unable to execute query to get list of bills due: %w", err)
	}

	return pokemons.ToPokemons(), nil
}
