package pokestrg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	db "github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

const (
	partner int = 1
	enemy   int = 0
)

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetAvailablePartners(ctx context.Context) ([]entity.Pokemon, error) {
	return s.getPokemonsByRole(ctx, partner)
}

func (s *Storage) GetPossibleEnemies(ctx context.Context) ([]entity.Pokemon, error) {
	return s.getPokemonsByRole(ctx, enemy)
}

func (s *Storage) GetPartner(ctx context.Context, partnerID string) (*entity.Pokemon, error) {
	var pokemon entity.Pokemon

	query := "SELECT id, name, health, max_health, attack, defense, speed, avatar_url FROM pokemons WHERE id = ?"

	if err := mappingPokemon(s.db.QueryRowContext(ctx, query, partnerID), &pokemon); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find partner with id %s", partnerID)
		}
		return nil, fmt.Errorf("unable to find partner with id %s: %v", partnerID, err)
	}

	return &pokemon, nil
}

func (s *Storage) getPokemonsByRole(ctx context.Context, isPartnerable int) ([]entity.Pokemon, error) {
	var pokemons []entity.Pokemon

	query := "SELECT id, name, health, max_health, attack, defense, speed, avatar_url FROM pokemons WHERE is_partnerable = ?"

	rows, err := s.db.QueryContext(ctx, query, isPartnerable)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pk entity.Pokemon
		if err := mappingPokemon(rows, &pk); err != nil {
			return pokemons, err
		}
		pokemons = append(pokemons, pk)
	}
	if err = rows.Err(); err != nil {
		return pokemons, err
	}

	return pokemons, nil
}

func mappingPokemon(row db.RowResultInterface, pk *entity.Pokemon) error {
	return row.Scan(
		&pk.ID, &pk.Name,
		&pk.BattleStats.Health, &pk.BattleStats.MaxHealth, &pk.BattleStats.Attack, &pk.BattleStats.Defense, &pk.BattleStats.Speed,
		&pk.AvatarURL)
}
