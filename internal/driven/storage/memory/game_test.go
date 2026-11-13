package memory

import (
	"context"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/stretchr/testify/require"
)

func TestSaveGetGame(t *testing.T) {
	// init storage & game
	strg := NewGameStorage()
	expGame := newGame()

	// get game, supposedly the returned game is nil
	game, err := strg.GetGame(context.Background(), expGame.ID)
	require.NoError(t, err)
	require.Nil(t, game, "game is not nil")

	// save game
	err = strg.SaveGame(context.Background(), *expGame)
	require.NoError(t, err)

	// get game
	game, err = strg.GetGame(context.Background(), expGame.ID)
	require.NoError(t, err)

	// ensure game is equal to newGame
	require.Equal(t, expGame, game, "unexpected game")
}

func newGame() *entity.Game {
	currentTs := time.Now().Unix()
	game, _ := entity.NewGame(entity.GameConfig{
		PlayerName: "Riandy R.N",
		Partner:    testutil.NewTestMonster(),
		CreatedAt:  currentTs,
	})
	return game
}
