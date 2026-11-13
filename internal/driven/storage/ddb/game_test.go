package ddb

import (
	"context"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveGetGame(t *testing.T) {
	storage := newGameStorage(t)

	game := entity.Game{
		ID:         uuid.NewString(),
		PlayerName: "Alfat",
		Partner: &entity.Monster{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.ScenarioBattle3,
	}

	err := storage.SaveGame(context.Background(), game)
	require.NoError(t, err)

	savedGame, err := storage.GetGame(context.Background(), game.ID)
	require.NoError(t, err)
	require.Equal(t, game, *savedGame)
}

func TestGetGameNotFound(t *testing.T) {
	storage := newGameStorage(t)
	game, err := storage.GetGame(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, game)
}

func newGameStorage(t *testing.T) *GameStorage {
	s, err := NewGameStorage(GameStorageConfig{
		DynamoClient: newTestDDBClient(),
		TableName:    os.Getenv(envKeyGameTableName),
	})
	require.NoError(t, err)

	return s
}
