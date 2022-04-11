package pokestrg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

const envKeySQLDSN = "SQL_DSN"

func newSQLClient() (*sqlx.DB, error) {
	sqlDSN := os.Getenv(envKeySQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
	}
	return sqlClient, nil
}

func TestGetPartner(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// insert pokemon
	p := newPokemon()
	id := insertPokemon(sqlClient, p, 1)
	// check whether pokemon exists on database
	savedPokemon, err := strg.GetPartner(context.Background(), id)
	require.NoError(t, err)
	// check whether pokemon data is match
	require.Equal(t, p, *savedPokemon)
}

func insertPokemon(db *sqlx.DB, p entity.Pokemon, is_partnerable int) string {
	_, err := db.Exec(
		"INSERT INTO pokemons (id, name, health, max_health, attack, defense, speed, avatar_url, is_partnerable) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.ID, p.Name, p.BattleStats.Health, p.BattleStats.MaxHealth, p.BattleStats.Attack, p.BattleStats.Defense, p.BattleStats.Speed, p.AvatarURL, is_partnerable,
	)
	if err != nil {
		panic(err)
	}

	return p.ID
}

func newPokemon() entity.Pokemon {
	return entity.Pokemon{
		ID:        uuid.NewString(),
		Name:      "Lala",
		AvatarURL: "https://example.com/025.png",
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    49,
			Defense:   49,
			Speed:     45,
		},
	}
}
