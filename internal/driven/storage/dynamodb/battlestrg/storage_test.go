package battlestrg

import (
	"context"
	"os"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveGetBattle(t *testing.T) {
	// initialize storage
	storage := newStorage(t)

	// save battle
	battle := entity.Battle{
		GameID: uuid.NewString(),
		State:  entity.StateDecideTurn,
		Partner: &entity.Monster{
			ID:   "partner-id",
			Name: "my-partner",
		},
		Enemy: &entity.Monster{
			ID:   "enemy-id",
			Name: "my-enemy",
		},
		LastDamage: entity.LastDamage{
			Partner: 10,
			Enemy:   10,
		},
	}
	err := storage.SaveBattle(context.Background(), battle)
	require.NoError(t, err)

	// get battle
	savedBattle, err := storage.GetBattle(context.Background(), battle.GameID)
	require.NoError(t, err)
	require.Equal(t, battle, *savedBattle)
}

func TestGetBattleNotFound(t *testing.T) {
	// initialize storage
	storage := newStorage(t)

	// get battle from storage
	bt, err := storage.GetBattle(context.Background(), uuid.NewString())
	require.NoError(t, err)
	require.Nil(t, bt)
}

func newStorage(t *testing.T) *Storage {
	s, err := New(Config{
		DynamoClient: shared.NewLocalTestDDBClient(),
		TableName:    os.Getenv(shared.TestConfig.EnvKeyBattleTableName),
	})
	require.NoError(t, err)

	return s
}
