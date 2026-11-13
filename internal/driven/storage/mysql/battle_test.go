package mysql

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveBattle(t *testing.T) {
	// initialize storage
	strg := newBattleStorage(t)

	// save battle
	b := newBattle()
	err := strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)

	// check whether battle exists on database
	savedBattle, err := strg.GetBattle(context.Background(), b.GameID)
	require.NoError(t, err)

	// check whether battle data is match
	require.Equal(t, b, *savedBattle)
}

func TestUpdateBattle(t *testing.T) {
	// initialize storage
	strg := newBattleStorage(t)

	// save battle
	b := newBattle()
	err := strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)

	// update battle state
	b.State = entity.StateEnemyTurn
	err = strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)

	// check whether battle exists on database
	savedBattle, err := strg.GetBattle(context.Background(), b.GameID)
	require.NoError(t, err)

	// check whether battle data is match
	require.Equal(t, b, *savedBattle)
}

func TestGetBattle(t *testing.T) {
	// initialize storage
	strg := newBattleStorage(t)

	// save battle
	b := newBattle()
	err := strg.SaveBattle(context.Background(), b)
	require.NoError(t, err)

	testCases := []struct {
		Name      string
		GameID    string
		ExpBattle *entity.Battle
	}{
		{
			Name:      "Battle Exists",
			GameID:    b.GameID,
			ExpBattle: &b,
		},
		{
			Name:      "Battle Not Exists",
			GameID:    uuid.NewString(),
			ExpBattle: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			battle, err := strg.GetBattle(context.Background(), testCase.GameID)
			require.NoError(t, err)
			require.Equal(t, testCase.ExpBattle, battle)
		})
	}
}

func newBattle() entity.Battle {
	return entity.Battle{
		GameID:  uuid.NewString(),
		State:   entity.StateDecideTurn,
		Partner: testutil.NewTestMonster(),
		Enemy:   testutil.NewTestMonster(),
		LastDamage: entity.LastDamage{
			Partner: 0,
			Enemy:   10,
		},
	}
}

func newBattleStorage(t *testing.T) *BattleStorage {
	// initialize sql client
	sqlClient, err := newTestSQLClient()
	require.NoError(t, err)

	// initialize storage
	strg, err := NewBattleStorage(BattleStorageConfig{SQLClient: sqlClient})
	require.NoError(t, err)

	return strg
}
