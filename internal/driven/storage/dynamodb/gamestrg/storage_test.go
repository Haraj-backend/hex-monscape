package gamestrg

import (
	"context"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveGetGame(t *testing.T) {
	storage := newStorage(t)

	game := entity.Game{
		ID:         uuid.NewString(),
		PlayerName: "Alfat",
		Partner: &entity.Monster{
			ID: "partner-id",
		},
		BattleWon: 2,
		Scenario:  entity.BATTLE_3,
	}

	err := storage.SaveGame(context.Background(), game)
	require.NoError(t, err)

	savedGame, err := storage.GetGame(context.Background(), game.ID)
	require.NoError(t, err)
	require.Equal(t, game, *savedGame)
}

func TestGetGameNotFound(t *testing.T) {
	storage := newStorage(t)
	game, err := storage.GetGame(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, game)
}

func newStorage(t *testing.T) *Storage {
	s, err := New(Config{
		DynamoClient: shared.NewLocalTestDDBClient(),
		TableName:    os.Getenv(shared.TestConfig.EnvKeyGameTableName),
	})
	require.NoError(t, err)

	return s
}
