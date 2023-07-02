package battlestrg

import (
	"context"
	"testing"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveGetBattle(t *testing.T) {
	// init storage & battle
	strg := New()
	expBattle := newBattle()

	// save battle
	err := strg.SaveBattle(context.Background(), *expBattle)
	require.NoError(t, err)

	// get battle
	battle, err := strg.GetBattle(context.Background(), expBattle.GameID)
	require.NoError(t, err)

	// ensure battle is equal to expBattle
	require.Equal(t, expBattle, battle, "unexpected battle")
}

func newBattle() *entity.Battle {
	game, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  uuid.NewString(),
		Partner: util.NewMonster(),
		Enemy:   util.NewMonster(),
	})
	return game
}
