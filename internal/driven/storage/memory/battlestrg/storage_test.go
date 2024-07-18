package battlestrg_test

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/core/testutil"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/battlestrg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveGetBattle(t *testing.T) {
	// init storage & battle
	strg := battlestrg.New()
	expBattle := newBattle()

	// get battle, supposedly the returned battle is nil
	battle, err := strg.GetBattle(context.Background(), expBattle.GameID)
	require.NoError(t, err)
	require.Nil(t, battle, "battle is not nil")

	// save battle
	err = strg.SaveBattle(context.Background(), *expBattle)
	require.NoError(t, err)

	// get battle
	battle, err = strg.GetBattle(context.Background(), expBattle.GameID)
	require.NoError(t, err)

	// ensure battle is equal to expBattle
	require.Equal(t, expBattle, battle, "unexpected battle")
}

func newBattle() *entity.Battle {
	game, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  uuid.NewString(),
		Partner: testutil.NewTestMonster(),
		Enemy:   testutil.NewTestMonster(),
	})
	return game
}
