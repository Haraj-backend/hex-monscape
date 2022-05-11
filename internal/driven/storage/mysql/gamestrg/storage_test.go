package gamestrg

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Haraj-backend/hex-pokebattle/internal/core/entity"
	"github.com/Haraj-backend/hex-pokebattle/internal/driven/storage/mysql/shared"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

func TestSaveGameGetGame(t *testing.T) {
	// initialize sql client
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// create new test pokemon
	partner := shared.NewTestPokemon()
	err = shared.InsertTestPokemon(sqlClient, partner, 1)
	require.NoError(t, err)
	// save Game
	g := newGame(partner)
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
	sqlClient, err := shared.NewTestSQLClient()
	require.NoError(t, err)
	// initialize storage
	strg, err := New(Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// check whether game exists on database
	savedGame, err := strg.GetGame(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, savedGame)
}

func newGame(partner entity.Pokemon) entity.Game {
	return entity.Game{
		ID:         uuid.NewString(),
		PlayerName: "player1",
		CreatedAt:  1646205996,
		BattleWon:  0,
		Scenario:   "BATTLE_3",
		Partner:    &partner,
	}
}
