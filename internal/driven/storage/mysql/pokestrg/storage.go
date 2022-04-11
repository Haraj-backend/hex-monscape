package pokestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	sqlClient *sqlx.DB
}

const (
	partner int = 1
	enemy   int = 0
)

func New(cfg shared.Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
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

	if err := s.sqlClient.GetContext(ctx, &pokemon, query, partnerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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

	if err := s.sqlClient.GetContext(ctx, &pokemons, query, isPartnerable); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to find pokemon's role %d: %v", isPartnerable, err)
	}

	return pokemons.ToPokemons(), nil
}
