package battlestrg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetBattle(t *testing.T) {
	strg := New()
	battle := initNewBattle()
	err := strg.SaveBattle(context.Background(), *battle)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}
	newBattle, err := strg.GetBattle(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to get battle, due: %v", err)
	}
	require.Equal(t, battle, newBattle, "battle is not equal")
}

func TestSaveBattle(t *testing.T) {
	strg := New()
	battle := initNewBattle()
	err := strg.SaveBattle(context.Background(), *battle)
	if err != nil {
		t.Fatalf("unable to save battle, due: %v", err)
	}
	newBattle, err := strg.GetBattle(context.Background(), battle.GameID)
	if err != nil {
		t.Fatalf("unable to get battle, due: %v", err)
	}
	require.Equal(t, battle, newBattle, "battle is not equal")
}

func initNewBattle() *entity.Battle {
	game, _ := entity.NewBattle(entity.BattleConfig{
		GameID:  "b1c87c5c-2ac3-471d-9880-4812552ee15d",
		Partner: newSamplePokemon(),
		Enemy:   newSamplePokemon(),
	})
	return game
}

func newSamplePokemon() *entity.Monster {
	currentTs := time.Now().Unix()
	return &entity.Monster{
		ID:   uuid.NewString(),
		Name: fmt.Sprintf("pokemon_%v", currentTs),
		BattleStats: entity.BattleStats{
			Health:    100,
			MaxHealth: 100,
			Attack:    100,
			Defense:   100,
			Speed:     100,
		},
		AvatarURL: fmt.Sprintf("https://example.com/%v", currentTs),
	}
}
