package gamestrg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
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

func TestSaveGame(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// save Game
	g := newGame()
	err = strg.SaveGame(context.Background(), g)
	require.NoError(t, err)
	// check whether Game exists on database
	savedGame, err := strg.GetGame(context.Background(), g.ID)
	fmt.Println(savedGame)
	require.NoError(t, err)
	// check whether Game data is match
	require.Equal(t, g, *savedGame)
}

func TestGetGame(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// save Game
	g := newGame()
	err = strg.SaveGame(context.Background(), g)
	require.NoError(t, err)
	// override Game battle won
	g.BattleWon = 1
	// save again
	err = strg.SaveGame(context.Background(), g)
	require.NoError(t, err)
	// check whether Game exists on database
	savedGame, err := strg.GetGame(context.Background(), g.ID)
	require.NoError(t, err)
	// check whether Game data is match
	require.Equal(t, g, *savedGame)
}

func TestGetGameNotFound(t *testing.T) {
	// initialize sql client
	sqlClient, err := newSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(shared.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// check whether game exists on database
	savedGame, err := strg.GetGame(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, savedGame)
}

func newGame() entity.Game {
	partner := entity.Pokemon{
		ID:        "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Name:      "Pikachu",
		AvatarURL: "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png",
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    49,
			Defense:   49,
			Speed:     45,
		},
	}

	return entity.Game{
		ID:         uuid.NewString(),
		PlayerName: "player1",
		CreatedAt:  1646205996,
		BattleWon:  0,
		Scenario:   "BATTLE_3",
		Partner:    &partner,
	}
}
