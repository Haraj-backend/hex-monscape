package pokestrg

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestGetPartner(t *testing.T) {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
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

func TestGetPartnerNotFound(t *testing.T) {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// insert pokemon
	p := newPokemon()
	// check whether pokemon exists on database
	savedPokemon, err := strg.GetPartner(context.Background(), p.ID)
	require.NoError(t, err)
	require.Nil(t, savedPokemon)
}

func TestGetAvailablePartners(t *testing.T) {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)

	partner := newPokemon()

	testCases := []struct {
		Name     string
		Partners []entity.Monster
	}{
		{
			Name:     "Test Empty Available Partners",
			Partners: []entity.Monster{},
		},
		{
			Name:     "Test Exists Available Partners",
			Partners: []entity.Monster{partner},
		},
	}

	for _, testCase := range testCases {
		truncateTable(sqlClient)
		t.Run(testCase.Name, func(t *testing.T) {
			for _, p := range testCase.Partners {
				insertPokemon(sqlClient, p, 1)
			}

			fetchedPartners, err := strg.GetAvailablePartners(context.Background())
			require.NoError(t, err)

			switch len(testCase.Partners) {
			case 0:
				require.Nil(t, fetchedPartners)
			case 1:
				require.Equal(t, testCase.Partners[0], fetchedPartners[0])
				truncateTable(sqlClient)
			}
		})
	}
}

func TestGetAvailableEnemies(t *testing.T) {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)

	enemy := newPokemon()

	testCases := []struct {
		Name    string
		Enemies []entity.Monster
	}{
		{
			Name:    "Test Empty Available Enemies",
			Enemies: []entity.Monster{},
		},
		{
			Name:    "Test Exists Available Enemies",
			Enemies: []entity.Monster{enemy},
		},
	}

	for _, testCase := range testCases {
		truncateTable(sqlClient)
		t.Run(testCase.Name, func(t *testing.T) {
			for _, p := range testCase.Enemies {
				insertPokemon(sqlClient, p, 0)
			}

			fetchedEnemies, err := strg.GetPossibleEnemies(context.Background())
			require.NoError(t, err)

			switch len(testCase.Enemies) {
			case 0:
				require.Nil(t, fetchedEnemies)
			case 1:
				require.Equal(t, testCase.Enemies[0], fetchedEnemies[0])
				truncateTable(sqlClient)
			}
		})
	}
}

func truncateTable(db *sqlx.DB) {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("TRUNCATE TABLE monster")
	if err != nil {
		panic(err)
	}
}

func insertPokemon(db *sqlx.DB, p entity.Monster, is_partnerable int) string {
	_, err := db.Exec(
		"INSERT INTO monster (id, name, health, max_health, attack, defense, speed, avatar_url, is_partnerable) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.ID, p.Name, p.BattleStats.Health, p.BattleStats.MaxHealth, p.BattleStats.Attack, p.BattleStats.Defense, p.BattleStats.Speed, p.AvatarURL, is_partnerable,
	)
	if err != nil {
		panic(err)
	}

	return p.ID
}

func newPokemon() entity.Monster {
	return entity.Monster{
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
