package shared

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const envKeySQLDSN = "TEST_SQL_DSN"

func NewTestSQLClient() (*sqlx.DB, error) {
	sqlDSN := os.Getenv(envKeySQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize sql client due: %w", err)
	}
	return sqlClient, nil
}

func InsertTestPokemon(db *sqlx.DB, p entity.Pokemon, isPartnerable int) error {
	_, err := db.Exec(
		"REPLACE INTO pokemons (id, name, health, max_health, attack, defense, speed, avatar_url, is_partnerable) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.ID, p.Name, p.BattleStats.Health, p.BattleStats.MaxHealth, p.BattleStats.Attack, p.BattleStats.Defense, p.BattleStats.Speed, p.AvatarURL, isPartnerable,
	)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}

func NewTestPokemon() entity.Pokemon {
	return entity.Pokemon{
		ID:        uuid.NewString(),
		Name:      fmt.Sprintf("Test_Pokemon_%v", rand.New(rand.NewSource(time.Now().UnixMilli())).Int63()),
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
